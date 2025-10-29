// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./Auction.sol";

contract AuctionFactory {
    address[] public auctions;
    mapping(address => bool) public isAuction;
    
    event AuctionCreated(address indexed auctionAddress, address indexed creator, uint256 timestamp);
    
    function createAuction() external returns (address) {
        Auction newAuction = new Auction();
        address auctionAddress = address(newAuction);
        
        auctions.push(auctionAddress);
        isAuction[auctionAddress] = true;
        
        emit AuctionCreated(auctionAddress, msg.sender, block.timestamp);
        return auctionAddress;
    }
    
    function getAllAuctions() external view returns (address[] memory) {
        return auctions;
    }
    
    function getAuctionsCount() external view returns (uint256) {
        return auctions.length;
    }
}
