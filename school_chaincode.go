package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing the school token
type SmartContract struct {
	contractapi.Contract
}

// Wallet represents a user's wallet
type Wallet struct {
	GithubID string `json:"githubID"`
	Balance  int    `json:"balance"`
}

// CommitFile represents a file in a GitHub commit
type CommitFile struct {
	SHA        string `json:"sha"`
	Filename   string `json:"filename"`
	Status     string `json:"status"`
	Additions  int    `json:"additions"`
	Deletions  int    `json:"deletions"`
	Changes    int    `json:"changes"`
	BlobURL    string `json:"blob_url"`
	RawURL     string `json:"raw_url"`
	ContentsURL string `json:"contents_url"`
	Patch      string `json:"patch"`
}

// Commit represents a GitHub commit
type Commit struct {
	Files []CommitFile `json:"files"`
}

// InitLedger adds a base set of wallets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	wallets := []Wallet{
		{GithubID: "admin", Balance: 1000000},
	}

	for _, wallet := range wallets {
		walletJSON, err := json.Marshal(wallet)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(wallet.GithubID, walletJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateWallet creates a new wallet for a github user
func (s *SmartContract) CreateWallet(ctx contractapi.TransactionContextInterface, githubID string) error {
	exists, err := s.WalletExists(ctx, githubID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the wallet for githubID %s already exists", githubID)
	}

	wallet := Wallet{
		GithubID: githubID,
		Balance:  0,
	}
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(githubID, walletJSON)
}

// WalletExists returns true when wallet with given ID exists in world state
func (s *SmartContract) WalletExists(ctx contractapi.TransactionContextInterface, githubID string) (bool, error) {
	walletJSON, err := ctx.GetStub().GetState(githubID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return walletJSON != nil, nil
}

// GetWallet returns the wallet stored in the world state with given id
func (s *SmartContract) GetWallet(ctx contractapi.TransactionContextInterface, githubID string) (*Wallet, error) {
	walletJSON, err := ctx.GetStub().GetState(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if walletJSON == nil {
		return nil, fmt.Errorf("the wallet %s does not exist", githubID)
	}

	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// validateCommitWithOllama validates a commit using Ollama
func (s *SmartContract) validateCommitWithOllama(commit Commit) (bool, error) {
	// 커밋 정보를 문자열로 변환
	commitInfo := fmt.Sprintf("Commit files:\n")
	for _, file := range commit.Files {
		commitInfo += fmt.Sprintf("- File: %s\n  Status: %s\n  Changes: +%d -%d\n  Patch: %s\n",
			file.Filename, file.Status, file.Additions, file.Deletions, file.Patch)
	}

	// Ollama API 호출
	cmd := exec.Command("curl", "-X", "POST", "http://localhost:11434/api/chat",
		"-H", "Content-Type: application/json",
		"-d", fmt.Sprintf(`{
			"model": "llama2",
			"messages": [
				{
					"role": "user",
					"content": "다음 커밋을 검증해주세요:\n%s\n\n다음 사항들을 검토해주세요:\n1. 커밋 메시지가 명확하고 설명적인가?\n2. 변경사항이 적절한가?\n3. 보안상의 문제는 없는가?\n4. 코드 품질은 적절한가?\n\n답변은 '적절' 또는 '부적절'로만 해주세요.",
					commitInfo
				}
			]
		}`))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to call Ollama: %v", err)
	}

	// Ollama 응답 파싱
	var response struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return false, fmt.Errorf("failed to parse Ollama response: %v", err)
	}

	// 응답에 '적절'이 포함되어 있으면 true 반환
	return strings.Contains(response.Message.Content, "적절"), nil
}

// ValidateAndRewardCommit validates a commit and rewards tokens if valid
func (s *SmartContract) ValidateAndRewardCommit(ctx contractapi.TransactionContextInterface, githubID string, commitJSON string) error {
	// 커밋 정보 파싱
	var commit Commit
	if err := json.Unmarshal([]byte(commitJSON), &commit); err != nil {
		return fmt.Errorf("failed to parse commit: %v", err)
	}

	// Ollama를 사용하여 커밋 검증
	isValid, err := s.validateCommitWithOllama(commit)
	if err != nil {
		return fmt.Errorf("failed to validate commit: %v", err)
	}

	if !isValid {
		return fmt.Errorf("commit validation failed")
	}

	// 검증이 성공하면 토큰 발급
	wallet, err := s.GetWallet(ctx, githubID)
	if err != nil {
		return err
	}

	wallet.Balance += 100 // 100 토큰 발급
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(githubID, walletJSON)
}

// Transfer transfers tokens from one wallet to another
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, fromGithubID string, toGithubID string, amount int) error {
	fromWallet, err := s.GetWallet(ctx, fromGithubID)
	if err != nil {
		return err
	}

	toWallet, err := s.GetWallet(ctx, toGithubID)
	if err != nil {
		return err
	}

	if fromWallet.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	fromWalletJSON, err := json.Marshal(fromWallet)
	if err != nil {
		return err
	}

	toWalletJSON, err := json.Marshal(toWallet)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(fromGithubID, fromWalletJSON)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(toGithubID, toWalletJSON)
}

// GetAllWallets returns all wallets found in world state
func (s *SmartContract) GetAllWallets(ctx contractapi.TransactionContextInterface) ([]*Wallet, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var wallets []*Wallet
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var wallet Wallet
		err = json.Unmarshal(queryResponse.Value, &wallet)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating school chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting school chaincode: %s", err.Error())
	}
} 