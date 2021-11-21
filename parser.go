package main

import (
	"encoding/json"
	"io"
	"log"
)

//Parser of nats messages
type Parser struct {
	Logger *log.Logger
	Type   string
	Output io.WriteCloser
	Path   string
}

// Parse nats message and take action
func (p *Parser) Parse(b []byte) {
	var programs Programs
	err := json.Unmarshal(b, &programs)
	if err != nil {
		p.Logger.Printf("Parser: %v", err)
	}

	if p.Type == "ffmpeg" {
		(&Supervisor{Logger: p.Logger, Output: p.Output, Path: p.Path}).GenerateFFMPEG(programs.VideoFfmpeg)
	} else {
		(&Supervisor{Logger: p.Logger, Output: p.Output, Path: p.Path}).GenerateCache(programs.VideoCache)
	}
}
