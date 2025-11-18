// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 private count;
    address public owner;
    
    event CountIncreased(uint256 newCount, address increasedBy);
    event CountReset(address resetBy);
    
    constructor() {
        owner = msg.sender;
        count = 0;
    }
    
    function getCount() public view returns (uint256) {
        return count;
    }
    
    function increase() public returns (uint256) {
        count += 1;
        emit CountIncreased(count, msg.sender);
        return count;
    }
    
    function reset() public {
        require(msg.sender == owner, "Only owner can reset");
        count = 0;
        emit CountReset(msg.sender);
    }
    
    function increaseBy(uint256 number) public returns (uint256) {
        count += number;
        emit CountIncreased(count, msg.sender);
        return count;
    }
}