package main

import (
	"fmt"
	"time"
)

// Sleeper interface abstracts time.Sleep
type Sleeper interface {
	Sleep(duration time.Duration)
}

// Worker handles operations that take time
type Worker struct {
	sleeper Sleeper
}

func NewWorker(sleeper Sleeper) *Worker {
	return &Worker{sleeper: sleeper}
}

// DoWork simulates work that takes time
func (w *Worker) DoWork() {
	fmt.Println("Starting work...")
	w.sleeper.Sleep(2 * time.Second)
	fmt.Println("Work complete!")
}

type realSleeper struct{}

func (r realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

func main() {
	start := time.Now()
	fmt.Printf("Started at: %s\n", start.Format("15:04:05"))

	worker := NewWorker(realSleeper{})
	worker.DoWork()

	end := time.Now()
	fmt.Printf("Finished at: %s\n", end.Format("15:04:05"))
	fmt.Printf("Duration: %s\n", end.Sub(start))
}
