const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Auction Functionality", function () {
  let auction, nft;
  let owner, seller, bidder1, bidder2;

  beforeEach(async function () {
    [owner, seller, bidder1, bidder2] = await ethers.getSigners();
    
    // 部署合约
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    nft = await AuctionNFT.deploy();
    
    const Auction = await ethers.getContractFactory("Auction");
    auction = await Auction.deploy();
    
    // 准备测试数据：给卖家铸造NFT并授权
    await nft.mint(seller.address);
    await nft.connect(seller).approve(auction.target, 0);
  });

  describe("创建拍卖", function () {
    it("应该成功创建拍卖并转移NFT", async function () {
      // 创建拍卖
      await auction.connect(seller).createAuction(nft.target, 0, 3600, ethers.parseEther("1.0"));
      
      // 检查拍卖信息
      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.seller).to.equal(seller.address);
      expect(auctionInfo.nftContract).to.equal(nft.target);
      expect(auctionInfo.tokenId).to.equal(0);
      expect(auctionInfo.startPrice).to.equal(ethers.parseEther("1.0"));
      expect(auctionInfo.ended).to.be.false;
      
      // 检查NFT所有权转移到拍卖合约
      expect(await nft.ownerOf(0)).to.equal(auction.target);
    });

    it("应该递增拍卖计数器", async function () {
      expect(await auction.auctionCount()).to.equal(0);
      
      await auction.connect(seller).createAuction(nft.target, 0, 3600, ethers.parseEther("1.0"));
      expect(await auction.auctionCount()).to.equal(1);
      
      // 再创建一个拍卖
      await nft.mint(seller.address);
      await nft.connect(seller).approve(auction.target, 1);
      await auction.connect(seller).createAuction(nft.target, 1, 3600, ethers.parseEther("2.0"));
      expect(await auction.auctionCount()).to.equal(2);
    });

    it("应该要求正数的拍卖时长", async function () {
      await expect(
        auction.connect(seller).createAuction(nft.target, 0, 0, ethers.parseEther("1.0"))
      ).to.be.revertedWith("Duration must be positive");
    });
  });

  describe("出价功能", function () {
    beforeEach(async function () {
      // 在每个出价测试前创建拍卖
      await auction.connect(seller).createAuction(nft.target, 0, 3600, ethers.parseEther("1.0"));
    });

    it("应该接受更高的出价", async function () {
      await auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("1.5") });
      
      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.highestBidder).to.equal(bidder1.address);
      expect(auctionInfo.highestBid).to.equal(ethers.parseEther("1.5"));
    });

    it("应该拒绝更低的出价", async function () {
      await auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("1.5") });
      
      await expect(
        auction.connect(bidder2).placeBid(1, { value: ethers.parseEther("1.4") })
      ).to.be.revertedWith("Bid too low");
    });

    it("应该退还前一个出价者的资金", async function () {
      const bidder1BalanceBefore = await ethers.provider.getBalance(bidder1.address);
      
      // 第一个出价
      await auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("1.5") });
      
      // 第二个更高的出价，应该退还第一个出价者的资金
      await auction.connect(bidder2).placeBid(1, { value: ethers.parseEther("2.0") });
      
      const bidder1BalanceAfter = await ethers.provider.getBalance(bidder1.address);
      // 注意：这里需要考虑到gas费用，所以余额应该接近原始余额
      expect(bidder1BalanceAfter).to.be.closeTo(bidder1BalanceBefore, ethers.parseEther("0.1"));
    });

    it("不应该在拍卖结束后接受出价", async function () {
      // 增加时间到拍卖结束
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("1.5") })
      ).to.be.revertedWith("Auction ended");
    });
  });

  describe("结束拍卖", function () {
    beforeEach(async function () {
      await auction.connect(seller).createAuction(nft.target, 0, 3600, ethers.parseEther("1.0"));
      await auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("1.5") });
    });

    it("应该成功结束有出价的拍卖", async function () {
      // 增加时间到拍卖结束
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);
      
      await auction.connect(seller).endAuction(1);
      
      // 检查NFT转移给出价最高者
      expect(await nft.ownerOf(0)).to.equal(bidder1.address);
      
      // 检查卖家收到资金
      const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);
      expect(sellerBalanceAfter).to.be.gt(sellerBalanceBefore);
      
      // 检查拍卖状态
      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.ended).to.be.true;
    });

    it("应该处理无出价的拍卖", async function () {
      // 创建新的无出价拍卖
      await nft.mint(seller.address);
      await nft.connect(seller).approve(auction.target, 1);
      await auction.connect(seller).createAuction(nft.target, 1, 3600, ethers.parseEther("1.0"));
      
      // 增加时间到拍卖结束
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await auction.connect(seller).endAuction(2);
      
      // 检查NFT归还给卖家
      expect(await nft.ownerOf(1)).to.equal(seller.address);
    });

    it("不应该在拍卖未结束时结束拍卖", async function () {
      await expect(
        auction.connect(seller).endAuction(1)
      ).to.be.revertedWith("Auction not ended");
    });

    it("只有卖家或合约所有者可以结束拍卖", async function () {
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      // 非卖家或所有者尝试结束拍卖
      await expect(
        auction.connect(bidder2).endAuction(1)
      ).to.be.revertedWith("Not authorized");
      
      // 所有者可以结束拍卖
      await auction.connect(owner).endAuction(1);
    });
  });

  describe("查询功能", function () {
    it("应该返回正确的拍卖信息", async function () {
      await auction.connect(seller).createAuction(nft.target, 0, 3600, ethers.parseEther("1.0"));
      
      const auctionInfo = await auction.getAuction(1);
      
      expect(auctionInfo.seller).to.equal(seller.address);
      expect(auctionInfo.nftContract).to.equal(nft.target);
      expect(auctionInfo.tokenId).to.equal(0);
      expect(auctionInfo.startPrice).to.equal(ethers.parseEther("1.0"));
      expect(auctionInfo.highestBidder).to.equal(ethers.ZeroAddress);
      expect(auctionInfo.highestBid).to.equal(ethers.parseEther("1.0"));
      expect(auctionInfo.ended).to.be.false;
    });
  });
});
