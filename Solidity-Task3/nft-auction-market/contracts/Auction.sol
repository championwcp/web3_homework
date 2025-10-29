// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Auction is ReentrancyGuard, Ownable {
    struct AuctionItem {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 startTime;
        uint256 endTime;
        uint256 startPrice;
        address highestBidder;
        uint256 highestBid;
        bool ended;
    }

    mapping(uint256 => AuctionItem) public auctions;
    uint256 public auctionCount;

    event AuctionCreated(uint256 indexed auctionId, address indexed seller, address nftContract, uint256 tokenId);
    event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount);
    event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 amount);

    constructor() Ownable(msg.sender) {}

    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 duration,
        uint256 startPrice
    ) external returns (uint256) {
        require(duration > 0, "Duration must be positive");
        
        IERC721(nftContract).transferFrom(msg.sender, address(this), tokenId);

        auctionCount++;
        auctions[auctionCount] = AuctionItem({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            startTime: block.timestamp,
            endTime: block.timestamp + duration,
            startPrice: startPrice,
            highestBidder: address(0),
            highestBid: startPrice,
            ended: false
        });

        emit AuctionCreated(auctionCount, msg.sender, nftContract, tokenId);
        return auctionCount;
    }

    function placeBid(uint256 auctionId) external payable nonReentrant {
        AuctionItem storage auction = auctions[auctionId];
        
        require(block.timestamp < auction.endTime, "Auction ended");
        require(!auction.ended, "Auction already ended");
        require(msg.value > auction.highestBid, "Bid too low");

        // Return previous bid
        if (auction.highestBidder != address(0)) {
            payable(auction.highestBidder).transfer(auction.highestBid);
        }

        auction.highestBidder = msg.sender;
        auction.highestBid = msg.value;

        emit BidPlaced(auctionId, msg.sender, msg.value);
    }

    function endAuction(uint256 auctionId) external nonReentrant {
        AuctionItem storage auction = auctions[auctionId];
        
        require(block.timestamp >= auction.endTime, "Auction not ended");
        require(!auction.ended, "Auction already ended");
        require(msg.sender == auction.seller || msg.sender == owner(), "Not authorized");

        auction.ended = true;

        if (auction.highestBidder != address(0)) {
            IERC721(auction.nftContract).transferFrom(address(this), auction.highestBidder, auction.tokenId);
            payable(auction.seller).transfer(auction.highestBid);
            emit AuctionEnded(auctionId, auction.highestBidder, auction.highestBid);
        } else {
            IERC721(auction.nftContract).transferFrom(address(this), auction.seller, auction.tokenId);
        }
    }

    function getAuction(uint256 auctionId) external view returns (AuctionItem memory) {
        return auctions[auctionId];
    }
}
