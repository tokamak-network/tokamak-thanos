import "@nomiclabs/hardhat-waffle";
import "@nomicfoundation/hardhat-verify";
import 'hardhat-deploy';

export default {
  solidity: {
    compilers: [
      {
        version: "0.8.9"
      },
      {
        version: "0.7.6"
      },
      {
        version: "0.7.0"
      },
      {
        version: "0.5.0"
      },
      {
        version: "0.8.17"
      }
    ]
  },
  networks: {
    devnetL1: {
      url: "http://localhost:9545",
      chainId: 901,
    },
    'thanos-sepolia-test': {
      url: "https://explorer.thanos-sepolia-test.tokamak.network",
      chainId: 111551118080,
    }
  },
  etherscan: {
    apiKey: "YOUR_API_KEY",
    customChains: [
      {
        network: "devnetL1",
        chainId: 901,
        urls: {
          apiURL: "http://localhost/api?",
          browserURL: "http://localhost"
        }
      },
      {
        network: "thanos-sepolia-test",
        chainId: 111551118080,
        urls: {
          apiURL: "https://rpc.titan-sepolia.tokamak.network",
          browserURL: "https://explorer.thanos-sepolia-test.tokamak.network"
        }
      }
    ]
  },
  sourcify: {
    enabled: false
  }
};
