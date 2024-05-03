package main

import (
	conc "ardan/conc_patterns"
)

func main() {
	conc.SyncWithMutex()
	conc.WaitForResult()
	conc.FanOut()
	conc.Pooling()
	conc.FanOutSemaphore()
	conc.FanOutBounded()
	conc.Drop()
	conc.Cancellation()
}
