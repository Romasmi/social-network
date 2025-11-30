package cli

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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
	i := 0
	for {
		i++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println()

		if i > 20 {
			break
		}

		parsed, err := recordIntoModel(record)
		if err != nil {
			fmt.Printf("error while parsing a line: %v", err)
		}
		fmt.Println("parsed", parsed)
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
