package fesl

import (
	"../GameSpy"
	"../log"
)

// NuGetAccount - General account information retrieved, based on parameters sent
func (fM *FeslManager) NuGetAccount(event GameSpy.EventClientTLSCommand) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}

	loginPacket := make(map[string]string)
	loginPacket["TXN"] = "NuGetAccount"
	loginPacket["heroName"] = event.Client.RedisState.Get("username")
	loginPacket["nuid"] = event.Client.RedisState.Get("username") + "@westheroes.com"
	loginPacket["DOBDay"] = "1"
	loginPacket["DOBMonth"] = "1"
	loginPacket["DOBYear"] = "2017"
	loginPacket["userId"] = event.Client.RedisState.Get("uID")
	loginPacket["globalOptin"] = "0"
	loginPacket["thidPartyOptin"] = "0"
	loginPacket["language"] = "enUS"
	loginPacket["country"] = "US"
	event.Client.WriteFESL(event.Command.Query, loginPacket, event.Command.PayloadID)
	fM.logAnswer(event.Command.Query, loginPacket, event.Command.PayloadID)
}
