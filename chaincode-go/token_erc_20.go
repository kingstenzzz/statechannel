/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/hyperledger/fabric-samples/token-erc-20/chaincode-go/chaincode"
	//"github.com/kingstenzzz/statechannel/chaincode-go/chaincode"
	"./chaincode"


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
