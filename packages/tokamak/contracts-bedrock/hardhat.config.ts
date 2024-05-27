import "@nomiclabs/hardhat-waffle";
// import "@nomicfoundation/hardhat-verify";
import "@nomiclabs/hardhat-etherscan";
import 'hardhat-deploy';

const LOWEST_OPTIMIZER_COMPILER_SETTINGS = {
  version: '0.7.6',
  settings: {
    evmVersion: 'istanbul',
    optimizer: {
      enabled: true,
      runs: 1_000_000,
    },
    metadata: {
      bytecodeHash: 'none',
    },
  },
};

export default {
  solidity: {
    compilers: [
      LOWEST_OPTIMIZER_COMPILER_SETTINGS,
    ],
  },
  networks: {
    devnetL1: {
      url: "http://localhost:9545",
      chainId: 901,
    },
    'thanos-sepolia-test': {
      url: "https://explorer.thanos-sepolia-test.tokamak.network",
      chainId: 111551118080,
    },
  },
  etherscan: {
    apiKey: "YOUR_API_KEY",
    customChains: [
      {
        network: "devnetL1",
        chainId: 901,
        urls: {
          apiURL: "http://localhost/api",
          browserURL: "http://localhost",
        },
      },
      {
        network: "thanos-sepolia-test",
        chainId: 111551118080,
        urls: {
          apiURL: "https://rpc.titan-sepolia.tokamak.network",
          browserURL: "https://explorer.thanos-sepolia-test.tokamak.network",
        },
      },
    ],
  },
  sourcify: {
    enabled: false,
  },
};
