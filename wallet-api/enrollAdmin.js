const { Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const fs = require('fs');
const path = require('path');
require('dotenv').config();

async function main() {
    try {
        // wallet 디렉토리 생성
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // CA 클라이언트 생성
        const caURL = `https://${process.env.PEER_HOST}:7054`;
        
        const ca = new FabricCAServices(caURL, {
            trustedRoots: [fs.readFileSync(process.env.CA_PATH)],
            verify: false
        }, process.env.CA_NAME);

        // 관리자 등록
        const adminExists = await wallet.get('admin');
        if (adminExists) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        // 관리자 등록
        const enrollment = await ca.enroll({
            enrollmentID: 'admin',
            enrollmentSecret: 'adminpw'
        });

        // X.509 인증서 생성
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'Org1MSP',
            type: 'X.509',
        };

        // wallet에 저장
        await wallet.put('admin', x509Identity);
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        process.exit(1);
    }
}

main(); 