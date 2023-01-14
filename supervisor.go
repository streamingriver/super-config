package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Supervisor generates real config files
type Supervisor struct {
	App    *App
	Logger *log.Logger
	Output io.WriteCloser
}

func (sv Supervisor) templateFFMPEG(name, url string) string {
	template := `[program:{{name}}]
command=/bin/bash -c "mkdir -p /dev/shm/{{name}}; cd /dev/shm/{{name}}; /ffmpeg -nostats -nostdin -user-agent "streamingriveriptv/1.0" -i "{{url}}" -codec copy -map 0:0 -map 0:1 -map_metadata 0  -f hls -hls_list_size 3 -hls_flags delete_segments -hls_time 5 -segment_list_size 3 -hls_segment_filename file%%07d.ts stream.m3u8"
autostart = true
startsec = 1
user = root
stdout_logfile=/dev/fd/1
stdout_logfile_maxbytes=0
stderr_logfile=/dev/fd/2
stderr_logfile_maxbytes=0
autorestart=true
startretries=5000000000
stopasgroup=true
killasgroup=true
stdout_events_enabled=true
stderr_events_enabled=true

`

	template = strings.ReplaceAll(template, "{{name}}", name)
	template = strings.ReplaceAll(template, "{{url}}", url)
	return template
}

func (sv Supervisor) templateCache(name, url string) string {
	template := `[program:proxy-{{name}}]
command=/opt/tools/hls-proxy_linux_amd64 --url "http://localhost:9005/{{name}}/stream.m3u8" --name "{{name}}" --frontend http://127.0.0.1:8085
autostart = true
startsec = 1
user = root
stdout_logfile=/dev/fd/1
stdout_logfile_maxbytes=0
stderr_logfile=/dev/fd/2
stderr_logfile_maxbytes=0
autorestart=true
startretries=5000000000
stopasgroup=true
killasgroup=true
stdout_events_enabled=true
stderr_events_enabled=true

`

	template = strings.ReplaceAll(template, "{{name}}", name)
	template = strings.ReplaceAll(template, "{{url}}", url)
	return template
}

// GenerateFFMPEG file
func (sv *Supervisor) GenerateFFMPEG(programs []Program) {
	sv.generate(programs, "ffmpeg")
}

// GenerateCache file
func (sv *Supervisor) GenerateCache(programs []Program) {
	sv.generate(programs, "cache")
}

func (sv *Supervisor) generate(programs []Program, t string) {

	output := ""
	for _, program := range programs {
		if t == "ffmpeg" {
			output += sv.templateFFMPEG(program.Name, program.URL)
		} else {
			output += sv.templateCache(program.Name, program.URL)
		}
	}

	if sv.Output != nil {
		sv.Output.Write([]byte(output))
		return
	}

	err := os.WriteFile(sv.App.path+"/"+t+sv.App.ext, []byte(output), 0755)
	if err != nil {
		sv.Logger.Printf("%s", err)
	}
	sv.reload()
}

func (sv *Supervisor) reload() {
	cmd := exec.Command(sv.App.SupervisorPath, "-c", sv.App.SupervisorConfig, "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%v (%s)", err, output)
	}
}
