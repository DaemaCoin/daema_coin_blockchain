package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PingContract : 테스트용 체인코드
type PingContract struct {
	contractapi.Contract
}

// Ping : 호출 시 "pong" 반환
func (pc *PingContract) Ping(ctx contractapi.TransactionContextInterface) (string, error) {
	return "pong", nil
}

func main() {
	cc, err := contractapi.NewChaincode(new(PingContract))
	if err != nil {
		panic("Could not create chaincode: " + err.Error())
	}
	if err := cc.Start(); err != nil {
		panic("Failed to start chaincode: " + err.Error())
	}
}
