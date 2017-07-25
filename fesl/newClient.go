package fesl

import (
	"time"

	"../GameSpy"
	"../log"
)

func (fM *FeslManager) newClient(event GameSpy.EventNewClientTLS) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}

	memCheck := make(map[string]string)
	memCheck["TXN"] = "MemCheck"
	memCheck["memcheck.[]"] = "0"
	memCheck["salt"] = "5"
	event.Client.WriteFESL("fsys", memCheck, 0xC0000000)
	fM.logAnswer("fsys", memCheck, 0xC0000000)

	// Start Heartbeat
	event.Client.State.HeartTicker = time.NewTicker(time.Second * 10)
	go func() {
		for {
			if !event.Client.IsActive {
				return
			}
			select {
			case <-event.Client.State.HeartTicker.C:
				if !event.Client.IsActive {
					return
				}
				memCheck := make(map[string]string)
				memCheck["TXN"] = "MemCheck"
				memCheck["memcheck.[]"] = "0"
				memCheck["salt"] = "5"
				event.Client.WriteFESL("fsys", memCheck, 0xC0000000)
				fM.logAnswer("fsys", memCheck, 0xC0000000)
			}
		}
	}()

	log.Noteln("Client connecting")

}
