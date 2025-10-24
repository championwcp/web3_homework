// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract {
    // 合约所有者地址
    address public owner;
    
    // 每个地址的捐赠总额
    mapping(address => uint256) public donations;
    
    // 总捐赠金额
    uint256 public totalDonations;
    
    // 事件：记录捐赠信息
    event DonationReceived(address indexed donor, uint256 amount);
    
    // 事件：记录提款信息
    event Withdrawal(address indexed owner, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    // 构造函数：在部署时设置合约所有者
    constructor() {
        owner = msg.sender;
    }

    // donate 函数：接收以太币捐赠
    function donate() external payable {
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        donations[msg.sender] += msg.value;
        totalDonations += msg.value;
        emit DonationReceived(msg.sender, msg.value);
    }

    // withdraw 函数：所有者提取所有资金
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        // 转移资金到所有者地址
        payable(owner).transfer(balance);
        
        emit Withdrawal(owner, balance);
    }

    // getDonation 函数：查询特定地址的捐赠金额
    function getDonation(address donor) external view returns (uint256) {
        return donations[donor];
    }

    // 获取合约当前余额
    function getContractBalance() external view returns (uint256) {
        return address(this).balance;
    }
}