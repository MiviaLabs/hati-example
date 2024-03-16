package transport

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
)

var (
	ErrInvalidTransportType = errors.New("invalid transport type")
)

const (
	REDIS_TYPE types.TransportType = "redis"
)

type TransportManager struct {
	modules       map[string]interfaces.Module
	serverName    string
	config        TransportManagerConfig
	redis         *Redis
	moduleManager interfaces.ModuleManager
}

type TransportManagerConfig struct {
	Redis RedisConfig `yaml:"redis" json:"redis"`
}

func NewTransportManager(serverName string, config TransportManagerConfig, moduleManager interfaces.ModuleManager) TransportManager {
	tm := TransportManager{
		serverName:    serverName,
		config:        config,
		redis:         NewRedis(config.Redis),
		moduleManager: moduleManager,
	}
	moduleManager.SetTransportManager(tm)

	return tm
}

func (tm TransportManager) SetModules(modules map[string]interfaces.Module) {
	tm.modules = modules
}

func (tm TransportManager) Start() error {
	if tm.config.Redis.On {
		tm.redis.Start()

		if err := tm.Subscribe(REDIS_TYPE, types.CHAN_MESSAGE, tm.ReceiveMessage); err != nil {
			return err
		}

		if err := tm.Subscribe(REDIS_TYPE, types.CHAN_MESSAGE_RESPONSE, tm.ReceiveMessageResponse); err != nil {
			return err
		}
	}

	return nil
}

func (tm TransportManager) Stop() error {
	log.Debug("stopping transport manager")

	if tm.config.Redis.On {
		if err := tm.redis.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (tm TransportManager) Send(transportType types.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) (any, error) {
	err := tm.Publish(transportType, types.CHAN_MESSAGE, targetServer, targetModule, targetAction, payload, waitForResponse)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (tm TransportManager) Publish(transportType types.TransportType, channel types.Channel, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) error {
	switch transportType {
	case REDIS_TYPE:
		{
			msg, err := tm.prepareMessage(targetServer, targetModule, targetAction, payload, waitForResponse)
			if err != nil {
				return err
			}

			msgBytes, err := msg.MarshalMessage()
			if err != nil {
				return err
			}

			if err := tm.redis.Publish(channel, msgBytes); err != nil {
				return err
			}
			return nil
		}
	default:
		return ErrInvalidTransportType
	}
}

func (tm TransportManager) Subscribe(transportType types.TransportType, channel types.Channel, callback func(payload []byte) (types.Response, error)) error {
	switch transportType {
	case REDIS_TYPE:
		{
			return tm.redis.Subscribe(channel, callback)
		}
	default:
		return ErrInvalidTransportType
	}
}

func (tm TransportManager) ReceiveMessage(payload []byte) (types.Response, error) {
	var message *structs.Message[[]byte] = &structs.Message[[]byte]{}

	err := json.Unmarshal(payload, message)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if message.TargetID != tm.serverName {
		return nil, errors.New("i am not the target")
	}

	module := (tm.moduleManager).GetModule(message.TargetAction.Module)
	if module == nil {
		log.Warning("module does not exist")
		return nil, errors.New("module does not exist")
	}

	response, err := module.CallAction(message.TargetAction.Action, message)
	if err != nil {
		log.Warning(err.Error())

		return nil, err
	}

	// send response
	fmt.Println(response)

	return nil, nil
}

func (tm TransportManager) ReceiveMessageResponse(payload []byte) (types.Response, error) {
	fmt.Println("<--- RECEIVE MESSAGE RESPONSE")
	fmt.Println(string(payload))

	return nil, nil
}

func (tm TransportManager) prepareMessage(targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) (*structs.Message[[]byte], error) {
	msg := &structs.Message[[]byte]{
		FromID:   tm.serverName,
		TargetID: targetServer,
		TargetAction: structs.TargetAction{
			Module: targetModule,
			Action: targetAction,
		},
		Payload:         payload,
		WaitForResponse: waitForResponse,
	}

	return msg, nil
}
