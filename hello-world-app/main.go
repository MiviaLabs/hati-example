package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MiviaLabs/hati"
	"github.com/MiviaLabs/hati-test/hello-world-app/modules"
	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/core"
	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/transport"
)

func main() {
	hati := hati.New(core.Config{
		Name: "hello-world-app",
		Transport: transport.TransportManagerConfig{
			Redis: transport.RedisConfig{
				On:       true,
				Host:     "localhost",
				Port:     "6379",
				Username: "",
				Password: "",
				Database: 0,
				Protocol: 3,
				PoolSize: 40,
			},
		},
	})

	var m interfaces.Module = modules.HelloWorldModule()

	if err := hati.AddModule(m); err != nil {
		panic(err)
	}

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
