package domain

type Paginated[Result any] struct {
	Results    []Result
	IsLastPage bool
}

func NewEmptyPage[Result any]() Paginated[Result] {
	return Paginated[Result]{
		Results:    []Result{},
		IsLastPage: true,
	}
}

func NewPage[Result any](data []Result, isLastPage bool) Paginated[Result] {
	return Paginated[Result]{
		Results:    data,
		IsLastPage: isLastPage,
	}
}

func (self Paginated[Result]) Len() int {
	return len(self.Results)
}
