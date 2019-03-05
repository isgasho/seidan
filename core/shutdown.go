package core

import (
	"log"
	"sync"
)

type ShutdownFunc func() error

var (
	shutdownFuncs []ShutdownFunc
	shutdownLock  sync.Mutex
)

func RegisterShutdown(f ShutdownFunc) {
	shutdownLock.Lock()
	defer shutdownLock.Unlock()
	shutdownFuncs = append(shutdownFuncs, f)
}

func Shutdown() {
	shutdownLock.Lock()
	defer shutdownLock.Unlock()

	if len(shutdownFuncs) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(shutdownFuncs))

	for _, f := range shutdownFuncs {
		go func(f ShutdownFunc) {
			err := f()
			if err != nil {
				log.Printf("[core] Shutdown error reported: %s", err)
			}
			wg.Done()
		}(f)
	}

	wg.Wait()
}
