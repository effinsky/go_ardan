package concpatterns

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var counter int32

func SyncWithMutex() {
	const grs = 2

	var wg sync.WaitGroup
	var mu sync.RWMutex
	wg.Add(grs)

	for i := 0; i < grs; i++ {
		go func() {
			defer wg.Done()

			for count := 0; count < 2; count++ {
				mu.Lock()
				{
					value := counter
					value++

					// This print statement enforces context switching on the scheduler
					// and there's writes into a dirty counter value.
					fmt.Println("logging")

					counter = value
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// NOTE: At the end here the counter should be 4, since we have 2 routines incrementing
	// it by 2 each.
	fmt.Printf("counter: %v\n", counter)
}

// NOTE: channel patterns
func WaitForResult() {
	ch := make(chan string)

	go func() {
		time.Sleep(time.Duration(rand.Intn(500) * int(time.Millisecond)))
		ch <- "paper"
		println("signal sent")
	}()

	p := <-ch
	fmt.Printf("signal received: %v\n", p)

	time.Sleep(time.Second)
}

func FanOut() {
	emps := 2000
	ch := make(chan string, emps)

	for e := range emps {
		go func(emp int) {
			time.Sleep(time.Duration(rand.Intn(200) * int(time.Millisecond)))
			ch <- "paper"
			fmt.Printf("sent signal: emp %d\n", emp)
		}(e)
	}

	for emps > 0 {
		p := <-ch
		emps--
		fmt.Printf("p: %v\n", p)
		fmt.Printf("received signal: %v\n", emps)
	}

	time.Sleep(time.Second)
	fmt.Println("---------------------------------------------------")
}

func FanOutSemaphore() {
	workers := 2000
	ch := make(chan string, workers)

	grs := runtime.NumCPU()
	sem := make(chan bool, grs)

	for w := range workers {
		go func() {
			sem <- true
			{
				time.Sleep(time.Duration(rand.Intn(200) * int(time.Millisecond)))
				ch <- "paper"
				// capture that loop variable
				fmt.Println("worker : signal sent: ", w)
			}
			<-sem
		}()
	}

	for workers > 0 {
		p := <-ch
		workers--
		fmt.Printf("%s\n", p)
		fmt.Printf("signal received: %d\n", workers)
	}
}

func Pooling() {
	ch := make(chan string)
	grs := runtime.NumCPU()

	// routine management
	for e := range grs {
		go func() {
			for p := range ch {
				fmt.Printf("emp %d received signal:) %v\n", e, p)
			}
			fmt.Printf("emp %d received shutdown signal\n", e)
		}()
	}

	// work management
	for work := range 100 {
		ch <- "paper"
		fmt.Printf("signal sent: %v\n", work)
	}

	close(ch)
	fmt.Printf("shutdown signal sent\n")
	fmt.Println("---------------------------------------------------")
}

func FanOutBounded() {
	work := []string{"paper", "paper", "paper", 2000: "paper"}
	fmt.Printf("work: %v\n", work)

	grs := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(grs)

	ch := make(chan string, grs) // buf chan for workers

	for worker := range grs {
		go func() {
			defer wg.Done()
			for p := range ch {
				fmt.Printf("%d : received signal : %s\n", worker, p)
			}
			fmt.Printf("%d : received shutdown signal\n", worker)
		}()
	}

	for _, w := range work {
		ch <- w
	}

	close(ch)
	wg.Wait()
	fmt.Println("---------------------------------------------------")
}

func Drop() {
	cap := 100
	ch := make(chan string, cap) // blocks over 100

	// a single routine to keep the code simple
	go func() {
		// receiving
		for p := range ch {
			fmt.Println("signal received:", p)
		}
	}()

	const work = 2000
	// sending
	for w := range work {
		select {
		// if this blocks, go to default and drop
		case ch <- "paper":
			fmt.Printf("signal sent: %d\n", w)
		default:
			fmt.Printf("data dropped\n")
		}
	}
	close(ch)

	time.Sleep(time.Second)
	fmt.Printf("sending shutdown signal\n")
	fmt.Println("---------------------------------------------------")
}

func Cancellation() {
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// on an unbuffered channel, a routine might leak when continues past the timeout
	// deadline and has no receiver to guarantee it's finished (goroutine leak)
	ch := make(chan string, 1)

	go func() {
		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
		ch <- "paper"
	}()

	select {
	// clock starts ticking on the 150ms timeout. now the first case has that time and no
	// more to get the work done as both cases block. the done case will unblock once
	// the timeout is up.
	case d := <-ch:
		fmt.Printf("work complete: %s\n", d)

	case <-ctx.Done():
		fmt.Println("work cancelled")
	}

	time.Sleep(time.Second)
	fmt.Printf("sending shutdown signal\n")
	fmt.Println("---------------------------------------------------")
}

// RunWorker with stop func
func RunWorker(process func(int, time.Time) error) func() {
	const (
		estimatedCheckFreq time.Duration = time.Second * 5
		checkSelectLimit   int           = 50
		shutdownTimeout    time.Duration = time.Second * 15
	)

	ticker := time.NewTicker(estimatedCheckFreq)
	processStopCh := make(chan chan struct{}, 1)
	go func() {
		for {
			select {
			case now := <-ticker.C:
				if err := process(checkSelectLimit, now); err != nil {
					fmt.Printf("Failed to run check: %v", err)
				}
			case ret := <-processStopCh:
				close(ret)
				return
			}
		}
	}()

	return func() {
		ticker.Stop()
		ret := make(chan struct{})
		processStopCh <- ret
		select {
		case <-time.After(shutdownTimeout):
			fmt.Printf("Failed to gracefully stop check worker\n")
		case <-ret:
			fmt.Printf("Gracefully stopped check worker\n")
		}
	}
}
