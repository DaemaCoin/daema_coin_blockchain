# 학교 블록체인 네트워크

이 프로젝트는 하이퍼레저 패브릭을 사용하여 학교 내부의 프라이빗 블록체인 네트워크를 구현한 것입니다. GitHub 커밋을 기반으로 토큰을 발급하고, 사용자 간 토큰 전송이 가능합니다.

## 주요 기능

1. GitHub ID 기반 지갑 생성
2. GitHub 커밋 검증 및 토큰 발급
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

3. 서버 실행:
```bash
cd server
go run server.go
```

## API 엔드포인트

### 1. 커밋 검증 및 토큰 발급
```bash
POST /validate-commit
{
    "githubID": "your_github_id",
    "repoName": "owner/repo",
    "commitSHA": "commit_sha"
}
```

### 2. 토큰 전송
```bash
POST /transfer
{
    "fromGithubID": "sender_github_id",
    "toGithubID": "receiver_github_id",
    "amount": 100
}
```

## 체인코드 기능

1. `CreateWallet`: GitHub ID로 지갑 생성
2. `RewardCommit`: 커밋 검증 후 토큰 발급
3. `Transfer`: 토큰 전송
4. `GetWallet`: 지갑 정보 조회
5. `GetAllWallets`: 모든 지갑 정보 조회

## 보안 주의사항

1. GitHub 토큰은 필요한 최소한의 권한만 부여하세요.
2. 실제 운영 환경에서는 HTTPS를 사용하세요.
3. 인증서와 키 파일은 안전하게 보관하세요. 