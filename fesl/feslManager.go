package fesl

import (
	"database/sql"
	"strings"
	"time"

	"../log"

	"../GameSpy"

	"../matchmaking"
	"github.com/go-redis/redis"
)

// FeslManager - handles incoming and outgoing FESL data
type FeslManager struct {
	name          string
	db            *sql.DB
	redis         *redis.Client
	socket        *GameSpy.SocketTLS
	eventsChannel chan GameSpy.SocketEvent
	batchTicker   *time.Ticker
	stopTicker    chan bool
	server        bool
}

// New creates and starts a new ClientManager
func (fM *FeslManager) New(name string, port string, certFile string, keyFile string, server bool, db *sql.DB, redis *redis.Client) {
	var err error

	fM.socket = new(GameSpy.SocketTLS)
	fM.db = db
	fM.redis = redis
	fM.name = name
	fM.eventsChannel, err = fM.socket.New(fM.name, port, certFile, keyFile)
	fM.stopTicker = make(chan bool, 1)
	fM.server = server

	if err != nil {
		log.Errorln(err)
	}

	go fM.run()
}

func (fM *FeslManager) run() {
	for {
		select {
		case event := <-fM.eventsChannel:
			log.Debugf(event.Name)
			switch {
			case event.Name == "newClient":
				fM.newClient(event.Data.(GameSpy.EventNewClientTLS))
			case event.Name == "client.command.Hello":
				fM.hello(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.NuLogin":
				fM.NuLogin(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.NuGetPersonas":
				fM.NuGetPersonas(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.NuGetAccount":
				fM.NuGetAccount(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.NuLoginPersona":
				fM.NuLoginPersona(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.GetStatsForOwners":
				fM.GetStatsForOwners(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.GetStats":
				fM.GetStats(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.NuLookupUserInfo":
				fM.NuLookupUserInfo(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.GetPingSites":
				fM.GetPingSites(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.UpdateStats":
				fM.UpdateStats(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.GetTelemetryToken":
				fM.GetTelemetryToken(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.command.Start":
				fM.Start(event.Data.(GameSpy.EventClientTLSCommand))
			case event.Name == "client.close":
				fM.close(event.Data.(GameSpy.EventClientTLSClose))
			case event.Name == "client.command":
				fM.LogCommand(event.Data.(GameSpy.EventClientTLSCommand))
				log.Debugf("Got event %s.%s: %v", event.Name, event.Data.(GameSpy.EventClientTLSCommand).Command.Message["TXN"], event.Data.(GameSpy.EventClientTLSCommand).Command)
			default:
				log.Debugf("Got event %s: %v", event.Name, event.Data)
			}
		}
	}
}

// LogCommand - logs detailed FESL command data to a file for further analysis
func (fM *FeslManager) LogCommand(event GameSpy.EventClientTLSCommand) {
	/* not necessary
	b, err := json.MarshalIndent(event.Command.Message, "", "	")
	if err != nil {
		panic(err)
	}

	commandType := "request"

	os.MkdirAll("./commands/"+event.Command.Query+"."+event.Command.Message["TXN"]+"", 0777)
	err = ioutil.WriteFile("./commands/"+event.Command.Query+"."+event.Command.Message["TXN"]+"/"+commandType, b, 0644)
	if err != nil {
		panic(err)
	}
	*/
}

func (fM *FeslManager) logAnswer(msgType string, msgContent map[string]string, msgType2 uint32) {
	/* not necessary
	b, err := json.MarshalIndent(msgContent, "", "	")
	if err != nil {
		panic(err)
	}

	commandType := "answer"

	os.MkdirAll("./commands/"+msgType+"."+msgContent["TXN"]+"", 0777)
	err = ioutil.WriteFile("./commands/"+msgType+"."+msgContent["TXN"]+"/"+commandType, b, 0644)
	if err != nil {
		panic(err)
	}*/
}

// Status - Basic fesl call to get overall service status (called before pnow?)
func (fM *FeslManager) Status(event GameSpy.EventClientTLSCommand) {
	if !event.Client.IsActive {
		log.Noteln("Client left")
		return
	}

	log.Noteln("STATUS CALLED")

	answer := make(map[string]string)
	answer["TXN"] = "Status"
	answer["id.id"] = "1"
	answer["id.partition"] = event.Command.Message["partition.partition"]
	answer["sessionState"] = "COMPLETE"
	answer["props.{}.[]"] = "2"
	answer["props.{resultType}"] = "JOIN"

	// Find latest game (do better later)
	gameID := matchmaking.FindAvailableGID()

	answer["props.{games}.0.lid"] = "1"
	answer["props.{games}.0.fit"] = "1001"
	answer["props.{games}.0.gid"] = gameID
	answer["props.{games}.[]"] = "1"

	event.Client.WriteFESL("pnow", answer, 0x80000000)
	fM.logAnswer("pnow", answer, 0x80000000)
}

// MysqlRealEscapeString - you know
func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

func (fM *FeslManager) error(event GameSpy.EventClientTLSError) {
	log.Noteln("Client threw an error: ", event.Error)
}
