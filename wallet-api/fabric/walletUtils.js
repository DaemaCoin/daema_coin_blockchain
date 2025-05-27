// fabric/walletUtils.js

const path = require('path');
const fs = require('fs');
const { Wallets, Gateway } = require('fabric-network');
const connectionProfile = require('./dynamicConnection'); // 동적으로 구성된 connection object
require('dotenv').config();

const walletPath = path.join(__dirname, '../wallet'); // 환경에 따라 조정 가능
const identityLabel = process.env.FABRIC_IDENTITY || 'Admin@org1.example.com';
const channelName = process.env.CHANNEL_NAME || 'mychannel';
const chaincodeName = process.env.CHAINCODE_NAME || 'wallet';

let cachedContract = null;

/**
 * Fabric Gateway를 초기화하고 컨트랙트를 반환합니다.
 */
async function getContract() {
    if (cachedContract) return cachedContract;

    // 1. wallet 로딩
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    // 2. 인증서 존재 확인
    const identity = await wallet.get(identityLabel);
    if (!identity) {
        throw new Error(
            `Identity '${identityLabel}' not found in wallet at ${walletPath}`
        );
    }

    // 3. gateway 연결
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, {
        wallet,
        identity: identityLabel,
        discovery: { enabled: true, asLocalhost: true }, // true: local 환경, false: 서버 배포 시 변경
    });

    const network = await gateway.getNetwork(channelName);
    const contract = network.getContract(chaincodeName);

    cachedContract = contract;
    return contract;
}

module.exports = {
    getContract,
};
