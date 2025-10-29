// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract AuctionNFT is ERC721, Ownable {
    uint256 private _nextTokenId;
    
    constructor() ERC721("AuctionNFT", "ANFT") Ownable(msg.sender) {}
    
    function mint(address to) public onlyOwner returns (uint256) {
        uint256 tokenId = _nextTokenId++;
        _mint(to, tokenId);
        return tokenId;
    }
    
    function getNextTokenId() public view returns (uint256) {
        return _nextTokenId;
    }
}
