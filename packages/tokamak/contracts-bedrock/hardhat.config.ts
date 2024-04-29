import "@nomiclabs/hardhat-waffle";
import "@nomicfoundation/hardhat-verify";
import 'hardhat-deploy';

export default {
  solidity: "0.8.9",
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
          apiURL: "http://localhost/api",
          browserURL: "http://localhost"
        }
      }
    ]
  },
  sourcify: {
    enabled: false
  }
};
