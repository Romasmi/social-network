package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
	"github.com/google/uuid"
)

type ParsedPost struct {
	ProfileId uuid.UUID
	Post      string
}

func (a *App) importPosts(context context.Context, link string, profileIdsRaw []string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed request: %v", resp.StatusCode)
	}

	profileIds := make([]uuid.UUID, len(profileIdsRaw))
	for i, v := range profileIdsRaw {
		id, _ := uuid.Parse(v)
		profileIds[i] = id
	}

	scanner := bufio.NewScanner(resp.Body)

	var wg sync.WaitGroup
	const maxWorkers = 50
	jobs := make(chan *ParsedPost, 100)
	errCh := make(chan error, 10)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer func() {
				if r := recover(); r != nil {
					errCh <- fmt.Errorf("recovered in goroutine: %w", r)
				}
			}()
			defer wg.Done()
			for post := range jobs {
				if err := a.importPost(context, post); err != nil {
					errCh <- fmt.Errorf("error while post creation: %v", err)
				}
			}
		}(i)
	}

	go func() {
		for scanner.Scan() {
			line := scanner.Text()

			if strings.TrimSpace(line) == "" {
				continue
			}

			post := &ParsedPost{
				ProfileId: profileIds[rand.Intn(len(profileIds))],
				Post:      line,
			}
			jobs <- post
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errors []string
	for err := range errCh {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		return fmt.Errorf("import completed with errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (a *App) importPost(context context.Context, post *ParsedPost) error {
	postRepo := repository.CreatePostsRepository(a.DbConn.DB, a.RedisConn.Rdb)
	postService := services.CreatePostService(postRepo)
	_, err := postService.CreatePost(context, post.ProfileId, post.Post)
	if err != nil {
		return err
	}
	return nil
}
