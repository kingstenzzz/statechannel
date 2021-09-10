package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"time"
)

// SmartContract provides functions for transferring tokens between accounts
type StateChannel struct {
	token TokenInterface
	contractapi.Contract
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
	Close   Status = 3
	Dispute Status = 4
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
	openTime     timestamp.Timestamp
	deadline     uint
}

/*
event
*/

// Create adds a new key with value to the world state
func (sc *StateChannel) CreateChannel(ctx contractapi.TransactionContextInterface, chName string, playerNum uint, players [][]string) (uint, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id: %v", err)
	}
	existing, err := ctx.GetStub().GetState(chName)
	var channel Channel
	var payment Payment
	if err != nil {
		return 0, errors.New("Unable to interact with world state\n")
	}
	if existing != nil {
		return 0, fmt.Errorf("Cannot create world state pair with key %s. Already exists\n", chName)
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
	//channel.openTime ,_ = ctx.GetStub().GetTxTimestamp()

	for id, playerObject := range players {
		var player Player
		player.addr = playerObject[id]
		deposit, _ := strconv.ParseUint(playerObject[1], 10, 32)
		player.deposit = uint(deposit)
		player.credit = player.deposit
		player.withdrawal = 0
		player.uid = uint(id)
		UpdatePlayer(ctx, player.uid, player)

	}
	UpdateChannel(ctx, chName, channel)
	return channelCounter, nil
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
	resultsIterator, err := ctx.GetStub().GetStateByRange("Player:1", "Player:3")
	defer resultsIterator.Close()
	if err != nil {
		return err, "Fail to read plays"
	}
	playersInfo := ""
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, ""
		}
		player := new(Player)
		json.Unmarshal(queryResponse.Value, &player)
		playerInfo := string(player.uid) + player.addr
		playersInfo += playerInfo
	}

	return nil, playersInfo
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

func (sc *StateChannel) CloseChannel(ctx contractapi.TransactionContextInterface, chId string) error {
	channelByte, err := ctx.GetStub().GetState(chId)
	if err != nil {
		return err
	}
	var channel Channel
	err = json.Unmarshal(channelByte, &channel)
	if channel.status == Pending {
		return errors.New("It is pending\r\n")
	}
	if channel.status == Dispute {
		return errors.New("Depositing\r\n")
	}
	if channel.bestRound > 10 {
		return errors.New("Beyond the deadline\r\n")
	}
	now, err := ctx.GetStub().GetTxTimestamp()
	now1 := time.Now().Unix()
	fmt.Println(now, now1)
	fmt.Println(time.Unix(now1, 0))
	var totalDeposit uint = 0
	var totalWithdrawal uint = 0
	for id := uint32(0); id <= (channel.playercount); id++ {
		var player Player
		totalDeposit += player.deposit
		player.withdrawal = player.credit
		totalWithdrawal += player.withdrawal
		player.deposit = 0
		player.credit = 0
		player.credit = player.deposit
		player.withdrawal = 0
		player.uid = uint(id)
		UpdatePlayer(ctx, player.uid, player)
	}
	if totalDeposit != totalWithdrawal {
		return errors.New("totalDeposit not equall totalWithdrawal")
	}
	channel.status = Close
	//return token
	return nil
}

/*
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
	if strings.Compare(clientID, ) == 0 || strings.Compare(clientID, channel.platers[mapAddress[from]].addr) == 0 {
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

*/
