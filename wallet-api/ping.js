// ping.js
const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const cors = require('cors');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 3000;

// CORS 허용
app.use(cors());
app.use(express.json());

// 체인코드 호출 함수
async function callChaincode() {
  try {
    const ccpPath = path.resolve(__dirname, 'connection-org1.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

    const walletPath = path.join(__dirname, 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const identity = await wallet.get('admin');
    if (!identity) {
      throw new Error('❌ Admin identity not found in wallet. Please run enrollAdmin.js first');
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, {
      wallet,
      identity: 'admin',
      discovery: { enabled: true, asLocalhost: true },
      clientTlsIdentity: 'admin',
      eventHandlerOptions: {
        commitTimeout: 300,
        strategy: null
      }
    });

    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('wallet');

    // 지갑 생성 테스트
    const testOwner = 'test_user_' + Date.now();
    await contract.submitTransaction('CreateWallet', testOwner, '100');
    
    // 지갑 조회 테스트
    const result = await contract.evaluateTransaction('ReadWallet', testOwner);

    await gateway.disconnect();

    return JSON.parse(result.toString());
  } catch (error) {
    console.error('⚠️ Failed to evaluate transaction:', error);
    return { error: error.message };
  }
}

// GET /
app.get('/', async (req, res) => {
  const response = await callChaincode();
  res.json(response);
});

// 서버 실행
app.listen(PORT, () => {
  console.log(`✅ Wallet API server is running at http://localhost:${PORT}`);
});
