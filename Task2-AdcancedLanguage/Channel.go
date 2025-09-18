package main

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 第一题：无缓冲通道通信 ====================

func unbufferedChannelDemo() {
	fmt.Println("=== 第一题：无缓冲通道通信 ===")
	fmt.Println("生产者生成1-10的整数，消费者接收并打印")
	fmt.Println("----------------------------------------")

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
	fmt.Println("=== 第二题：缓冲通道通信 ===")
	fmt.Println("生产者发送100个整数到缓冲通道，消费者接收并打印")
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

// ==================== 扩展：多生产者和多消费者模式 ====================

func multipleProducersConsumersDemo() {
	fmt.Println("=== 扩展：多生产者和多消费者模式 ===")
	fmt.Println("2个生产者，3个消费者，缓冲通道容量为20")
	fmt.Println("----------------------------------------")

	ch := make(chan int, 20)
	const numbersPerProducer = 30
	totalNumbers := 2 * numbersPerProducer

	var wg sync.WaitGroup

	// 启动2个生产者
	for producerID := 1; producerID <= 2; producerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			start := (id-1)*numbersPerProducer + 1
			end := id * numbersPerProducer

			for i := start; i <= end; i++ {
				ch <- i
				if i%10 == 0 {
					fmt.Printf("生产者%d发送: %d\n", id, i)
				}
				time.Sleep(15 * time.Millisecond)
			}
			fmt.Printf("生产者%d完成\n", id)
		}(producerID)
	}

	// 启动3个消费者
	consumerWG := sync.WaitGroup{}
	for consumerID := 1; consumerID <= 3; consumerID++ {
		consumerWG.Add(1)
		go func(id int) {
			defer consumerWG.Done()
			count := 0
			for num := range ch {
				count++
				// 这里使用num变量，避免"declared and not used"错误
				if count%8 == 0 {
					fmt.Printf("消费者%d接收: %d\n", id, num)
				}
				time.Sleep(20 * time.Millisecond)
			}
			fmt.Printf("消费者%d完成，接收了 %d 个数字\n", id, count)
		}(consumerID)
	}

	// 等待所有生产者完成，然后关闭通道
	go func() {
		wg.Wait()
		close(ch)
	}()

	consumerWG.Wait()
	fmt.Printf("所有任务完成，总共处理了 %d 个数字\n", totalNumbers)
}

func main() {
	// 第一题：无缓冲通道
	unbufferedChannelDemo()

	// 第二题：缓冲通道
	bufferedChannelDemo()

	// 扩展：多生产者和多消费者
	multipleProducersConsumersDemo()

	fmt.Println("\n=== 通道特性总结 ===")
	fmt.Println("1. 无缓冲通道: 发送和接收操作会阻塞，直到另一方准备好")
	fmt.Println("2. 缓冲通道: 发送操作只在缓冲区满时阻塞，接收操作只在缓冲区空时阻塞")
	fmt.Println("3. close(ch): 关闭通道表示没有更多数据发送")
	fmt.Println("4. range ch: 自动从通道接收数据直到通道关闭")
	fmt.Println("5. len(ch): 获取通道当前元素数量")
	fmt.Println("6. cap(ch): 获取通道容量")
}