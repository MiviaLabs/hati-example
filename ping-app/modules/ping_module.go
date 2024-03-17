package modules

import (
	"fmt"
	"sync"
	"time"

	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/module"
	"github.com/MiviaLabs/hati/transport"
)

// PingModule sends ping message every 1 second
func PingModule() *interfaces.Module {
	var m interfaces.Module = module.New("ping-module")

	var wg sync.WaitGroup
	var closeChan chan bool = make(chan bool)

	// callback to call before module starts
	m.BeforeStart(func(m interfaces.Module) {
		wg.Add(1)

		go func(w *sync.WaitGroup, cc chan bool, mod interfaces.Module) {
			ticker := time.NewTicker(25 * time.Microsecond)
			ticker2 := time.NewTicker(5 * time.Second)
			counter := 0

			for {
				select {
				case <-ticker2.C:
					fmt.Println(counter)
					continue
				case <-ticker.C:
					tm := mod.GetTransportManager()
					if tm == nil {
						log.Error("transport manager is nil")
						continue
					}

					_, err := (tm).Send(transport.REDIS_TYPE, "pong-app", "pingpong", "pong", []byte(`ping`), true, "")
					if err != nil {
						log.Error(err.Error())
					}
					counter++
					// fmt.Println("response from pong:")
					// fmt.Println(response)

				case <-cc:
					w.Done()
				}
			}
		}(&wg, closeChan, m)
	})

	// callback to call before module stops
	m.BeforeStop(func(m interfaces.Module) {
		closeChan <- true

		wg.Wait()
	})

	return &m
}
