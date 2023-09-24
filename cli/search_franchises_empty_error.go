package cli

type searchFranchisesEmptySearchError struct {
}

func SearchFranchisesEmptySearchError() searchFranchisesEmptySearchError {
	return searchFranchisesEmptySearchError{}
}

func (searchFranchisesEmptySearchError) Error() string {
	return "Cannot search for an empty string."
}
