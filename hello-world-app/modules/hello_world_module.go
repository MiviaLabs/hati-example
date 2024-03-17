package modules

import (
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/module"
)

// HelloWorldModule
func HelloWorldModule() *module.Module {
	m := module.New("helloworld")

	m.AddAction("hi", func(payload structs.Message[[]byte]) (types.Response, error) {

		return "hello", nil
	}, &structs.ActionRoute{
		Methods: []string{types.GET.String()},
		Path:    "/",
	})

	return m
}
