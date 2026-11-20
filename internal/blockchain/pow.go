package blockchain

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
)

type Target [32]byte

const (
	TargetBlockTime  int64 = 5
	RetargetInterval       = 50
	AdjustmentFactor int64 = 4
)

func InitialTarget() Target {
	var target Target

	for i := 2; i < len(target); i++ {
		target[i] = 0xff
	}

	return target
}

func MaxTarget() Target {
	var target Target

	target[1] = 0x0f

	for i := 2; i < len(target); i++ {
		target[i] = 0xff
	}

	return target
}

func meetsTarget(hash Hash, target Target) bool {
	return bytes.Compare(hash[:], target[:]) <= 0
}

func mineBlock(data string, previousHash Hash, target Target, timestamp int64) (Block, error) {
	block := Block{
		Data: data,
		Header: BlockHeader{
			Version:      blockVersion,
			PreviousHash: previousHash,
			DataHash:     hashBytes([]byte(data)),
			Target:       target,
			Timestamp:    timestamp,
		},
	}

	for {
		hash := block.Header.Hash()

		if meetsTarget(hash, target) {
			return block, nil
		}

		if block.Header.Nonce == math.MaxUint64 {
			return Block{}, fmt.Errorf("nonce space exhaustedz")
		}

		block.Header.Nonce++
	}
}

func calculateNextTarget(oldTarget Target, actualSeconds int64) Target {
	expectedSeconds := int64(RetargetInterval) * TargetBlockTime
	oldValue := new(big.Int).SetBytes(oldTarget[:])
	factor := big.NewInt(AdjustmentFactor)

	newTarget := new(big.Int).Mul(oldValue, big.NewInt(actualSeconds))
	newTarget.Div(newTarget, big.NewInt(expectedSeconds))

	minValue := new(big.Int).Add(oldValue, big.NewInt(AdjustmentFactor-1))
	minValue.Div(minValue, factor)

	maxValue := new(big.Int).Mul(oldValue, factor)

	if newTarget.Cmp(minValue) < 0 {
		newTarget.Set(minValue)
	}
	if newTarget.Cmp(maxValue) > 0 {
		newTarget.Set(maxValue)
	}

	maxTarget := MaxTarget()
	powLimit := new(big.Int).SetBytes(maxTarget[:])
	if newTarget.Cmp(powLimit) > 0 {
		newTarget.Set(powLimit)
	}

	var target Target
	newTarget.FillBytes(target[:])

	return target
}
