package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ChannelInterface interface {
	UpdateChannel(ctx contractapi.TransactionContextInterface, chId string, channel Channel)
	UpdatePlayer(ctx contractapi.TransactionContextInterface, uId int, player Player)
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

func UpdatePlayer(ctx contractapi.TransactionContextInterface, playerId int, player Player) error {
	playerKey := fmt.Sprintf("%s%d", "Player", playerId)
	playerByte, err := json.Marshal(player)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(playerKey, playerByte)
	if err != nil {
		return err
	}
	return nil
}
func UpdatePayment(ctx contractapi.TransactionContextInterface, payId int, payment Payment) error {
	return nil
}
