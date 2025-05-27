const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
require('dotenv').config();

async function getContract() {
    try {
        // connection profile 로드
        const ccpPath = path.resolve(__dirname, '../connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // 인증서 및 키 로드
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Gateway 연결
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: process.env.FABRIC_IDENTITY,
            discovery: { enabled: true, asLocalhost: false }
        });

        // 네트워크 및 컨트랙트 가져오기
        const network = await gateway.getNetwork(process.env.CHANNEL_NAME);
        const contract = network.getContract(process.env.CHAINCODE_NAME);

        return { gateway, contract };
    } catch (error) {
        console.error(`Error in getContract: ${error}`);
        throw error;
    }
}

module.exports = { getContract }; 