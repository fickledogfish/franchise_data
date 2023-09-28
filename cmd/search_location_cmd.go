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

	db               LocationSaver
	services         []LocationService
	addressRecoverer AddressRecoverer

	resultsChan chan searchResult
	errorsChan  chan error
}

func NewSearchLocationCmd(
	query string,
	db LocationSaver,
	services []LocationService,
	addressRecoverer AddressRecoverer,
) searchLocationCmd {
	chanBufferSize := len(services)

	resultsChan := make(chan searchResult, chanBufferSize)
	errorsChan := make(chan error, chanBufferSize)

	return searchLocationCmd{
		query: query,

		db:               db,
		services:         services,
		addressRecoverer: addressRecoverer,

		resultsChan: resultsChan,
		errorsChan:  errorsChan,
	}
}

func (self searchLocationCmd) Run() error {
	var processWaitGroup sync.WaitGroup
	go processResults(self.db, self.addressRecoverer, self.resultsChan, self.errorsChan, &processWaitGroup)
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
	addressRecoverer AddressRecoverer,
	resultsChan <-chan searchResult,
	errorsChan chan<- error,
	waitGroup *sync.WaitGroup,
) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	for result := range resultsChan {
		log.Info("Got page with %d elements", result.page.Len())

		recoverAddressesInResult(addressRecoverer, &result, errorsChan)

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

func recoverAddressesInResult(
	addressRecoverer AddressRecoverer,
	result *searchResult,
	errorsChan chan<- error,
) {
	if addressRecoverer == nil {
		return // The user didn't ask to recover, so don't.
	}

	for index, location := range result.page.Results {
		if location.Address.PostalCode == "" {
			log.Info("Found location with missing postal code, so it cannot be recovered: %#v", location)
			continue // Not much else we can do here
		}

		if location.Address.Street == "" || location.Address.City == "" || location.Address.State == "" {
			newAddress, err := addressRecoverer.RecoverInfoFromPostalCode(location.Address.PostalCode)
			if err != nil {
				errorsChan <- err
				continue
			}

			log.Debug("Recovered info for %s", location.Address.PostalCode)

			result.page.Results[index].Address = newAddress
		}
	}
}
