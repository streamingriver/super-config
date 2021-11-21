package main

import (
	"encoding/json"
	"io"
	"log"
)

//Parser of nats messages
type Parser struct {
	App    *App
	Logger *log.Logger
	Type   string
	Output io.WriteCloser
}

// Parse nats message and take action
func (p *Parser) Parse(b []byte) {
	var programs Programs
	err := json.Unmarshal(b, &programs)
	if err != nil {
		p.Logger.Printf("Parser: %v", err)
	}

	if p.Type == "ffmpeg" {
		(&Supervisor{App: p.App, Logger: p.Logger, Output: p.Output}).GenerateFFMPEG(programs.VideoFfmpeg)
	} else {
		(&Supervisor{App: p.App, Logger: p.Logger, Output: p.Output}).GenerateCache(programs.VideoCache)
	}
}
