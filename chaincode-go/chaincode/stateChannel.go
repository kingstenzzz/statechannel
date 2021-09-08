package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

// SmartContract provides functions for transferring tokens between accounts
type StateChannel struct {
	token TokenInterface
	contractapi.Contract
}

type token struct {
}

var channelCounter uint //channelCount
var mapAddress = make(map[string]uint32)

// event provides an organized struct for emitting events
type Player struct {
	uid        uint   `json:"uid"`
	addr       string `json:"from"`
	credit     uint   `json:"credit"`
	withdrawal uint   `json:"withdrawal"`
	withdrawn  uint   `json:"withdrawn"`
	deposit    uint   `json:"deposit"`
}

type Status uint8

const (
	Init    Status = 0
	OK      Status = 1
	Pending Status = 2
)

var channels []Channel

type Payment struct {
	amount    uint
	expiry    uint
	recipient string
	//preimageHash
}

type Channel struct {
	tokenAddress string
	playercount  uint32
	bestRound    int
	status       Status
	deadline     uint
}

/*
event
*/

// Create adds a new key with value to the world state
func (sc *StateChannel) CreateChannel(ctx contractapi.TransactionContextInterface, channelId string) (uint, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id: %v", err)
	}
	existing, err := ctx.GetStub().GetState(channelId)
	var channel Channel
	var payment Payment
	if err != nil {
		return 0, errors.New("Unable to interact with world state\n")
	}
	if existing != nil {
		return 0, fmt.Errorf("Cannot create world state pair with key %s. Already exists\n", channelId)
	}

	chId := channelCounter
	channel.tokenAddress = clientID

	channel.bestRound = -1
	channel.status = Init
	channel.deadline = 0
	channel.playercount = 2

	payment.expiry = 0
	payment.amount = 0
	payment.recipient = ""
	channels[chId] = channel
	channelCounter += 1

	UpdateChannel(ctx, channelId, channel)
	return channelCounter, nil
}
func (sc *StateChannel) JoinChannel(ctx contractapi.TransactionContextInterface, other string, channelName string, channelAdd string, amount uint) (uint, error) {

}

func (sc *StateChannel) CreateWithDeposit(ctx contractapi.TransactionContextInterface, other string, channelName string, channelAdd string, amount uint) (uint, error) {
	chId, err := sc.CreateChannel(ctx, other, channelName, channelAdd)
	if err != nil {
		return 0, errors.New("Create channel failed\n")
	}
	err = sc.deposit(ctx, chId, amount)
	if err != nil {
		return 0, err
	}

	return 0, err
}
func (sc *StateChannel) deposit(ctx contractapi.TransactionContextInterface, chId uint, amount uint) error {
	channel := channels[chId]
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	err = sc.token.TransferFrom(ctx, clientID, channel.tokenAddress, int(amount))
	if err != nil {
		return err
	}
	return nil
}

//return both playerID
func (sc *StateChannel) GetPlayers(ctx contractapi.TransactionContextInterface, channelName string) (error, string) {
	channelByte, err := ctx.GetStub().GetState(channelName)
	if err != nil {
		return errors.New("Create channel failed\n"), ""
	}
	resultsIterator, err := stub.GetStateByRange("Player:1", "Player:3")

	players := new(Channel)
	err = json.Unmarshal(channelByte, &channel)
	if err != nil {
		return errors.New("Unmarshal json faild"), ""
	}
	Players := ""
	for _, player := range platers {
		Players += player.addr
		if len(player.addr) == 0 {
			break
		}
	}
	return nil, Players
}

func (sc *StateChannel) ReadAsset(ctx contractapi.TransactionContextInterface, chId string) (*Channel, error) {
	assetJSON, err := ctx.GetStub().GetState(chId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", chId)
	}

	var channel Channel
	err = json.Unmarshal(assetJSON, &channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (sc *StateChannel) SendTokenTo(ctx contractapi.TransactionContextInterface, from, to, channelName string, amount uint) error {
	channelByte, err := ctx.GetStub().GetState(channelName)
	if err != nil {
		return errors.New("Get channel failed\n")
	}
	channel := new(Channel)
	err = json.Unmarshal(channelByte, &channel)

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	if strings.Compare(clientID, channel.platers[mapAddress[from]].addr) == 0 || strings.Compare(clientID, channel.platers[mapAddress[from]].addr) == 0 {
		return errors.New("error client address")
	}

	for i, player := range channel.platers {
		if strings.Compare(player.addr, clientID) == 0 {
			break
		}
		if i == len(channel.platers) {
			return errors.New("not in player")
		}
	}
	if (channel.platers[mapAddress[from]].credit)-amount <= 0 {
		return errors.New("unenough credit")
	}
	channel.platers[mapAddress[from]].credit -= 10
	channel.platers[mapAddress[to]].credit += 10
	///event

	return nil
}
