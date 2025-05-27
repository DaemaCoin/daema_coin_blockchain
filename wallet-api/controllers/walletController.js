const { getContract } = require('../fabric/connection');

exports.createWallet = async (req, res) => {
    const { owner, initialBalance } = req.body;
    try {
        const { contract, gateway } = await getContract();
        await contract.submitTransaction(
            'CreateWallet',
            owner,
            String(initialBalance)
        );
        gateway.disconnect();
        res.status(201).json({ message: '지갑이 생성되었습니다.' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

exports.readWallet = async (req, res) => {
    try {
        const { contract, gateway } = await getContract();
        const result = await contract.evaluateTransaction(
            'ReadWallet',
            req.params.owner
        );
        gateway.disconnect();
        res.json(JSON.parse(result.toString()));
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

exports.rewardForCommit = async (req, res) => {
    const { owner, commitHash, amount } = req.body;
    try {
        const { contract, gateway } = await getContract();
        await contract.submitTransaction(
            'RewardForCommit',
            owner,
            commitHash,
            String(amount)
        );
        gateway.disconnect();
        res.json({ message: '보상이 지급되었습니다.' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

exports.transfer = async (req, res) => {
    const { from, to, amount } = req.body;
    try {
        const { contract, gateway } = await getContract();
        await contract.submitTransaction(
            'Transfer',
            from,
            to,
            String(amount)
        );
        gateway.disconnect();
        res.json({ message: '전송이 완료되었습니다.' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

exports.burn = async (req, res) => {
    const { owner, amount } = req.body;
    try {
        const { contract, gateway } = await getContract();
        await contract.submitTransaction(
            'Burn',
            owner,
            String(amount)
        );
        gateway.disconnect();
        res.json({ message: '토큰이 소각되었습니다.' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

exports.getCommitRecord = async (req, res) => {
    try {
        const { contract, gateway } = await getContract();
        const result = await contract.evaluateTransaction(
            'GetCommitRecord',
            req.params.hash
        );
        gateway.disconnect();
        res.json(JSON.parse(result.toString()));
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};

exports.getMyBalance = async (req, res) => {
    try {
        const { contract, gateway } = await getContract();
        const result = await contract.evaluateTransaction('MyBalance');
        gateway.disconnect();
        res.json({ balance: parseFloat(result.toString()) });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
};
