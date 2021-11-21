package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	app := new(App)

	os.Setenv("NATS_URL", "nats-server")
	os.Setenv("NATS_TOPIC", "nats-topic")
	os.Setenv("PARSER_TYPE", "cache")
	os.Setenv("CONFIG_PATH", "etc")
	os.Setenv("CONFIG_EXT", ".conf")

	app.init()

	if app.natsTopic != "nats-topic" {
		t.Errorf("Wrong nats topic %s", app.natsTopic)
	}
	if app.natsURL != "nats-server" {
		t.Errorf("Wrong nats server %s", app.natsURL)
	}
	if app.parserType != "cache" {
		t.Errorf("Wrong parser type %s", app.parserType)
	}
	if app.path != "etc" {
		t.Errorf("Wrong parser type %s", app.parserType)
	}
	if app.ext != ".conf" {
		t.Errorf("Wrong parser type %s", app.parserType)
	}
}

type LogWriter struct {
	t    *testing.T
	buff *bytes.Buffer
}

func (mw LogWriter) Write(b []byte) (int, error) {
	mw.buff.Write(b)
	return len(b), nil
}
func TestConnectToNatsShouldFail(t *testing.T) {
	os.Setenv("NATS_URL", "127.0.0.1:4567")
	os.Setenv("NATS_TOPIC", "nats-topic")
	os.Setenv("PARSER_TYPE", "cache")

	writer := &LogWriter{t, bytes.NewBuffer(nil)}
	app := new(App)
	app.init()
	app.testing = true
	app.logger = log.New(writer, "", 0)
	go func() {
		app.ConnectToNats()
	}()
	time.Sleep(time.Second)

	if strings.Trim(writer.buff.String(), "\n") != "nats: no servers available for connection" {
		t.Errorf("ConnectToNats %s", writer.buff.Bytes())
	}
}

func TestNatsSubscribeShouldFail(t *testing.T) {
	os.Setenv("NATS_URL", "127.0.0.1:4567")
	os.Setenv("NATS_TOPIC", "nats-topic")
	os.Setenv("PARSER_TYPE", "cache")

	writer := &LogWriter{t, bytes.NewBuffer(nil)}
	app := new(App)
	app.testing = true
	app.init()
	app.logger = log.New(writer, "", 0)
	go func() {
		app.Subscribe()
	}()
	time.Sleep(time.Second)

	if strings.Trim(writer.buff.String(), "\n") != "nats: invalid connection" {
		t.Errorf("Subscribe %s", writer.buff.Bytes())
	}
}
