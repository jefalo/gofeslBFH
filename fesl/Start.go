package fesl

import (
	"../GameSpy"
	"../log"
)

// Start - a method of pnow
func (fM *FeslManager) Start(event GameSpy.EventClientTLSCommand) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}
	log.Noteln("START CALLED")
	log.Noteln(event.Command.Message["partition.partition"])
	answer := make(map[string]string)
	answer["TXN"] = "Start"
	answer["id.id"] = "1"
	answer["id.partition"] = event.Command.Message["partition.partition"]
	event.Client.WriteFESL(event.Command.Query, answer, event.Command.PayloadID)
	fM.logAnswer(event.Command.Query, answer, event.Command.PayloadID)

	fM.Status(event)
}
