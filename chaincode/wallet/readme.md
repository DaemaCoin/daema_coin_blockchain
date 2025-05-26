# Wallet 체인코드 - 서버 연동 가이드

## 📀 개요

이 문서는 Hyperledger Fabric 기반의 "Wallet" 체인코드와 백엔드 서버(Node.js 등)가 어떻게 상호작용할 수 있는지 설명합니다. 지갑 생성, 잔액 조회, 커밋 기반 보상 지급, 토큰 전송, 트랜잭션 이력 조회 등을 다룹니다.

---

## 🌟 지원하는 체인코드 함수

| 함수 이름                                    | 설명                                      | 호출 주체                      |
| -------------------------------------------- | ----------------------------------------- | ------------------------------ |
| `InitLedger()`                               | 기본 지갑 초기화 및 권한 있는 MSP ID 등록 | 관리자만                       |
| `CreateWallet(owner, initialBalance)`        | 지정한 사용자의 지갑 생성                 | 백엔드 또는 사용자 클라이언트  |
| `ReadWallet(owner)`                          | 특정 사용자의 지갑 잔액 조회              | 백엔드 또는 사용자             |
| `MyBalance()`                                | 현재 호출자의 지갑 잔액 조회              | 사용자 클라이언트              |
| `RewardForCommit(owner, commitHash, amount)` | 커밋에 대한 토큰 보상 지급                | **백엔드만 가능 (인증된 MSP)** |
| `Transfer(from, to, amount)`                 | 한 지갑에서 다른 지갑으로 토큰 전송       | 지갑 소유자 (사용자)           |
| `GetCommitRecord(commitHash)`                | 보상된 커밋 기록 조회                     | 백엔드 또는 사용자             |

---

## 🚑 백엔드 서버가 수행해야 할 역할

### 1. GitHub Webhook 수신 핸들러 구축

- GitHub 커밋 푸시 이벤트를 수신할 Webhook 설정
- `commitHash`, `owner`(GitHub ID) 등 필요한 데이터 추출

### 2. Fabric 네트워크 연결

Fabric SDK 사용 예시 (Node.js):

```js
const { Gateway, Wallets } = require("fabric-network");
const ccp = require("./connection-org1.json");

const wallet = await Wallets.newFileSystemWallet("./wallet");
const gateway = new Gateway();
await gateway.connect(ccp, {
  wallet,
  identity: "Admin@org1.example.com",
  discovery: { enabled: true, asLocalhost: true },
});
const network = await gateway.getNetwork("mychannel");
const contract = network.getContract("wallet");
```

### 3. 체인코드 함수 호출 예시

#### 지갑 생성

```js
await contract.submitTransaction("CreateWallet", "githubUser123", "0");
```

#### 지갑 잔액 조회

```js
const result = await contract.evaluateTransaction("ReadWallet", "githubUser123");
console.log(JSON.parse(result.toString()));
```

#### 커밋 보상 지급 (서버만 호출 가능)

```js
await contract.submitTransaction("RewardForCommit", "githubUser123", "shaabc123...", "10");
```

#### 토큰 전송

```js
await contract.submitTransaction("Transfer", "githubUser123", "friend123", "5");
```

#### 본인 지갑 잔액 확인

```js
await contract.evaluateTransaction("MyBalance");
```

#### 커밋 보상 기록 조회

```js
await contract.evaluateTransaction("GetCommitRecord", "shaabc123...");
```

---

## 🔑 백엔드 연동에 필요한 파일

다음 파일들을 압축(zip)하여 백엔드 서버로 전달합니다:

```
wallet-zip/
├── connection-org1.json
├── signcerts/cert.pem
├── keystore/<private_key>.pem
└── ca.crt (TLS 검증용, 선택사항)
```

---

## 🔒 보안 고려사항

- 백엔드는 체인코드에 등록된 `AUTHORIZED_SERVER_MSPID`와 일치하는 MSP를 사용해야 함
- `RewardForCommit`은 서버 MSP만 호출 가능
- GitHub Webhook은 `X-Hub-Signature-256` 등으로 유효성 검증 필요
- 개인 키 및 인증서는 외부에 노출되지 않도록 보관

---

## 🚨 현재 제약사항

- 토큰 외부 출금 기능 없음
- 보상 토큰 만료일 없음
- 보상 정책은 코드에 하드코딩됨 (동적 아님)
