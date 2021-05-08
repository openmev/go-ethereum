package teller

import "gorm.io/gorm"

type BreakPointType int64

const (
	BreakPointTypeUndefined BreakPointType = 0
	BreakPointTypeUniswapGetReserve
	BreakPointTypeSushiswapGetReserve
)

type TellerLog struct {
	gorm.Model
	// TxHash      common.Hash
	TxHash      string
	Origin      string
	Caller      string
	Callee      string
	Input       string
	BlockNumber int64
	Mutated     bool

	ObservedCount  int64
	BreakPointType BreakPointType

	MutateDetail MutateDetail `gorm:"embedded"`
}

type MutateDetail struct {
	IsDifference bool
	TxStatus     bool
	TxErrMsg     string
}

type DBWatchAddress struct {
	Signature string
	Address   string
	Type      string
	Remark    string
}
