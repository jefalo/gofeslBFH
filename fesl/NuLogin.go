package fesl

import (
	"strconv"

	"../GameSpy"
	"../log"
)

// NuLogin - master login command
func (fM *FeslManager) NuLogin(event GameSpy.EventClientTLSCommand) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}

	if event.Client.RedisState.Get("clientType") == "server" {
		// Server login
		stmt, err := fM.db.Prepare("SELECT id, name, id FROM west_heroes_servers WHERE secretKey = ?")
		defer stmt.Close()
		if err != nil {
			log.Debugln(err)
			return
		}

		var sID, uID int
		var username string

		err = stmt.QueryRow(event.Command.Message["password"]).Scan(&sID, &username, &uID)
		if err != nil {
			log.Debugf(event.Command.Message["password"])
			loginPacket := make(map[string]string)
			loginPacket["TXN"] = "NuLogin"
			loginPacket["localizedMessage"] = "\"The password the user specified is incorrect\""
			loginPacket["errorContainer.[]"] = "0"
			loginPacket["errorCode"] = "122"
			event.Client.WriteFESL(event.Command.Query, loginPacket, event.Command.PayloadID)
			return
		}

		saveRedis := make(map[string]interface{})
		saveRedis["uID"] = strconv.Itoa(uID)
		saveRedis["username"] = username
		saveRedis["apikey"] = event.Command.Message["encryptedInfo"]
		saveRedis["keyHash"] = event.Command.Message["password"]
		event.Client.RedisState.SetM(saveRedis)

		loginPacket := make(map[string]string)
		loginPacket["TXN"] = "NuLogin"
		loginPacket["profileId"] = strconv.Itoa(uID)
		loginPacket["userId"] = strconv.Itoa(uID)
		loginPacket["nuid"] = username
		loginPacket["lkey"] = event.Command.Message["password"]
		event.Client.WriteFESL(event.Command.Query, loginPacket, event.Command.PayloadID)
		fM.logAnswer(event.Command.Query, loginPacket, event.Command.PayloadID)
		return
	}

	stmt, err := fM.db.Prepare("SELECT t1.uid, t1.sessionid, t1.ip, t2.username, t2.banned, t2.is_admin, t2.is_tester, t2.confirmed_em, t2.key_hash, t2.email, t2.country FROM web_sessions t1 LEFT JOIN web_users t2 ON t1.uid=t2.id WHERE t1.sessionid = ?")
	defer stmt.Close()
	if err != nil {
		log.Debugln(err)
		return
	}

	var uID int
	var ip, username, sessionID, keyHash, email, country string
	var banned, isAdmin, isTester, confirmedEm bool

	err = stmt.QueryRow(event.Command.Message["encryptedInfo"]).Scan(&uID, &sessionID, &ip, &username, &banned, &isAdmin, &isTester, &confirmedEm, &keyHash, &email, &country)
	if err != nil {
		loginPacket := make(map[string]string)
		loginPacket["TXN"] = "NuLogin"
		loginPacket["localizedMessage"] = "\"The password the user specified is incorrect\""
		loginPacket["errorContainer.[]"] = "0"
		loginPacket["errorCode"] = "122"
		event.Client.WriteFESL(event.Command.Query, loginPacket, event.Command.PayloadID)
		return
	}

	// ready for all
	if sessionID != event.Command.Message["encryptedInfo"] || !confirmedEm || banned {
		log.Noteln("User not worthy: " + username)
		loginPacket := make(map[string]string)
		loginPacket["TXN"] = "NuLogin"
		loginPacket["localizedMessage"] = "\"The user is not entitled to access this game\""
		loginPacket["errorContainer.[]"] = "0"
		loginPacket["errorCode"] = "120"
		event.Client.WriteFESL(event.Command.Query, loginPacket, event.Command.PayloadID)
		return
	}

	saveRedis := make(map[string]interface{})
	saveRedis["uID"] = strconv.Itoa(uID)
	saveRedis["username"] = username
	saveRedis["ip"] = ip
	saveRedis["sessionID"] = sessionID
	saveRedis["keyHash"] = keyHash
	saveRedis["email"] = email
	saveRedis["country"] = country
	event.Client.RedisState.SetM(saveRedis)

	loginPacket := make(map[string]string)
	loginPacket["TXN"] = "NuLogin"
	loginPacket["profileId"] = strconv.Itoa(uID)
	loginPacket["userId"] = strconv.Itoa(uID)
	loginPacket["nuid"] = username
	loginPacket["lkey"] = keyHash
	event.Client.WriteFESL(event.Command.Query, loginPacket, event.Command.PayloadID)
	fM.logAnswer(event.Command.Query, loginPacket, event.Command.PayloadID)
}
