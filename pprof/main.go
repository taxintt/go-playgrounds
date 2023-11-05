package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var s = &struct {
		st string
	}{
		st: "test",
	}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			select {
			case <-sigs:
				fmt.Println("signal has come")
				done <- struct{}{}
				break
			default:
				fmt.Printf("working with %+v\n", s)
			}
		}
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
