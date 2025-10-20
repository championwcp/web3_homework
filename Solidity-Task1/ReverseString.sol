// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StringReverser {
    /**
     * @dev 反转字符串
     * @param str 输入的字符串
     * @return 反转后的字符串
     */
    function reverseString(string memory str) public pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        uint256 length = strBytes.length;
        
        if (length <= 1) {
            return str;
        }
        
        bytes memory reversed = new bytes(length);
        
        for (uint256 i = 0; i < length; i++) {
            reversed[i] = strBytes[length - 1 - i];
        }
        
        return string(reversed);
    }
    
}