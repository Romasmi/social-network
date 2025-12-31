package cli

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Romasmi/social-network/internal/domain/profile"
)

type ParsedUser struct {
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	Birthdate  time.Time `json:"birthdate"`
	City       string    `json:"city"`
}

func (h *ImportHandler) ImportUserByLink(context context.Context, link string) error {
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

	reader := csv.NewReader(resp.Body)
	var wg sync.WaitGroup
	const maxWorkers = 100
	jobs := make(chan *ParsedUser, 100)
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
			for user := range jobs {
				if err := h.importUser(context, user); err != nil {
					errCh <- fmt.Errorf("error while user creation: %v", err)
				}
			}
		}(i)
	}

	go func() {
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errCh <- err
				break
			}

			parsed, err := recordIntoModel(record)
			if err != nil {
				errCh <- fmt.Errorf("error while parsing h line: %v", err)
			}

			jobs <- parsed
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

func recordIntoModel(record []string) (*ParsedUser, error) {
	if len(record) != 3 {
		return nil, fmt.Errorf("invalid record")
	}

	user := ParsedUser{}

	fullname := record[0]
	fullnameParts := strings.Split(fullname, " ")
	if len(fullnameParts) != 2 {
		return nil, fmt.Errorf("invalid fullname: %v", fullname)
	}
	user.FirstName = fullnameParts[0]
	user.SecondName = fullnameParts[1]
	birthdate, err := time.Parse("2006-01-02", record[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing birthdate %v: %w", record[1], err)
	}
	user.Birthdate = birthdate
	user.City = record[2]

	return &user, err
}

func (h *ImportHandler) importUser(context context.Context, user *ParsedUser) error {
	payload := &profile.CreateProfileModel{
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Birthdate:  user.Birthdate,
		City:       user.City,
	}
	_, err := h.userService.RegisterUser(context, payload)
	return err
}
