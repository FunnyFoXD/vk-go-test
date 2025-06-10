package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	pools "vk-go-test/pool"
)

func main() {
	// Init worker pool with 2 workers
	pool := pools.NewWorkerPool(2)

	// Goroutine to submit jobs to the pool
	go func() {
		for i := 1; ; i++ {
			pool.SubmitJob(fmt.Sprintf("task-%d", i))
			time.Sleep(500 * time.Millisecond) // Add new task every 500 ms
		}
	}()

	// Interface for console
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Worker pool started. Commands:")
	fmt.Println("  add - add worker")
	fmt.Println("  remove - remove worker")
	fmt.Println("  count - show worker count")
	fmt.Println("  status - show current tasks")
	fmt.Println("  exit - stop program")
	fmt.Print("> ")

	// Main program loop
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "add":
			pool.AddWorker()
			fmt.Print("> ")
		case "remove":
			pool.RemoveWorker()
			fmt.Print("> ")
		case "count":
			fmt.Printf("Current workers: %d\n", pool.WorkerCount())
			fmt.Print("> ")
		case "status":
			pool.PrintStatus()
			fmt.Print("> ")
		case "exit":
			pool.Stop()
			fmt.Println("Program stopped")
			return
		default:
			fmt.Println("Unknown command. Available: add, remove, count, status, exit")
			fmt.Print("> ")
		}
	}
}
