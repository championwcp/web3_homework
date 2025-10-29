const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Basic Tests", function () {
  it("应该部署NFT合约", async function () {
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    const nft = await AuctionNFT.deploy();
    
    expect(nft.target).to.be.a("string");
  });

  it("应该部署拍卖合约", async function () {
    const Auction = await ethers.getContractFactory("Auction");
    const auction = await Auction.deploy();
    
    expect(auction.target).to.be.a("string");
  });

  it("应该铸造NFT", async function () {
    const [owner] = await ethers.getSigners();
    const AuctionNFT = await ethers.getContractFactory("AuctionNFT");
    const nft = await AuctionNFT.deploy();
    
    await nft.mint(owner.address);
    expect(await nft.ownerOf(0)).to.equal(owner.address);
  });
});
