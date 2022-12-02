package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() <-chan struct{} {
	stopSignals := make(chan os.Signal, 2)
	stop := make(chan struct{})

	signal.Notify(stopSignals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stopSignals
		close(stop)
		<-stopSignals
		os.Exit(1)
	}()

	return stop
}
