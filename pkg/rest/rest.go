package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type RestfulServing interface {
	Run(stopCh <-chan struct{}) error
	SetReady(isReady bool)
}

type daemon struct {
	listener net.Listener
	ready    *Ready
	server   http.Server
}

func NewRestfulServing(addr string) (RestfulServing, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("new listener error: %v", err)
	}

	return &daemon{listener: listener, ready: NewReady()}, nil
}

func (d *daemon) Run(stopCh <-chan struct{}) error {
	return d.serve(stopCh)
}

func (d *daemon) serve(stopCh <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ready", d.ready.Handler)
	d.server = http.Server{Handler: mux}

	stoppedCh := make(chan struct{})
	go func() {
		defer close(stoppedCh)
		<-stopCh
		fmt.Println("Start to Shutdown rest server")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		d.server.Shutdown(ctx)
		cancel()
	}()

	msg := fmt.Sprintf("Stopped listening on %s", d.listener.Addr().String())
	go func() {
		err := d.server.Serve(d.listener)
		select {
		case <-stopCh:
			fmt.Println(msg)
		default:
			panic(fmt.Sprintf("%s due to error: %v", msg, err))
		}
	}()

	return nil
}

func (d *daemon) SetReady(isReady bool) {
	d.ready.SetReady(isReady)
}
