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

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/Romasmi/social-network/internal/services"
)

type ParsedUser struct {
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	Birthdate  time.Time `json:"birthdate"`
	City       string    `json:"city"`
}

func (a *App) importUserByLink(link string) error {
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
	errCh := make(chan error, 10)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		parsed, err := recordIntoModel(record)
		if err != nil {
			fmt.Printf("error while parsing a line: %v", err)
		}

		// TODO replace with pool worker pool
		wg.Add(1)
		go func(user *ParsedUser) {
			defer wg.Done()
			if err := a.importUser(parsed); err != nil {
				errCh <- fmt.Errorf("error while user creation: %v", err)
			}
		}(parsed)
	}

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

func (a *App) importUser(user *ParsedUser) error {
	cityRepo := repository.CreateCityRepository(a.DbConn.DB)
	profileRepo := repository.CreateProfileRepository(a.DbConn.DB)
	uow := repository.CreateUnitOfWork(a.DbConn.DB)
	userService := services.CreateUserService(cityRepo, profileRepo, uow)

	payload := &models.CreateProfileModel{
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Birthdate:  user.Birthdate,
		City:       user.City,
	}
	_, err := userService.RegisterUser(context.Background(), payload)
	return err
}
