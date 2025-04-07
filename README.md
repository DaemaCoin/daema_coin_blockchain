# 학교 블록체인 네트워크

이 프로젝트는 하이퍼레저 패브릭을 사용하여 학교 내부의 프라이빗 블록체인 네트워크를 구현한 것입니다. GitHub 커밋을 기반으로 토큰을 발급하고, 사용자 간 토큰 전송이 가능합니다.

## 주요 기능

1. GitHub ID 기반 지갑 생성
2. GitHub 커밋 검증 및 토큰 발급 (Ollama를 통한 검증)
3. 사용자 간 토큰 전송

## 설치 방법

1. 하이퍼레저 패브릭 네트워크 설정:

```bash
cd test-network
./network.sh up createChannel -c mychannel -ca
```

2. 체인코드 배포:

```bash
cd scripts
./deployCC.sh
```

3. Ollama 설치 및 실행:

```bash
# Ollama 설치 (Ubuntu)
curl https://ollama.ai/install.sh | sh

# Ollama 실행
ollama serve
```

## API 엔드포인트

### 1. 지갑 생성

```bash
POST /wallets
{
    "githubId": "your_github_id"
}
```

### 2. 모든 지갑 조회

```bash
GET /wallets
```

### 3. 특정 지갑 조회

```bash
GET /wallets/:githubId
```

### 4. 커밋 검증 및 토큰 발급

```bash
POST /wallets/:githubId/commits
{
    "commitData": {
        "files": [
            {
                "sha": "b64f6e8",
                "filename": "README.md",
                "status": "modified",
                "additions": 1,
                "deletions": 0,
                "changes": 1,
                "patch": "@@ -1 +1,2 @@\n # 내가 공부한거 올리는 레포\n+# this"
            }
        ]
    }
}
```

### 5. 토큰 전송

```bash
POST /wallets/:githubId/transfers
{
    "toGithubId": "receiver_github_id",
    "amount": 100
}
```

## 체인코드 기능

### 1. 지갑 관리

- `CreateWallet`: GitHub ID로 지갑 생성
- `GetWallet`: 특정 지갑 정보 조회
- `GetAllWallets`: 모든 지갑 정보 조회
- `WalletExists`: 지갑 존재 여부 확인

### 2. 커밋 검증 및 토큰 발급

- `ValidateAndRewardCommit`: 커밋 검증 및 토큰 발급
- `validateCommitWithOllama`: Ollama를 사용한 커밋 검증
  - 커밋 메시지 명확성
  - 변경사항 적절성
  - 보안 문제 검사
  - 코드 품질 검사

### 3. 토큰 전송

- `Transfer`: 지갑 간 토큰 전송
  - 잔액 확인
  - 송신자 잔액 차감
  - 수신자 잔액 증가

## 보안 주의사항

1. GitHub 토큰은 필요한 최소한의 권한만 부여하세요.
2. 실제 운영 환경에서는 HTTPS를 사용하세요.
3. 인증서와 키 파일은 안전하게 보관하세요.
4. Ollama API는 로컬에서만 접근 가능하도록 설정하세요.
