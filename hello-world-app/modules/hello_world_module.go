package modules

import (
	"github.com/MiviaLabs/hati/module"
)

// HelloWorldModule
func HelloWorldModule() *module.Module {
	m := module.New("hello-world-module")

	// m.AddAction("ping", func(payload transport.Message[[]byte]) (any, error) {
	// 	for {
	// 		select {}
	// 	}
	// 	return nil, nil
	// })

	return m
}
