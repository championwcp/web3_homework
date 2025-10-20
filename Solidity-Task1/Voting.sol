// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting {
    // 存储候选人得票数的 mapping
    mapping(string => uint256) public candidateVotes;
    
    // 存储所有候选人的数组，用于重置功能
    string[] private candidateList;
    
    
    // 投票函数
    function vote(string memory candidate) public {
        require(bytes(candidate).length > 0, "Candidate name cannot be empty");
        
        if (candidateVotes[candidate] == 0) {
            candidateList.push(candidate);
        }
        
        candidateVotes[candidate]++;
    }
    
    // 获取候选人得票数
    function getVotes(string memory candidate) public view returns (uint256) {
        require(bytes(candidate).length > 0, "Candidate name cannot be empty");
        return candidateVotes[candidate];
    }
    
    // 重置所有候选人的得票数
    function resetVotes() public {
        for (uint256 i = 0; i < candidateList.length; i++) {
            candidateVotes[candidateList[i]] = 0;
        }
    }
}