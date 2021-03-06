package fesl

import (
	"strconv"

	"../GameSpy"
	"../log"
)

// UpdateStats - updates stats about a soldier
func (fM *FeslManager) UpdateStats(event GameSpy.EventClientTLSCommand) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}

	answer := event.Command.Message
	answer["TXN"] = "UpdateStats"

	users, _ := strconv.Atoi(event.Command.Message["u.[]"])
	for i := 0; i < users; i++ {
		query := ""
		owner, ok := event.Command.Message["u."+strconv.Itoa(i)+".o"]

		if !ok {
			return
		}

		statsNum, _ := strconv.Atoi(event.Command.Message["u."+strconv.Itoa(i)+".s.[]"])
		for j := 0; j < statsNum; j++ {
			if event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".t"] != "" {
				query += event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"] + "='" + MysqlRealEscapeString(event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".t"]) + "', "
			} else if event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"] == "c_wallet_hero" || event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"] == "c_wallet_valor" {
				query += event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"] + "= " + event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"] + MysqlRealEscapeString(event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".v"]) + ", "
			} else {
				query += event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"] + "='" + MysqlRealEscapeString(event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".v"]) + "', "
			}
		}

		if owner != "0" && owner != event.Client.RedisState.Get("uID") {
			sql := "UPDATE `west_heroes_stats` SET " + query + "pid=" + owner + " WHERE pid = " + owner + ""
			_, err := fM.db.Exec(sql)
			if err != nil {
				log.Errorln(err)
			}
		} else {
			sql := "UPDATE `west_heroes_accounts` SET " + query + "uid=" + owner + " WHERE uid = " + owner + ""
			_, err := fM.db.Exec(sql)
			if err != nil {
				log.Errorln(err)
			}
		}
	}

	event.Client.WriteFESL(event.Command.Query, answer, event.Command.PayloadID)
	fM.logAnswer(event.Command.Query, answer, event.Command.PayloadID)
}
