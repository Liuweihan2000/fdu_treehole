package main

import (
	"GoProject/fudan_bbs/internal/di"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_, err := di.InitApp()
	if err != nil {
		fmt.Println(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Printf("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
