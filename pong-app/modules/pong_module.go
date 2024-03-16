package modules

import (
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/module"
)

// PongModule
func PongModule() *module.Module {
	m := module.New("pingpong")

	m.AddAction("pong", func(payload structs.Message[[]byte]) (types.Response, error) {
		log.Warning("received payload from: " + payload.FromID)
		log.Warning("  -> sending response: 'pong' to: " + payload.FromID)

		return "pong", nil
	})

	return m
}
