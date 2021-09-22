package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
	"time"
)

// SmartContract provides functions for transferring tokens between accounts
type StateChannel struct {
	//token TokenInterface
	channelInterface ChannelInterface
	tokenInterface   TokenInterface
	contractapi.Contract
}

var channelCounter int //channelCount

// event provides an organized struct for emitting events
type Player struct {
	Uid        int    `json:"Uid"`
	Addr       string `json:"Addr"`
	Credit     int    `json:"Credit"`
	Withdrawal int    `json:"Withdrawal"`
	Withdrawn  int    `json:"Withdrawn"`
	Deposit    int    `json:"Deposit"`
}

type Status int8

const (
	Init    Status = 0
	OK      Status = 1
	Pending Status = 2
	Close   Status = 3
	Dispute Status = 4
)

type Payment struct {
	Amount    int    `json:"amount"`
	Expiry    int    `json:"expiry"`
	Recipient string `json:"recipient"`
	//preimageHash
}

type Channel struct {
	ChannelAddress string              `json:"channel_address"`
	PlayerCount    int                 `json:"player_count"`
	BestRound      int                 `json:"best_round"`
	Status         Status              `json:"status"`
	OpenTime       timestamp.Timestamp `json:"open_time"`
	Deadline       int                 `json:"deadline"`
}

/*
event
*/

// Create adds a new key with value to the world state
func (sc *StateChannel) CreateChannel(ctx contractapi.TransactionContextInterface, chName string, deadLine int) (int, error) {
	fmt.Println("CreateChannel")
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
	channel.ChannelAddress = clientID

	channel.BestRound = -1
	channel.Status = Init
	channel.Deadline = 0
	channel.PlayerCount = 0

	payment.Expiry = 0
	payment.Amount = 0
	payment.Recipient = ""
	channelCounter += 1
	OpenTime, _ := ctx.GetStub().GetTxTimestamp()
	channel.OpenTime = *OpenTime
	channel.Deadline = deadLine

	var player []Player
	fmt.Printf("%+v\n", player)
	UpdateChannel(ctx, chName, channel)
	fmt.Println("create channel")
	return channelCounter, nil
}
func (sc *StateChannel) JoinChanel(ctx contractapi.TransactionContextInterface, chName string, deposit int) error {
	clientID, err := ctx.GetClientIdentity().GetID()
	channelByte, err := ctx.GetStub().GetState(chName)
	var channel Channel
	err = json.Unmarshal(channelByte, &channel)
	if err != nil {
		return err
	}
	channel.PlayerCount += 1
	var player Player
	player.Uid = channel.PlayerCount
	player.Addr = clientID
	player.Deposit = deposit
	player.Credit = deposit
	UpdatePlayer(ctx, player.Uid, player)
	UpdateChannel(ctx, chName, channel)
	return nil
}

func (sc *StateChannel) EixtChanel(ctx contractapi.TransactionContextInterface, chName, playerID string, deposit int) error {
	playerKey := "Player" + playerID
	playerByte, err := ctx.GetStub().GetState(playerKey)
	var player Player
	err = json.Unmarshal(playerByte, &player)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(playerKey)
	if err != nil {
		return err

	}
	channelCounter -= 1
	channelByte, err := ctx.GetStub().GetState(chName)
	var channel Channel
	err = json.Unmarshal(channelByte, &channel)
	if err != nil {
		return err
	}
	channelCounter -= 1
	UpdateChannel(ctx, chName, channel)
	return nil
}

//update a batch of player according to the status of offchain
func (sc *StateChannel) UpdateBatchStatus(ctx contractapi.TransactionContextInterface, chName string, playerGroup string) error {

	return nil

}

func (sc *StateChannel) deposit(ctx contractapi.TransactionContextInterface, chName string, amount int) error {
	channelByte, err := ctx.GetStub().GetState(chName)
	var channel Channel
	json.Unmarshal(channelByte, &channel)
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	err = sc.tokenInterface.TransferFrom(ctx, clientID, channel.ChannelAddress, amount)
	if err != nil {
		return err
	}
	return nil
}

//return both playerID
func (sc *StateChannel) GetPlayers(ctx contractapi.TransactionContextInterface, channelName string) (string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("Player1", "Player5")
	if err != nil {
		return "Fail to read plays", err
	}
	playersInfo := ""
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "Fail to read plays in Iterator ", err
		}
		player := new(Player)

		json.Unmarshal(queryResponse.Value, &player)
		playerInfo := fmt.Sprintf("player:%d--%s\r\n", player.Uid, player.Addr)
		playersInfo += playerInfo
	}

	return playersInfo, nil
}

func (sc *StateChannel) ReadAsset(ctx contractapi.TransactionContextInterface, chId string) (string, error) {
	assetJSON, err := ctx.GetStub().GetState(chId)
	if err != nil {
		return "nil", fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return "nil", fmt.Errorf("the asset %s does not exist", chId)
	}

	return string(assetJSON), nil
}

func (sc *StateChannel) CloseChannel(ctx contractapi.TransactionContextInterface, chId string) error {
	channelByte, err := ctx.GetStub().GetState(chId)
	if err != nil {
		return err
	}
	var channel Channel
	err = json.Unmarshal(channelByte, &channel)
	if channel.Status == Pending {
		return errors.New("It is pending\r\n")
	}
	if channel.Status == Dispute {
		return errors.New("Depositing\r\n")
	}
	if channel.BestRound > 10 {
		return errors.New("Beyond the Deadline\r\n")
	}
	now, err := ctx.GetStub().GetTxTimestamp()
	now1 := time.Now().Unix()
	fmt.Println(now, now1)
	fmt.Println(time.Unix(now1, 0))
	totalDeposit := 0
	totalWithdrawal := 0
	for id := 0; id <= (channel.PlayerCount); id++ {
		var player Player
		totalDeposit += player.Deposit
		player.Withdrawal = player.Credit
		totalWithdrawal += player.Withdrawal
		player.Deposit = 0
		player.Credit = 0
		player.Credit = player.Deposit
		player.Withdrawal = 0
		player.Uid = id
		UpdatePlayer(ctx, player.Uid, player)
	}
	if totalDeposit != totalWithdrawal {
		return errors.New("totalDeposit not equall totalWithdrawal")
	}
	channel.Status = Close
	//return token
	return nil
}

/*
func (sc *StateChannel) SendTokenTo(ctx contractapi.TransactionContextInterface, from, to, channelName string, Ammount uint) error {
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
	if strings.Compare(clientID, ) == 0 || strings.Compare(clientID, channel.platers[mapAddress[from]].Addr) == 0 {
		return errors.New("error client address")
	}

	for i, player := range channel.platers {
		if strings.Compare(player.Addr, clientID) == 0 {
			break
		}
		if i == len(channel.platers) {
			return errors.New("not in player")
		}
	}
	if (channel.platers[mapAddress[from]].Credit)-Ammount <= 0 {
		return errors.New("unenough Credit")
	}
	channel.platers[mapAddress[from]].Credit -= 10
	channel.platers[mapAddress[to]].Credit += 10
	///event

	return nil
}

*/

func Query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("指定的参数错误，必须且只能指定相应的Key")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据指定的" + args[0] + "查询数据时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的" + args[0] + "没有查询到相应的数据")
	}
	return shim.Success(result)
}
