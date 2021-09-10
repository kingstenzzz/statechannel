module github.com/kingstenzzz/statechannel

go 1.14

require (
	github.com/golang/protobuf v1.3.2
	github.com/hyperledger/fabric-contract-api-go v1.1.0
	github.com/kingstenzzz/statechannel v0.0.1 // indirect
	golang.org/x/tools v0.1.0 // indirect

)
replace (
		github.com/kingstenzzz/statechannel v0.0.1 => ./chaincode
)
