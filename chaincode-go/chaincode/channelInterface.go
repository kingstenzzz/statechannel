package chaincode

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ChannelInterface interface {
	AddChannel(ChannelInterface)
}

func UpdateChannel(ctx contractapi.TransactionContextInterface, chId string, channel Channel) error {
	channelByte, err := json.Marshal(channel)
	if err != nil {
		return errors.New("Marshal channel struct failed\n")
	}
	err = ctx.GetStub().PutState(chId, channelByte)

	if err != nil {
		return errors.New("Unable to interact with world state\n")
	}

	return nil
}

func UpdatePlayer(ctx contractapi.TransactionContextInterface, uId uint32, player Player) error {
	return nil
}
func UpdatePayment(ctx contractapi.TransactionContextInterface, payId uint32, payment Payment) error {
	return nil
}
