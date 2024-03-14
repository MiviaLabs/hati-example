package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MiviaLabs/hati"
	"github.com/MiviaLabs/hati/core"
	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/module"
	"github.com/MiviaLabs/hati/transport"
)

func main() {
	hati := hati.New(core.Config{
		Name: "example-app",
		Transport: transport.TransportManagerConfig{
			Redis: transport.RedisConfig{
				On:       true,
				Host:     "localhost",
				Port:     "6379",
				Username: "",
				Password: "",
				Database: 0,
			},
		},
	})

	m := module.New("test")

	hati.AddModule(m)

	if err := hati.Start(); err != nil {
		panic(err)
	}

	var osSignal chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-osSignal:
			log.Warning("shutting down, please wait")

			if err := hati.Stop(); err != nil {
				panic(err)
			}

			os.Exit(0)
		}
	}
}
