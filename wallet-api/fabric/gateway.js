const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function getContract() {
    const ccpPath = path.resolve(__dirname, '..', 'connection-org1.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

    const walletPath = path.join(__dirname, '..', 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const identity = await wallet.get('Admin@org1.example.com');
    if (!identity) throw new Error('Admin identity not found in wallet');

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'Admin@org1.example.com',
        discovery: { enabled: true, asLocalhost: true },
    });

    const network = await gateway.getNetwork('mychannel');
    return network.getContract('wallet');
}

module.exports = getContract;
