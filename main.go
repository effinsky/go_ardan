package main

import (
	"fmt"
	"strings"
)

func main() {
	// conc.SyncWithMutex()
	// conc.WaitForResult()
	// conc.FanOut()
	// conc.Pooling()
	// conc.FanOutSemaphore()
	// conc.FanOutBounded()
	// conc.Drop()
	// conc.Cancellation()
	_ = processRole("nothingness")
}

type role string

func (r *role) reverse() {
	fmt.Printf("r: %v\n", *r)
	*r = role(strings.ToUpper(string(*r)))
	fmt.Printf("r: %v\n", *r)
}

func processRole(r role) error {
	r.reverse()
	return nil
}
