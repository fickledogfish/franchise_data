package cmd

import "example.com/franchises/domain"

type AddressRecoverer interface {
	RecoverInfoFromPostalCode(string) (domain.Address, error)
}
