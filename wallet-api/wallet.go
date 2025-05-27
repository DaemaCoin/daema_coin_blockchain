package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// WalletContract 스마트컨트랙트 정의
type WalletContract struct {
	contractapi.Contract
}

// Wallet 지갑 구조체
type Wallet struct {
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

// CommitRecord 커밋 보상 기록
type CommitRecord struct {
	Owner      string    `json:"owner"`
	CommitHash string    `json:"commitHash"`
	Reward     float64   `json:"reward"`
	Timestamp  time.Time `json:"timestamp"`
}

// RewardEvent 이벤트 구조체
type RewardEvent struct {
	Owner      string  `json:"owner"`
	Amount     float64 `json:"amount"`
	CommitHash string  `json:"commitHash"`
}

// InitLedger 초기화
func (wc *WalletContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	wallet := Wallet{Owner: "user1", Balance: 100.0}
	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("wallet_user1", data)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("AUTHORIZED_SERVER_MSPID", []byte("ServerMSP"))
}

// CreateWallet 지갑 생성
func (wc *WalletContract) CreateWallet(ctx contractapi.TransactionContextInterface, owner string, initialBalance float64) error {
	if initialBalance < 0 {
		return fmt.Errorf("initial balance cannot be negative")
	}
	existing, err := ctx.GetStub().GetState("wallet_" + owner)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("wallet already exists for %s", owner)
	}
	wallet := Wallet{Owner: owner, Balance: initialBalance}
	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("wallet_"+owner, data)
}

// ReadWallet 지갑 조회
func (wc *WalletContract) ReadWallet(ctx contractapi.TransactionContextInterface, owner string) (*Wallet, error) {
	data, err := ctx.GetStub().GetState("wallet_" + owner)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("wallet not found for %s", owner)
	}
	var wallet Wallet
	err = json.Unmarshal(data, &wallet)
	return &wallet, err
}

// RewardForCommit 커밋 보상 지급
func (wc *WalletContract) RewardForCommit(ctx contractapi.TransactionContextInterface, owner string, commitHash string, amount float64) error {
	if err := wc.checkAuthorizedServer(ctx); err != nil {
		return err
	}
	if amount <= 0 || owner == "" || commitHash == "" {
		return fmt.Errorf("invalid input")
	}

	key := "commit_" + commitHash
	existing, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("commit already rewarded")
	}

	// 지갑 조회 및 업데이트
	wallet, err := wc.ReadWallet(ctx, owner)
	if err != nil {
		return err
	}
	wallet.Balance += amount
	updated, _ := json.Marshal(wallet)
	err = ctx.GetStub().PutState("wallet_"+owner, updated)
	if err != nil {
		return err
	}

	// 커밋 기록 저장
	record := CommitRecord{Owner: owner, CommitHash: commitHash, Reward: amount, Timestamp: time.Now()}
	rdata, _ := json.Marshal(record)
	err = ctx.GetStub().PutState(key, rdata)
	if err != nil {
		return err
	}

	// 이벤트 발생
	event := RewardEvent{Owner: owner, Amount: amount, CommitHash: commitHash}
	eventJSON, _ := json.Marshal(event)
	return ctx.GetStub().SetEvent("RewardEvent", eventJSON)
}

// Transfer 토큰 전송
func (wc *WalletContract) Transfer(ctx contractapi.TransactionContextInterface, from, to string, amount float64) error {
	if from == to {
		return fmt.Errorf("cannot transfer to same account")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	sender, err := wc.ReadWallet(ctx, from)
	if err != nil {
		return err
	}
	receiver, err := wc.ReadWallet(ctx, to)
	if err != nil {
		return err
	}
	if sender.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	sender.Balance -= amount
	receiver.Balance += amount

	sdata, _ := json.Marshal(sender)
	rdata, _ := json.Marshal(receiver)
	err = ctx.GetStub().PutState("wallet_"+from, sdata)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("wallet_"+to, rdata)
}

// Burn 토큰 차감
func (wc *WalletContract) Burn(ctx contractapi.TransactionContextInterface, owner string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	wallet, err := wc.ReadWallet(ctx, owner)
	if err != nil {
		return err
	}
	if wallet.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	wallet.Balance -= amount
	data, _ := json.Marshal(wallet)
	return ctx.GetStub().PutState("wallet_"+owner, data)
}

// MyBalance 현재 사용자의 잔액 확인
func (wc *WalletContract) MyBalance(ctx contractapi.TransactionContextInterface) (float64, error) {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, err
	}
	wallet, err := wc.ReadWallet(ctx, id)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

// GetCommitRecord 커밋 보상 내역 조회
func (wc *WalletContract) GetCommitRecord(ctx contractapi.TransactionContextInterface, commitHash string) (*CommitRecord, error) {
	data, err := ctx.GetStub().GetState("commit_" + commitHash)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("record not found")
	}
	var record CommitRecord
	err = json.Unmarshal(data, &record)
	return &record, err
}

// checkAuthorizedServer 호출자 권한 확인
func (wc *WalletContract) checkAuthorizedServer(ctx contractapi.TransactionContextInterface) error {
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	expected, err := ctx.GetStub().GetState("AUTHORIZED_SERVER_MSPID")
	if err != nil {
		return err
	}
	if mspid != string(expected) {
		return fmt.Errorf("unauthorized MSP: %s", mspid)
	}
	return nil
}

// main 진입점
func main() {
	cc, err := contractapi.NewChaincode(new(WalletContract))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
