package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"gitlab.com/avarf/getenvs"
)

// App is main struct
type App struct {
	natsURL   string
	natsTopic string

	logger *log.Logger
	conn   *nats.Conn

	parser           *Parser
	parserType       string
	path             string
	ext              string
	SupervisorPath   string
	SupervisorConfig string

	testing bool
	exit    chan struct{}
}

func (app *App) initLogger() {
	app.logger = log.Default()
}

func (app *App) getEnvVariables() {
	app.natsURL = getenvs.GetEnvString("SC_NATS_URL", "nats://nats:4222")
	app.natsTopic = getenvs.GetEnvString("SC_NATS_TOPIC", "super-config")
	app.parserType = getenvs.GetEnvString("SC_PARSER_TYPE", "ffmpeg")
	app.path = getenvs.GetEnvString("SC_CONFIG_PATH", "/etc/supervisor/conf.d/")
	app.ext = getenvs.GetEnvString("SC_CONFIG_EXT", ".conf")
	app.SupervisorPath = getenvs.GetEnvString("SC_SUPERVISORCTL", "supervisorctl")
	app.SupervisorConfig = getenvs.GetEnvString("SC_SUPERVISOR_CONFIG", "/data/conf/supervisor.conf")
}

//ConnectToNats server
func (app *App) ConnectToNats() {
	var err error
	app.conn, err = nats.Connect(app.natsURL)
	if err != nil {
		if app.testing {
			app.logger.Printf("%v", err)
			return
		}
		app.logger.Fatalf("%v", err)
	}
}

// Subscribe to nats channel
func (app *App) Subscribe() {
	_, err := app.conn.Subscribe(app.natsTopic, func(msg *nats.Msg) {
		app.parser.Parse(msg.Data)
	})
	if err != nil {
		if app.testing {
			app.logger.Printf("%v", err)
			return
		}
		app.logger.Fatalf("%v", err)
	}
}

func (app *App) init() {
	app.getEnvVariables()
	app.initLogger()
	app.exit = make(chan struct{})

	app.parser = &Parser{app, app.logger, app.parserType, nil}
}

// Run starts application
func (app *App) Run() {
	app.init()

	app.ConnectToNats()
	app.Subscribe()

	<-app.exit
}
