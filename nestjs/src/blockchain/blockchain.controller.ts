import { Controller, Post, Body, Get, Param, Put } from '@nestjs/common';
import { BlockchainService } from './blockchain.service';

@Controller('wallets')
export class BlockchainController {
  constructor(private readonly blockchainService: BlockchainService) {}

  @Post()
  async createWallet(@Body('githubId') githubId: string) {
    await this.blockchainService.createWallet(githubId);
    return { message: 'Wallet created successfully' };
  }

  @Get()
  async getAllWallets() {
    return await this.blockchainService.getAllWallets();
  }

  @Get(':githubId')
  async getWallet(@Param('githubId') githubId: string) {
    return await this.blockchainService.getWallet(githubId);
  }

  @Post(':githubId/commits')
  async validateAndRewardCommit(
    @Param('githubId') githubId: string,
    @Body('commitData') commitData: any,
  ) {
    await this.blockchainService.validateAndRewardCommit(githubId, commitData);
    return { message: 'Commit validated and rewarded successfully' };
  }

  @Post(':githubId/transfers')
  async transfer(
    @Param('githubId') githubId: string,
    @Body('toGithubId') toGithubId: string,
    @Body('amount') amount: number,
  ) {
    await this.blockchainService.transfer(githubId, toGithubId, amount);
    return { message: 'Transfer successful' };
  }
} 