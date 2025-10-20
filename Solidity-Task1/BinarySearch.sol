// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    /**
     * @dev 在有序数组中二分查找目标值
     * @param nums 有序数组（升序）
     * @param target 目标值
     * @return 目标值的索引，如果不存在则返回 type(uint256).max
     */
    function binarySearch(
        uint256[] memory nums,
        uint256 target
    ) public pure returns (uint256) {
        uint256 left = 0;
        uint256 right = nums.length;
        
        while (left < right) {
            uint256 mid = left + (right - left) / 2;
            
            if (nums[mid] == target) {
                return mid;
            } else if (nums[mid] < target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        
        return type(uint256).max; // 表示未找到
    }
}