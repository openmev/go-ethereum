package teller

import (
	"bytes"
	"encoding/json"
	"math"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type WatchAddress struct {
	Address   common.Address
	Signature []byte
}

type UniswapData struct {
	Pairs []UniswapPair `jons:"pairs"`
}

type UniswapPair struct {
	ID string `json:"id"`
}

var globalTellerCore *tellerCore
var once sync.Once

type WatchFunction struct {
	Signature []byte
	Address   map[string]BreakPointType
}

// core Teller that all tellers share.
type tellerCore struct {
	WatchList []WatchFunction
	Log       []TellerLog
	mu        *sync.Mutex

	logIndex int
	logSize  int

	mutateMapList *MutateMapList

	// map[txHash]TxInfo
	txInfoCache map[string]txInfo
}

type txInfo struct {
	observedCount int64
}

func newTellerCore() *tellerCore {
	once.Do(func() {
		logSize := 100
		data := struct {
			Data UniswapData `json:"data"`
		}{}
		json.Unmarshal([]byte(uniswapParisJSON), &data)

		// 0902f1ac  =>  getReserves()
		// 5909c0d5  =>  price0CumulativeLast()
		// 5a3d5493  =>  price1CumulativeLast()
		// 7464fc3d  =>  kLast()
		// getReserve := WatchFunction{
		// 	Signature: common.FromHex("0902f1ac"),
		// 	Address:   make(map[string]bool),
		// }
		// for _, pair := range data.Data.Pairs {

		// 	getReserve.Address[common.HexToAddress(pair.ID).Hex()] = true
		// }
		globalTellerCore = &tellerCore{
			WatchList:   []WatchFunction{},
			mu:          &sync.Mutex{},
			Log:         make([]TellerLog, logSize),
			logSize:     logSize,
			logIndex:    0,
			txInfoCache: make(map[string]txInfo),
		}
		globalTellerCore.loadConstantFunc()
		// globalTellerCore.loadWatchAddressFromDB("0x0902f1ac", 0)
	})
	return globalTellerCore
}

func (w WatchFunction) Match(address common.Address, input []byte) (BreakPointType, bool) {
	if len(input) < len(w.Signature) {
		return 0, false
	}
	if bytes.Equal(w.Signature, input[:len(w.Signature)]) {
		return BreakPointTypeUndefined, true
	}
	return 0, false
}

func (t *tellerCore) stop() {
}

func DecodeHelper(contractAbi string, signature []byte, ret []byte) (interface{}, error) {
	abi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return nil, err
	}
	method, err := abi.MethodById(signature)
	if err != nil {
		return nil, err
	}
	return abi.Unpack(method.Name, ret)
}

func decodeInputHelper(contractAbi string, input []byte) (interface{}, error) {
	abi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return nil, err
	}
	method, err := abi.MethodById(input[:4])
	if err != nil {
		return nil, err
	}
	return method.Inputs.Unpack(input[4:])
}

func encodeHelper(contractAbi string, signature []byte, args []interface{}) ([]byte, error) {
	abi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return nil, err
	}
	method, err := abi.MethodById(signature)
	if err != nil {
		return nil, err
	}
	return method.Outputs.PackValues(args)
}

func (t *tellerCore) insertMutateState(txHash common.Hash, detail MutateDetail) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, l := range t.Log {
		if l.TxHash == txHash.Hex() {
			t.Log[i].MutateDetail = detail
			t.Log[i].Mutated = true
		}
	}
}

func (t *tellerCore) setMutateMapList(mutateMapList *MutateMapList) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.mutateMapList = mutateMapList
}

func (t *tellerCore) clearMutateMapList() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.mutateMapList = nil
}

func mulFloat(a *big.Int, rate *big.Float) *big.Int {
	precision := 6
	denominator := math.Pow10(precision)
	rateInt, _ := rate.Mul(rate, big.NewFloat(denominator)).Int64()
	mul := big.NewInt(rateInt)

	a = a.Mul(a, mul)
	a = a.Div(a, big.NewInt(int64(denominator)))
	return a
}

func quoFloat(a *big.Int, rate *big.Float) *big.Int {
	precision := 6
	denominator := math.Pow10(precision)
	rateInt, _ := rate.Quo(rate, big.NewFloat(denominator)).Int64()
	mul := big.NewInt(rateInt)
	a = a.Mul(a, mul)
	return a
}

func mutateFloat(a *big.Int, b *big.Int, rate *big.Float) (*big.Int, *big.Int) {
	precision := 6
	denominator := math.Pow10(precision)
	rateInt, _ := rate.Mul(rate, big.NewFloat(denominator)).Int64()
	mul := big.NewInt(rateInt)

	a = a.Mul(a, mul)
	a = a.Div(a, big.NewInt(int64(denominator)))

	b = b.Mul(b, big.NewInt(int64(denominator)))
	b = b.Div(b, mul)
	return a, b
}

func (t *tellerCore) checkAndMutate(res []byte, caller common.Address, callee common.Address, input []byte, txHash common.Hash, txOrigin common.Address, blockNumber int64) (ret []byte, isMutate bool) {
	if len(input) >= 4 {
		// getReserve
		if bytes.Equal(input[:4], common.FromHex("0x0902f1ac")) {
			// if the mutateMap is define, the mutator only mutates specific calls.
			return t.mutateGetReserve(res, caller, callee, input)
		}

		// calc_withdraw_one_coin(address,uint256,int128)
		if bytes.Equal(input[:4], common.FromHex("0x41b028f3")) {
			return t.mutateCalcWithdrawOneCoin(res, caller, callee, input)
		}

		// calc_withdraw_one_coin(uint256,int128)
		if bytes.Equal(input[:4], common.FromHex("0xcc2b27d7")) {
			return t.mutateCalcWithdrawOneCoinStableSwap(res, caller, callee, input)
		}

		// calc_token_amount(uint256[2], bool)
		if bytes.Equal(input[:4], common.FromHex("0xed8e84f3")) {
			return t.mutateCalcTokenAmount(res, caller, callee, input)
		}

		// calc_token_amount(uint256[3], bool)
		if bytes.Equal(input[:4], common.FromHex("0x3883e119")) {
			return t.mutateCalcTokenAmount3Crv(res, callee, callee, input)
		}
		// calc_token_amount(address,uint256[4],bool)
		if bytes.Equal(input[:4], common.FromHex("0x861cdef0")) {
			// we simply use calWithdrawOneCoin as its the same
			return t.mutateCalcWithdrawOneCoin(res, caller, callee, input)
		}

		// getTokenToEthInputPrice(uint256)
		if bytes.Equal(input[:4], common.FromHex("0x95b68fe7")) {
			// we simply use calWithdrawOneCoin as its the same
			return t.mutateTokenToEthInputPrice(res, caller, callee, input)
		}

		// getEthToTokenInputPrice(uint256)
		if bytes.Equal(input[:4], common.FromHex("0xcd7724c3")) {
			// we simply use calWithdrawOneCoin as its the same
			return t.mutateTokenToEthInputPrice(res, caller, callee, input)
		}

		// getExpectedRate(address,address,uint256)
		if bytes.Equal(input[:4], common.FromHex("0x809a9e55")) {
			return t.mutateKyberGetExpectedRate(res, caller, callee, input)
		}

		// else if bytes.Compare(input[:4], common.FromHex("0x5a3d5493")) == 0 {
		// 	 if ret, err := DecodeHelper(input[:4], res); err == nil {
		// 	 	fmt.Printf("Type: %T, %v", ret, ret)
		// 	 }
		// }
	}
	return res, false
}
