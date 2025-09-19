package main

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 第一题：无缓冲通道通信 ====================

func unbufferedChannelDemo() {
	fmt.Println("生产者生成1-10的整数，消费者接收")
	// 创建无缓冲通道
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	// 生产者协程
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 1; i <= 10; i++ {
			fmt.Printf("生产者发送: %d\n", i)
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("生产者完成")
	}()
	// 消费者协程
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("消费者接收: %d\n", num)
			time.Sleep(150 * time.Millisecond)
		}
		fmt.Println("消费者完成")
	}()
	wg.Wait()
	fmt.Println()
}

// ==================== 第二题：缓冲通道通信 ====================
func bufferedChannelDemo() {
	fmt.Println("生产者发送100个整数到缓冲通道")
	fmt.Println("----------------------------------------")
	ch := make(chan int, 10)
	const totalNumbers = 100
	var wg sync.WaitGroup
	wg.Add(2)
	// 生产者协程
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 1; i <= totalNumbers; i++ {
			ch <- i
			if i%10 == 0 {
				fmt.Printf("生产者已发送: %d/%d (通道长度: %d/%d)\n", 
					i, totalNumbers, len(ch), cap(ch))
			}
			time.Sleep(10 * time.Millisecond)
		}
		fmt.Printf("生产者完成，共发送 %d 个数字\n", totalNumbers)
	}()

	// 消费者协程
	go func() {
		defer wg.Done()
		count := 0
		for num := range ch {
			count++
			fmt.Printf("仅仅打印: %d\n", num)
			if count%15 == 0 {
				fmt.Printf("消费者已接收: %d/%d (通道长度: %d/%d)\n", 
					count, totalNumbers, len(ch), cap(ch))
			}
			time.Sleep(25 * time.Millisecond)
		}
		fmt.Printf("消费者完成，共接收 %d 个数字\n", count)
	}()

	wg.Wait()
	fmt.Println()
}


func main() {
	// 第一题：无缓冲通道
	unbufferedChannelDemo()
	// 第二题：缓冲通道
	bufferedChannelDemo()
}
