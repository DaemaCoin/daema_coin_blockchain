package main

import (
    "encoding/json"
    "fmt"
    "strconv"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TokenChaincode는 코인 관련 기능을 구현하는 체인코드입니다.
type TokenChaincode struct {
    contractapi.Contract
}

// 계정 구조체 정의
type Account struct {
    Balance int `json:"balance"`
}

// Mint(코인 발행) - 관리자만 실행 가능
func (t *TokenChaincode) Mint(ctx contractapi.TransactionContextInterface, account string, amount int) error {
    clientMSP, err := ctx.GetClientIdentity().GetMSPID()
    if err != nil {
        return fmt.Errorf("failed to get client MSP: %v", err)
    }

    // 관리자만 Mint 실행 가능
    if clientMSP != "Org1MSP" {
        return fmt.Errorf("only Org1 admin can mint tokens")
    }

    accBytes, _ := ctx.GetStub().GetState(account)
    acc := Account{Balance: 0}
    if accBytes != nil {
        json.Unmarshal(accBytes, &acc)
    }

    acc.Balance += amount
    accBytes, _ = json.Marshal(acc)
    ctx.GetStub().PutState(account, accBytes)

    return nil
}

// Transfer(코인 전송)
func (t *TokenChaincode) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, amount int) error {
    fromBytes, _ := ctx.GetStub().GetState(from)
    toBytes, _ := ctx.GetStub().GetState(to)

    if fromBytes == nil {
        return fmt.Errorf("account %s not found", from)
    }

    fromAcc := Account{}
    json.Unmarshal(fromBytes, &fromAcc)

    if fromAcc.Balance < amount {
        return fmt.Errorf("insufficient balance")
    }

    toAcc := Account{Balance: 0}
    if toBytes != nil {
        json.Unmarshal(toBytes, &toAcc)
    }

    fromAcc.Balance -= amount
    toAcc.Balance += amount

    fromBytes, _ = json.Marshal(fromAcc)
    toBytes, _ = json.Marshal(toAcc)

    ctx.GetStub().PutState(from, fromBytes)
    ctx.GetStub().PutState(to, toBytes)

    return nil
}

// BalanceOf(잔액 조회)
func (t *TokenChaincode) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {
    accBytes, _ := ctx.GetStub().GetState(account)
    if accBytes == nil {
        return 0, fmt.Errorf("account %s not found", account)
    }

    acc := Account{}
    json.Unmarshal(accBytes, &acc)
    return acc.Balance, nil
}

// main 함수
func main() {
    chaincode, err := contractapi.NewChaincode(new(TokenChaincode))
    if err != nil {
        fmt.Printf("Error creating token chaincode: %s", err)
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting token chaincode: %s", err)
    }
}
