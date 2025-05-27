// index.js
const express = require('express');
const cors = require('cors');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');
require('dotenv').config();

const app = express();
app.use(cors());
app.use(express.json());

// Fabric network config
const ccpPath = path.resolve(__dirname, 'connection-org1.json');
const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

// Main helper: getContract
async function getContract() {
    const walletPath = path.join(__dirname);
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: 'Admin@org1.example.com',
        discovery: { enabled: true, asLocalhost: false },
    });

    const network = await gateway.getNetwork('mychannel');
    return network.getContract('wallet');
}

// Create wallet
app.post('/wallet', async (req, res) => {
    const { owner, balance } = req.body;
    try {
        const contract = await getContract();
        await contract.submitTransaction(
            'CreateWallet',
            owner,
            balance.toString()
        );
        res.status(201).json({ message: 'Wallet created' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Read wallet
app.get('/wallet/:owner', async (req, res) => {
    try {
        const contract = await getContract();
        const result = await contract.evaluateTransaction(
            'ReadWallet',
            req.params.owner
        );
        res.json(JSON.parse(result.toString()));
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Reward for commit (backend use only)
app.post('/reward', async (req, res) => {
    const { owner, commitHash, amount } = req.body;
    try {
        const contract = await getContract();
        await contract.submitTransaction(
            'RewardForCommit',
            owner,
            commitHash,
            amount.toString()
        );
        res.json({ message: 'Reward issued' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Token transfer
app.post('/transfer', async (req, res) => {
    const { from, to, amount } = req.body;
    try {
        const contract = await getContract();
        await contract.submitTransaction(
            'Transfer',
            from,
            to,
            amount.toString()
        );
        res.json({ message: 'Transfer successful' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Burn tokens
app.post('/burn', async (req, res) => {
    const { owner, amount } = req.body;
    try {
        const contract = await getContract();
        await contract.submitTransaction('Burn', owner, amount.toString());
        res.json({ message: 'Tokens burned' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Get commit record
app.get('/commit/:hash', async (req, res) => {
    try {
        const contract = await getContract();
        const result = await contract.evaluateTransaction(
            'GetCommitRecord',
            req.params.hash
        );
        res.json(JSON.parse(result.toString()));
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// 라우트 설정
const walletRoutes = require('./routes/walletRoutes');
app.use('/api/wallet', walletRoutes);

// 에러 핸들링 미들웨어
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({
        error: '서버 에러가 발생했습니다.',
        message: err.message
    });
});

// 서버 시작
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`🚀 지갑 API 서버가 포트 ${PORT}에서 실행 중입니다.`);
});
