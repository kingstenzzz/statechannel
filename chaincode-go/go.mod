module github.com/kingstenzzz/statechannel/chaincode-go

go 1.14

require (
	github.com/hyperledger/fabric-contract-api-go v1.1.0
	golang.org/x/tools v0.1.0 // indirect
	github.com/kingstenzzz/statechannel/chaincode-go/chaincode v1.0

)
replace (
github.com/kingstenzzz/statechannel/chaincode-go/chaincode v1.0 => "./chaincode"
)

