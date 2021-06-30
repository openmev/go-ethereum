package teller

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (t *tellerCore) mutateInverseCalcWithdrawOneCoin(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}

	if ret, err := DecodeHelper(crv_deposit_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		amount, ok := args[0].(*big.Int)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			amount = quoFloat(amount, rate)
			args[0] = amount
			if res, err := encodeHelper(crv_deposit_abi, input[:4], args); err == nil {
				return res, true
			}
		}
	}
	return nil, false
}

func (t *tellerCore) mutateCalcTokenAmount3Crv(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	//calc_token_amount
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	if ret, err := DecodeHelper(crv3_stable_swap_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		amount, ok := args[0].(*big.Int)
		fmt.Println("before mutate", amount)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			amount = mulFloat(amount, rate)

			args[0] = amount
			if res, err := encodeHelper(crv3_stable_swap_abi, input[:4], args); err == nil {
				fmt.Println("after mutate", amount)
				return res, true
			}
		}
	}
	return res, false
}

func (t *tellerCore) mutateCalcTokenAmount(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	//calc_token_amount
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	if ret, err := DecodeHelper(crv_stable_swap_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		amount, ok := args[0].(*big.Int)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			amount = mulFloat(amount, rate)
			args[0] = amount
			if res, err := encodeHelper(crv_stable_swap_abi, input[:4], args); err == nil {
				return res, true
			}
		}
	}
	return nil, false
}

func (t *tellerCore) mutateEthToTokenInputPrice(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	if ret, err := DecodeHelper(uniswap_v1_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		_, ok := args[0].(*big.Int)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			fmt.Println("before mutating", args[0])
			_, args[0] = mutateFloat(big.NewInt(10000000), args[0].(*big.Int), rate)
			fmt.Println("after mutating", mutateRate, args[0])
			if res, err := encodeHelper(uniswap_v1_abi, input[:4], args); err == nil {
				return res, true
			}
		}
	}
	return res, false
}

func (t *tellerCore) mutateKyberGetExpectedRate(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {

	getRateFromSrc := func(src common.Address, mapList *MutateMapList) (bool, *big.Float) {
		if mapList != nil {
			for _, mutateMap := range *t.mutateMapList {
				if mutateMap.Address.Hex() == src.Hex() {
					rate, _ := big.NewFloat(0).SetString(mutateMap.Rate)
					return true, rate
				}
			}
		}
		rate, _ := big.NewFloat(0).SetString("1.05")
		return false, rate
	}

	if inputArgs, err := decodeInputHelper(kyber_network_abi, input); err == nil {

		args := inputArgs.([]interface{})
		src, _ := args[0].(common.Address)
		ok, rate := getRateFromSrc(src, t.mutateMapList)
		if ok {

			ret, err := DecodeHelper(kyber_network_abi, input[:4], res)
			if err != nil {
				return res, false
			}

			returnArgs := ret.([]interface{})

			fmt.Println("before mutating", returnArgs[0])
			returnArgs[0], _ = mutateFloat(returnArgs[0].(*big.Int), big.NewInt(10000000), rate)
			fmt.Println("before mutating to eth input price", rate, returnArgs[0])
			returnArgs[1], _ = mutateFloat(returnArgs[1].(*big.Int), big.NewInt(10000000), rate)
			if res, err := encodeHelper(kyber_network_abi, input[:4], returnArgs); err == nil {
				return res, true
			}
		}
	}
	return res, false
}

func (t *tellerCore) mutateTokenToEthInputPrice(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	if ret, err := DecodeHelper(uniswap_v1_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		_, ok := args[0].(*big.Int)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			fmt.Println("before mutating", args[0])
			args[0], _ = mutateFloat(args[0].(*big.Int), big.NewInt(10000000), rate)
			fmt.Println("before mutating to eth input price", mutateRate, args[0])

			if res, err := encodeHelper(uniswap_v1_abi, input[:4], args); err == nil {
				return res, true
			}
		}
	}
	return res, false
}

func (t *tellerCore) mutateCalcWithdrawOneCoinStableSwap(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	fmt.Println("mutate calc withdraw one coin", caller.Hex(), mutateRate)

	if ret, err := DecodeHelper(crv_stable_swap_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		amount, ok := args[0].(*big.Int)
		fmt.Println("current amount", amount)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			amount = mulFloat(amount, rate)
			args[0] = amount
			if res, err := encodeHelper(crv_stable_swap_abi, input[:4], args); err == nil {
				fmt.Println("after mutate", amount, "rate", rate)
				return res, true
			} else {
				fmt.Println("error", err)
			}
		}
	}
	return res, false
}

func (t *tellerCore) mutateCalcWithdrawOneCoin(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	fmt.Println("mutate calc withdraw one coin", caller.Hex(), mutateRate)

	if ret, err := DecodeHelper(crv_deposit_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		amount, ok := args[0].(*big.Int)
		fmt.Println("current amount", amount)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			amount = mulFloat(amount, rate)
			args[0] = amount
			if res, err := encodeHelper(crv_deposit_abi, input[:4], args); err == nil {
				fmt.Println("after mutate", amount)
				return res, true
			}
		}
	}
	return nil, false
}

func (t *tellerCore) mutateGetReserve(res []byte, caller common.Address, callee common.Address, input []byte) (ret []byte, isMutate bool) {
	mutateRate := "1"
	if t.mutateMapList != nil {
		isMatch := false
		for _, mutateMap := range *t.mutateMapList {
			if mutateMap.Address.Hex() == callee.Hex() {
				isMatch = true
				mutateRate = mutateMap.Rate
				break
			}
		}

		if !isMatch {
			// fmt.Println("skip this because the caller does not match", callee.Hex())
			return res, false
		}
	}
	if ret, err := DecodeHelper(uniswap_pair_abi, input[:4], res); err == nil {
		args := ret.([]interface{})
		_, ok := args[0].(*big.Int)
		if ok && mutateRate != "1" {
			rate, _ := big.NewFloat(0).SetString(mutateRate)
			args[0], args[1] = mutateFloat(args[0].(*big.Int), args[1].(*big.Int), rate)

			if res, err := encodeHelper(uniswap_pair_abi, input[:4], args); err == nil {
				return res, true
			}
		}
	}
	return nil, false
}
