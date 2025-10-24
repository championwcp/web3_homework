// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleToken {
    // 代币基本信息
    string public constant name = "MyTestToken";
    string public constant symbol = "MTT";
    uint8 public constant decimals = 18;
    
    // 状态变量
    uint256 public totalSupply;
    address public owner;
    
    // 映射
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    
    // 事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Mint(address indexed to, uint256 value);
    
    // 修饰器
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    // 简化构造函数 - 无参数
    constructor() {
        owner = msg.sender;
        totalSupply = 1000000 * 10 ** decimals; // 100万代币
        balanceOf[msg.sender] = totalSupply;
        emit Transfer(address(0), msg.sender, totalSupply);
    }
    
    // 转账函数
    function transfer(address to, uint256 value) external returns (bool) {
        return _transfer(msg.sender, to, value);
    }
    
    // 授权函数
    function approve(address spender, uint256 value) external returns (bool) {
        allowance[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }
    
    // 代扣转账
    function transferFrom(address from, address to, uint256 value) external returns (bool) {
        uint256 currentAllowance = allowance[from][msg.sender];
        require(currentAllowance >= value, "Allowance too low");
        
        allowance[from][msg.sender] = currentAllowance - value;
        return _transfer(from, to, value);
    }
    
    // 增发代币
    function mint(address to, uint256 value) external onlyOwner {
        require(to != address(0), "Zero address");
        
        totalSupply += value;
        balanceOf[to] += value;
        
        emit Mint(to, value);
        emit Transfer(address(0), to, value);
    }
    
    // 内部转账函数
    function _transfer(address from, address to, uint256 value) internal returns (bool) {
        require(from != address(0), "From zero address");
        require(to != address(0), "To zero address");
        require(balanceOf[from] >= value, "Insufficient balance");
        
        balanceOf[from] -= value;
        balanceOf[to] += value;
        
        emit Transfer(from, to, value);
        return true;
    }
}