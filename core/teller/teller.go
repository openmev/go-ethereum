package teller

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Teller struct {
	core *tellerCore

	isMutate bool
	isFound  bool
	mutated  bool
}

// NewTeller returns a teller object that wraps the global shared tellerCore
// isMuate is true the teller would mutate the return data.
func NewTeller(isMutate bool) *Teller {
	return &Teller{
		core:     newTellerCore(),
		isMutate: isMutate,
		mutated:  false,
	}
}

func (t *Teller) Stop() {
	t.core.stop()
}

func (t *Teller) InsertMutateState(txHash common.Hash, detail MutateDetail) {
	t.core.insertMutateState(txHash, detail)
}

// callTrace is the result of a callTracer run.
type callTrace struct {
	Type    string          `json:"type"`
	From    common.Address  `json:"from"`
	To      common.Address  `json:"to"`
	Input   hexutil.Bytes   `json:"input"`
	Output  hexutil.Bytes   `json:"output"`
	Gas     *hexutil.Uint64 `json:"gas,omitempty"`
	GasUsed *hexutil.Uint64 `json:"gasUsed,omitempty"`
	Value   *hexutil.Big    `json:"value,omitempty"`
	Error   string          `json:"error,omitempty"`
	Calls   []callTrace     `json:"calls,omitempty"`
}

const LogDetailPath = "/home/jonah1005/contract/tellerDetail/"

func (t *Teller) LogDetail(result json.RawMessage, txHash common.Hash) {
	// ioutil.WriteFile(fmt.Sprintf("%s/%s", LogDetailPath, txHash.Hex()), result, 0644)
}

func (t *Teller) CheckAndMutate(res []byte, caller common.Address, callee common.Address, input []byte, txHash common.Hash, txOrigin common.Address, blockNumber int64) (ret []byte) {
	if t.isMutate {
		ret, mutated := t.core.checkAndMutate(res, caller, callee, input, txHash, txOrigin, blockNumber)
		if mutated {
			t.mutated = true
		}
		return ret
	}
	return res
}

func (t *Teller) Mutated() bool {
	return t.mutated
}

func (t *Teller) IsFound() bool {
	return t.isFound
}

func (t *Teller) IsMutate() bool {
	return t.isMutate
}

func (t *Teller) ResetMutateMapList() {
	t.core.setMutateMapList(nil)
}

// SetMutateMap sets map for teller's mutator
func (t *Teller) SetMutateMapList(mutateMapList *MutateMapList) {
	if mutateMapList != nil {
		t.core.setMutateMapList(mutateMapList)

	}
}
