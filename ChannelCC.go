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
		log.Panicf("Error creating sattechanel chaincode: %v", err)
	}

	if err := tokenChaincode.Start(); err != nil {
		log.Panicf("Error starting sattechanel chaincode: %v", err)
	}
}
