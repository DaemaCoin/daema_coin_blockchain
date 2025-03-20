package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type HelloWorldChaincode struct {
	contractapi.Contract
}

func (c *HelloWorldChaincode) HelloWorld(ctx contractapi.TransactionContextInterface) (string, error) {
	return "Hello, World!", nil
}

func main() {
	cc, err := contractapi.NewChaincode(&HelloWorldChaincode{})
	if err != nil {
		fmt.Printf("체인코드 생성 실패: %s", err.Error())
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("체인코드 시작 실패: %s", err.Error())
	}
}
