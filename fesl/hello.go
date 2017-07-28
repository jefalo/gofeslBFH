package fesl

import (
	"../GameSpy"
	"../core"
	"../log"
)

func (fM *FeslManager) hello(event GameSpy.EventClientTLSCommand) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}

	redisState := new(core.RedisState)
	redisState.New(fM.redis, event.Command.Message["clientType"]+"-"+event.Client.IpAddr.String())

	event.Client.RedisState = redisState

	if !fM.server {
		getSession := make(map[string]string)
		getSession["TXN"] = "GetSessionId"
		event.Client.WriteFESL("gsum", getSession, 0)
	}

	saveRedis := make(map[string]interface{})
	saveRedis["SDKVersion"] = event.Command.Message["SDKVersion"]
	saveRedis["clientPlatform"] = event.Command.Message["clientPlatform"]
	saveRedis["clientString"] = event.Command.Message["clientString"]
	saveRedis["clientType"] = event.Command.Message["clientType"]
	saveRedis["clientVersion"] = event.Command.Message["clientVersion"]
	saveRedis["locale"] = event.Command.Message["locale"]
	saveRedis["sku"] = event.Command.Message["sku"]
	event.Client.RedisState.SetM(saveRedis)

	helloPacket := make(map[string]string)
	helloPacket["TXN"] = "Hello"
	helloPacket["domainPartition.domain"] = "eagames"
	if fM.server {
		helloPacket["domainPartition.subDomain"] = "bfwest-server"
	} else {
		helloPacket["domainPartition.subDomain"] = "bfwest-dedicated"
	}
	helloPacket["curTime"] = "Jun-30-2017 00:00:00 UTC"
	helloPacket["activityTimeoutSecs"] = "10"
	helloPacket["messengerIp"] = "127.0.0.1" //changehere
	helloPacket["messengerPort"] = "13505"
	helloPacket["theaterIp"] = "127.0.0.1" //changehere
	if fM.server {
		helloPacket["theaterPort"] = "18056"
	} else {
		helloPacket["theaterPort"] = "18275"
	}
	event.Client.WriteFESL("fsys", helloPacket, 0xC0000001)
	fM.logAnswer("fsys", helloPacket, 0xC0000001)

}
