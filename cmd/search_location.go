package cmd

import (
	"fmt"
	"sync"

	"example.com/franchises/db"
	"example.com/franchises/domain"
	"example.com/franchises/service"
)

type searchResult struct {
	query   string
	page    domain.Paginated[domain.Location]
	service service.LocationService
}

type searchLocationCmd struct {
	query       string
	services    []service.LocationService
	resultsChan chan searchResult
	errorsChan  chan error
}

func NewSearchLocationCmd(
	query string,
	db db.LocationSaver,
	services []service.LocationService,
) Command {
	chanBufferSize := len(services)

	resultsChan := make(chan searchResult, chanBufferSize)
	errorsChan := make(chan error, chanBufferSize)

	go processResults(db, resultsChan)
	go processErrors(errorsChan)

	return searchLocationCmd{
		query:       query,
		services:    services,
		resultsChan: resultsChan,
		errorsChan:  errorsChan,
	}
}

func (self searchLocationCmd) Run() error {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(self.services))

	for _, service := range self.services {
		go runSearch(
			self.resultsChan,
			self.errorsChan,
			service,
			self.query,
			&waitGroup,
		)
	}

	waitGroup.Wait()
	close(self.resultsChan)
	close(self.errorsChan)

	return nil
}

func runSearch(
	resultsChan chan<- searchResult,
	errorsChan chan<- error,
	s service.LocationService,
	query string,
	waitGroup *sync.WaitGroup,
) {
	result, err := s.SearchLocation(query)
	if err != nil {
		errorsChan <- err
		close(errorsChan)
		close(resultsChan)

		return
	}

	resultsChan <- searchResult{
		query:   query,
		page:    result,
		service: s,
	}
	if !result.IsLastPage {
		go runSearch(resultsChan, errorsChan, s, query, waitGroup)
	} else {
		waitGroup.Done()
	}
}

func processResults(
	database db.LocationSaver,
	resultsChan <-chan searchResult,
) {
	for result := range resultsChan {
		fmt.Printf("Got page with %d elements\n", result.page.Len())

		for _, location := range result.page.Results {
			database.SaveLocation(location)
		}
	}
}

func processErrors(errorsChan <-chan error) {
	for error := range errorsChan {
		fmt.Println(error)
	}
}
