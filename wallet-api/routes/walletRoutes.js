const express = require('express');
const router = express.Router();
const walletController = require('../controllers/walletController');

// 지갑 생성
router.post('/', walletController.createWallet);

// 지갑 조회
router.get('/:owner', walletController.readWallet);

// 커밋 보상 지급
router.post('/reward', walletController.rewardForCommit);

// 토큰 전송
router.post('/transfer', walletController.transfer);

// 토큰 소각
router.post('/burn', walletController.burn);

// 커밋 기록 조회
router.get('/commit/:hash', walletController.getCommitRecord);

// 내 잔액 조회
router.get('/balance/me', walletController.getMyBalance);

module.exports = router;
