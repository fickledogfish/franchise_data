package cmd

import (
	"sync"

	"example.com/franchises/domain"
	"example.com/franchises/log"
)

type searchResult struct {
	query string
	page  domain.Paginated[domain.Location]
}

type searchLocationCmd struct {
	query string

	db       LocationSaver
	services []LocationService

	resultsChan chan searchResult
	errorsChan  chan error
}

func NewSearchLocationCmd(
	query string,
	db LocationSaver,
	services []LocationService,
) searchLocationCmd {
	chanBufferSize := len(services)

	resultsChan := make(chan searchResult, chanBufferSize)
	errorsChan := make(chan error, chanBufferSize)

	return searchLocationCmd{
		query: query,

		db:       db,
		services: services,

		resultsChan: resultsChan,
		errorsChan:  errorsChan,
	}
}

func (self searchLocationCmd) Run() error {
	var processWaitGroup sync.WaitGroup
	go processResults(self.db, self.resultsChan, &processWaitGroup)
	go processErrors(self.errorsChan, &processWaitGroup)

	var servicesWaitGroup sync.WaitGroup
	servicesWaitGroup.Add(len(self.services))

	for _, service := range self.services {
		go runSearch(
			self.resultsChan,
			self.errorsChan,
			service,
			self.query,
			&servicesWaitGroup,
		)
	}

	servicesWaitGroup.Wait()

	close(self.resultsChan)
	close(self.errorsChan)
	processWaitGroup.Wait()

	return nil
}

func runSearch(
	resultsChan chan<- searchResult,
	errorsChan chan<- error,
	service LocationService,
	query string,
	waitGroup *sync.WaitGroup,
) {
	result, err := service.SearchLocation(query)
	if err != nil {
		errorsChan <- err
		waitGroup.Done()
		return
	}

	resultsChan <- searchResult{
		query: query,
		page:  result,
	}
	if !result.IsLastPage {
		service.SleepBetweenRequests()
		go runSearch(resultsChan, errorsChan, service, query, waitGroup)
	} else {
		waitGroup.Done()
	}
}

func processResults(
	database LocationSaver,
	resultsChan <-chan searchResult,
	waitGroup *sync.WaitGroup,
) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	for result := range resultsChan {
		log.Info("Got page with %d elements\n", result.page.Len())

		for _, location := range result.page.Results {
			database.SaveLocation(location)
		}
	}
}

func processErrors(
	errorsChan <-chan error,
	waitGroup *sync.WaitGroup,
) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	for error := range errorsChan {
		log.Error("%s", error)
	}
}
