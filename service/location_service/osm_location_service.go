package locationservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/franchises/domain"
	httpheaders "example.com/franchises/service/http_headers"
)

const (
	osmBaseUrl = "https://nominatim.openstreetmap.org/search?"

	osmDataOrigin = "osm"

	osmSearchQueryParamName         = "q"
	osmAddressDetailsQueryParamName = "addressdetails"
	osmFormatQueryParamName         = "format"
	osmLimitQueryParamName          = "limit"
	osmEmailQueryParamName          = "email"
	osmCountryCodeQueryParamName    = "countrycodes"
	osmExclidesQueryParamName       = "exclude_place_ids"

	osmLocationSearchMaxItemsPerRequest = 50
)

type osmService struct {
	userAgent         string
	preferredLanguage string
	client            *http.Client

	countryCodes []string
	excludingIds []domain.LocationId
}

func NewOsmLocationService(
	userAgent string,
	preferredLanguage string,
	countryCodes []string,
	excludingIds []domain.LocationId,
) *osmService {
	return &osmService{
		userAgent:         userAgent,
		preferredLanguage: preferredLanguage,
		client:            &http.Client{},

		countryCodes: countryCodes,
		excludingIds: excludingIds,
	}
}

func (self *osmService) SleepBetweenRequests() {
	// Force the caller to wait a bit to avoid the trigger-happy pagination
	// firing too quickly.
	//
	// This is part of Nominatim's usage policy:
	//
	//     https://operations.osmfoundation.org/policies/nominatim/

	time.Sleep(1 * time.Second)
}

func (self *osmService) SearchLocation(
	query string,
) (domain.Paginated[domain.Location], error) {
	req, err := self.buildBaseRequest()
	if err != nil {
		return domain.NewEmptyPage[domain.Location](), err
	}

	urlQueryParams := req.URL.Query()
	urlQueryParams.Add(osmSearchQueryParamName, query)
	urlQueryParams.Add(osmAddressDetailsQueryParamName, "1")
	urlQueryParams.Add(osmFormatQueryParamName, "json")
	if len(self.countryCodes) > 0 {
		urlQueryParams.Add(osmCountryCodeQueryParamName, strings.Join(self.countryCodes, ","))
	}
	urlQueryParams.Add(osmLimitQueryParamName, strconv.Itoa(osmLocationSearchMaxItemsPerRequest))

	if len(self.excludingIds) > 0 {
		excludeIds := make([]string, len(self.excludingIds))
		for index, id := range self.excludingIds {
			excludeIds[index] = strconv.FormatUint(id, 10)
		}
		urlQueryParams.Add(osmExclidesQueryParamName, strings.Join(excludeIds, ","))
	}

	req.URL.RawQuery = urlQueryParams.Encode()

	fmt.Println(req.URL.String())

	response, err := self.client.Do(req)
	if err != nil {
		return domain.NewEmptyPage[domain.Location](), err
	}

	return self.parseSearchResponse(response)
}

func (self osmService) buildBaseRequest() (request *http.Request, err error) {
	request, err = http.NewRequest("GET", osmBaseUrl, nil)
	if err != nil {
		return
	}

	request.Header.Set(httpheaders.UserAgentHeader, self.userAgent)
	request.Header.Set(httpheaders.AcceptLanguageHeader, self.preferredLanguage+", *;q=0.7")

	return
}

func (self *osmService) parseSearchResponse(
	response *http.Response,
) (domain.Paginated[domain.Location], error) {
	defer response.Body.Close()

	var locationsDto []osmLocationDto
	err := json.NewDecoder(response.Body).Decode(&locationsDto)
	if err != nil {
		return domain.NewEmptyPage[domain.Location](), err
	}

	locationCount := len(locationsDto)

	locations := make([]domain.Location, locationCount)
	for index, dto := range locationsDto {
		newLocation := dto.asLocation()

		locations[index] = newLocation
		self.excludingIds = append(self.excludingIds, newLocation.Id)
	}

	return domain.NewPage[domain.Location](
		locations,
		locationCount < osmLocationSearchMaxItemsPerRequest,
	), nil
}
