async function main() {
  console.log("开始部署NFT拍卖市场合约到Sepolia...");
  
  // Sepolia ETH/USD Price Feed地址
  const SEPOLIA_ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
  
  const [deployer] = await ethers.getSigners();
  console.log("部署者地址:", deployer.address);
  console.log("部署者余额:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)), "ETH");

  // 部署AuctionNFT
  console.log("正在部署AuctionNFT...");
  const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
  const nft = await AuctionNFT.deploy();
  await nft.waitForDeployment();
  console.log("✅ AuctionNFT部署地址:", await nft.getAddress());

  // 部署基础Auction合约
  console.log("正在部署Auction...");
  const Auction = await ethers.getContractFactory("Auction");
  const auction = await Auction.deploy();
  await auction.waitForDeployment();
  console.log("✅ Auction部署地址:", await auction.getAddress());

  // 部署带预言机的AuctionWithOracle合约
  console.log("正在部署AuctionWithOracle...");
  const AuctionWithOracle = await ethers.getContractFactory("AuctionWithOracle");
  const auctionWithOracle = await AuctionWithOracle.deploy(SEPOLIA_ETH_USD_PRICE_FEED);
  await auctionWithOracle.waitForDeployment();
  console.log("✅ AuctionWithOracle部署地址:", await auctionWithOracle.getAddress());

  // 部署工厂合约
  console.log("正在部署AuctionFactory...");
  const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
  const factory = await AuctionFactory.deploy();
  await factory.waitForDeployment();
  console.log("✅ AuctionFactory部署地址:", await factory.getAddress());

  console.log("🎉 所有合约部署完成！");
  console.log("=== 合约地址汇总 ===");
  console.log("AuctionNFT:", await nft.getAddress());
  console.log("Auction:", await auction.getAddress());
  console.log("AuctionWithOracle:", await auctionWithOracle.getAddress());
  console.log("AuctionFactory:", await factory.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
