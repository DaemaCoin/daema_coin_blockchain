package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// WalletContract 정의
type WalletContract struct {
	contractapi.Contract
}

// Wallet 구조체
type Wallet struct {
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

// InitLedger 초기 데이터 등록
func (wc *WalletContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	wallet := Wallet{Owner: "user1", Balance: 100.0}
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("wallet_user1", walletJSON)
}

// CreateWallet 지갑 생성
func (wc *WalletContract) CreateWallet(ctx contractapi.TransactionContextInterface, owner string, initialBalance float64) error {
	wallet := Wallet{Owner: owner, Balance: initialBalance}
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("wallet_"+owner, walletJSON)
}

// ReadWallet 지갑 조회
func (wc *WalletContract) ReadWallet(ctx contractapi.TransactionContextInterface, owner string) (*Wallet, error) {
	walletJSON, err := ctx.GetStub().GetState("wallet_" + owner)
	if err != nil {
		return nil, err
	}
	if walletJSON == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// RewardToken: 커밋에 대한 토큰 지급
func (wc *WalletContract) RewardToken(ctx contractapi.TransactionContextInterface, owner string, amount float64) error {
	walletJSON, err := ctx.GetStub().GetState("wallet_" + owner)
	if err != nil {
		return err
	}
	if walletJSON == nil {
		return fmt.Errorf("wallet not found for owner %s", owner)
	}

	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	if err != nil {
		return err
	}

	wallet.Balance += amount

	updatedJSON, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("wallet_" + owner, updatedJSON)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(WalletContract))
	if err != nil {
		panic("Could not create chaincode: " + err.Error())
	}
	if err := chaincode.Start(); err != nil {
		panic("Failed to start chaincode: " + err.Error())
	}
}
