package main

import (
	. "../socket"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func processmsg() error {

	for {

		uwd := <-Userwalkdata_chan
		fmt.Println("uid upload msg : ", uwd)

	}
	return nil
}

func init() {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sigc
		switch sig {
		case syscall.SIGINT:
			fmt.Println("catch SIGINT ")
			os.Exit(1)
		case syscall.SIGQUIT:
			fmt.Println("handle SIGQUIT")
			os.Exit(1)
		}
	}()
}

func main() {

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	go func() {
		processmsg()
	}()

	select {}
}
