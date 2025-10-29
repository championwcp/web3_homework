cat > README.md << 'EOF'
# NFT拍卖市场

基于Hardhat开发的NFT拍卖市场，支持ERC721 NFT拍卖、Chainlink价格预言机和工厂模式管理。

## 功能特性

### 核心功能
- **NFT合约** (AuctionNFT): ERC721标准的NFT代币
- **拍卖合约** (Auction): 支持创建拍卖、出价、结束拍卖
- **预言机集成** (AuctionWithOracle): Chainlink ETH/USD价格转换
- **工厂模式** (AuctionFactory): 管理多个拍卖实例

### 技术特性
- Solidity 0.8.20
- OpenZeppelin合约库
- Chainlink价格预言机
- 重入攻击防护
- 完整的测试覆盖

## 合约地址 (Sepolia测试网)

| 合约 | 地址 |
|------|------|
| AuctionNFT | `0x710A2CEB504EEEA5B303a41d8fdae7E23dA5AeB4` |
| Auction | `0x2839fF636f4957B6b1719a0dd3428C8250eF8391` |
| AuctionWithOracle | `0x2eDEAE1c5ecb116BC925C5f93204D3eEEa2EcA25` |
| AuctionFactory | `0xfa1314a9C664cFb18395776F10Fc58512F572c9F` |

## 项目结构
├── contracts/ # 智能合约
├── scripts/ # 部署脚本
├── test/ # 测试文件
├── hardhat.config.js # Hardhat配置
└── README.md # 项目文档
└── TEST-REPORT.md # 测试报告


## 快速开始

### 环境要求
- Node.js 16+
- npm 或 yarn

### 安装依赖
```bash
npm install

### 编译合约
npx hardhat compile

### 运行测试
npx hardhat test

### 部署到Sepolia
npx hardhat run scripts/deploy.js --network sepolia


## 部署地址
本次部署成功在Sepolia测试网部署了以下合约：

合约	               地址	
AuctionNFT	0x710A2CEB504EEEA5B303a41d8fdae7E23dA5AeB4
Auction	0x2839fF636f4957B6b1719a0dd3428C8250eF8391
AuctionWithOracle	0x2eDEAE1c5ecb116BC925C5f93204D3eEEa2EcA25
AuctionFactory	0xfa1314a9C664cFb18395776F10Fc58512F572c9F
