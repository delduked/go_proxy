package utils

import "prx/internal/config"

var semaphore = make(chan struct{}, config.MaxRequests)

func Acquire() {
	semaphore <- struct{}{}
}

func Release() {
	<-semaphore
}
