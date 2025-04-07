import { Injectable } from '@nestjs/common';
import { Gateway, Wallets, X509Identity } from '@hyperledger/fabric-gateway';
import * as grpc from '@grpc/grpc-js';
import * as path from 'path';
import * as fs from 'fs';

@Injectable()
export class BlockchainService {
  private gateway: Gateway;
  private network: any;
  private contract: any;

  constructor() {
    this.initializeGateway();
  }

  private async initializeGateway() {
    // 인증서 및 키 파일 경로
    const certPath = path.resolve(__dirname, '../../organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt');
    const keyPath = path.resolve(__dirname, '../../organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk');
    const cert = fs.readFileSync(certPath);
    const key = fs.readFileSync(keyPath);

    // X509 Identity 생성
    const identity: X509Identity = {
      credentials: {
        certificate: cert.toString(),
        privateKey: key.toString(),
      },
      mspId: 'Org1MSP',
    };

    // gRPC 연결 설정
    const clientTlsCredentials = grpc.credentials.createSsl(cert);
    const peerEndpoint = 'localhost:7051';

    // Gateway 연결
    this.gateway = await Gateway.connect({
      identity,
      signer: identity.credentials.privateKey,
      clientTlsCredentials,
      peerEndpoint,
    });

    // 네트워크 및 컨트랙트 설정
    this.network = this.gateway.getNetwork('mychannel');
    this.contract = this.network.getContract('school');
  }

  async createWallet(githubId: string): Promise<void> {
    try {
      await this.contract.submitTransaction('CreateWallet', githubId);
    } catch (error) {
      throw new Error(`Failed to create wallet: ${error.message}`);
    }
  }

  async validateAndRewardCommit(githubId: string, commitData: any): Promise<void> {
    try {
      const commitJSON = JSON.stringify(commitData);
      await this.contract.submitTransaction('ValidateAndRewardCommit', githubId, commitJSON);
    } catch (error) {
      throw new Error(`Failed to validate and reward commit: ${error.message}`);
    }
  }

  async transfer(fromGithubId: string, toGithubId: string, amount: number): Promise<void> {
    try {
      await this.contract.submitTransaction('Transfer', fromGithubId, toGithubId, amount.toString());
    } catch (error) {
      throw new Error(`Failed to transfer tokens: ${error.message}`);
    }
  }

  async getWallet(githubId: string): Promise<any> {
    try {
      const result = await this.contract.evaluateTransaction('GetWallet', githubId);
      return JSON.parse(result.toString());
    } catch (error) {
      throw new Error(`Failed to get wallet: ${error.message}`);
    }
  }

  async getAllWallets(): Promise<any[]> {
    try {
      const result = await this.contract.evaluateTransaction('GetAllWallets');
      return JSON.parse(result.toString());
    } catch (error) {
      throw new Error(`Failed to get all wallets: ${error.message}`);
    }
  }

  async onModuleDestroy() {
    if (this.gateway) {
      await this.gateway.close();
    }
  }
} 