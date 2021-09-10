/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/kingstenzzz/statechannel/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	tokenChaincode, err := contractapi.NewChaincode(&chaincode.StateChannel{})
	if err != nil {
		log.Panicf("Error creating token-erc-20 chaincode: %v", err)
	}

	if err := tokenChaincode.Start(); err != nil {
		log.Panicf("Error starting token-erc-20 chaincode: %v", err)
	}
}
