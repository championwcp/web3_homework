package main

import "fmt"

// 接收整数指针作为参数，将指向的值增加10
func increaseByTen(ptr *int) {
    *ptr += 10 // 解引用指针并修改值
}

// 接收整数切片指针，将每个元素乘以2
func doubleSliceElements(slicePtr *[]int) {
    // 解引用切片指针
    slice := *slicePtr
    
    // 遍历切片并修改每个元素
    for i := 0; i < len(slice); i++ {
        slice[i] *= 2
    }
}

func main() {
    // num := 5
    // fmt.Println("修改前的值:", num)
    // // 传递变量的地址给函数
    // increaseByTen(&num)
    // fmt.Println("修改后的值:", num)


    numbers := []int{1, 2, 3, 4, 5}
    fmt.Println("修改前的切片:", numbers)
    // 传递切片指针给函数
    doubleSliceElements(&numbers)
    fmt.Println("修改后的切片:", numbers)
}