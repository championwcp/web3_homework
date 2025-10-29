require("@nomicfoundation/hardhat-ethers");
require("@nomicfoundation/hardhat-chai-matchers");
require("dotenv").config();

module.exports = {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  },
  networks: {
    hardhat: {
      chainId: 1337
    },
    sepolia: {
      url: "https://ethereum-sepolia-rpc.publicnode.com", // 备用节点
      accounts: [process.env.PRIVATE_KEY],
      chainId: 11155111,
      gas: 5000000,
      gasPrice: 25000000000,
      timeout: 60000  // 增加超时时间
    }
  }
};
