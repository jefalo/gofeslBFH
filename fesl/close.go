package fesl

import (
	"../GameSpy"
	"../log"
)

//close connection
func (fM *FeslManager) close(event GameSpy.EventClientTLSClose) {
	log.Noteln("Client closed.")

	if event.Client.RedisState != nil {
		event.Client.RedisState.Delete()
	}

	if !event.Client.State.HasLogin {
		return
	}

}
