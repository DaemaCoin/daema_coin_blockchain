require('dotenv').config();
const fs = require('fs');

const connectionProfile = {
    name: 'fabric-network',
    version: '1.0.0',
    client: {
        organization: 'Org1',
        connection: {
            timeout: {
                peer: { endorser: '300' },
                orderer: '300',
            },
        },
    },
    organizations: {
        Org1: {
            mspid: 'Org1MSP',
            peers: ['peer0.org1.example.com'],
            certificateAuthorities: ['ca.org1.example.com'],
        },
    },
    peers: {
        'peer0.org1.example.com': {
            url: `grpcs://${process.env.PEER_HOST}:${process.env.PEER_PORT}`,
            tlsCACerts: {
                path: process.env.CA_PATH,
            },
            grpcOptions: {
                'ssl-target-name-override': 'peer0.org1.example.com',
                hostnameOverride: 'peer0.org1.example.com',
            },
        },
    },
    orderers: {
        'orderer.example.com': {
            url: `grpcs://${process.env.ORDERER_HOST}:${process.env.ORDERER_PORT}`,
            tlsCACerts: {
                path: process.env.CA_PATH,
            },
            grpcOptions: {
                'ssl-target-name-override': 'orderer.example.com',
                hostnameOverride: 'orderer.example.com',
            },
        },
    },
    certificateAuthorities: {
        'ca.org1.example.com': {
            url: `https://${process.env.PEER_HOST}:7054`,
            caName: 'ca-org1',
            tlsCACerts: {
                path: process.env.CA_PATH,
            },
            httpOptions: {
                verify: false,
            },
        },
    },
};

module.exports = connectionProfile;
