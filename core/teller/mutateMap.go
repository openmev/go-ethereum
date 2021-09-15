package teller

import "github.com/ethereum/go-ethereum/common"

// MutatorType defines mutators type.  Each type has a specific rules.
type MutatorType int64

// MutatorType enumerator
const (
	UndefinedMutatorType MutatorType = iota
	GetReserveMutator
)

// MutateMap defines breakpoints where the teller mutators would mutate
type MutateMap struct {
	Address common.Address
	Rate    string
}

// MutateMapList is a list of MutateMap which indicates the address and mutate ratio for
// the teller's mutator.
type MutateMapList []MutateMap
