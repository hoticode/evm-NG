// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
package evm

import (
	"testing"

	"encoding/hex"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/blockchain/config"
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"math/big"
)

var (
	callerAddress   = types.HexToAddress("0x8a8c58e424f4a6d2f0b2270860c96dfe34f10c78")
	contractAddress = types.HexToAddress("0xf74cc8824a00bcb96e8546bf3b4dc47ace9cab2c")
	code, _         = hex.DecodeString("6080604052348015600f57600080fd5b5060998061001e6000396000f300608060405260043610603e5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416634f2be91f81146043575b600080fd5b348015604e57600080fd5b5060556067565b60408051918252519081900360200190f35b610378905600a165627a7a723058205d540f3e87376532c076a230eb73eee4aa46c0df1a71cdba5a33cda64a8e6f400029")
	input1, _       = hex.DecodeString("4f2be91f")
	input2, _       = hex.DecodeString("4f2be91f")
)

// mock a blockchain.
func mockPreBlockChain() *blockchain.BlockChain {
	// init chain
	blockchain.InitBlockChain(config.BlockChainConfig{
		PluginName:    blockchain.PLUGIN_MEMDB,
		StateDataPath: "",
		BlockDataPath: "",
	})
	// create chain instance
	bc, _ := blockchain.NewLatestStateBlockChain()

	//create caller account
	bc.CreateAccount(callerAddress)
	bc.AddBalance(callerAddress, big.NewInt(1000))

	//create contract account
	bc.CreateAccount(contractAddress)
	bc.SetCode(contractAddress, code)
	return bc
}

// mock a evm instance
func mockEVM(bc *blockchain.BlockChain) *EVM {
	msg := Message{
		from:     callerAddress,
		gasPrice: big.NewInt(0x5af3107a4000),
	}
	header := &types.Header{
		PrevBlockHash: types.HexToHash(""),
		Height:        1,
		Timestamp:     1,
	}
	author := types.HexToAddress("0x0000000000000000000000000000000000000000")
	context := NewEVMContext(msg, header, bc, author)
	return NewEVM(context, bc)
}

// test execute contract
func TestVM(t *testing.T) {
	assert := assert.New(t)
	//init statedb state
	bc := mockPreBlockChain()

	//execute contract code
	evmInst := mockEVM(bc)

	//specify the caller address
	callerRef := AccountRef(callerAddress)
	_, _, error := evmInst.Call(callerRef, contractAddress, input1, 3000, big.NewInt(0))
	assert.Nil(error)
}
