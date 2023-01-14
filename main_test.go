package main

import (
	"bytes"
	"log"
	"testing"
)

type FakeLogOutput struct {
	t *testing.T
}

func (log FakeLogOutput) Write(b []byte) (int, error) {
	log.t.Logf("%s", b)
	return len(b), nil
}

var programsDataFfmpeg = []byte(`
{
	"videoffmpeg" : [
		{
			"name": "test",
			"url": "http://test"
		},
		{
			"name": "anothertest1",
			"url": "http://anothertest1"
		}
	]
}
`)

var programsDataCache = []byte(`
{
	"videocache" : [
		{
			"name": "test",
			"url": "http://test"
		},
		{
			"name": "anothertest1",
			"url": "http://anothertest1"
		}
	]
}
`)

var programsDataFfmpegExpected = `[program:test]
command=/bin/bash -c "mkdir -p /dev/shm/test; cd /dev/shm/test; /ffmpeg -nostats -nostdin -user-agent "streamingriveriptv/1.0" -i "http://test" -codec copy -map 0:0 -map 0:1 -map_metadata 0  -f hls -hls_list_size 3 -hls_flags delete_segments -hls_time 5 -segment_list_size 3 -hls_segment_filename file%%07d.ts stream.m3u8"
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

[program:anothertest1]
command=/bin/bash -c "mkdir -p /dev/shm/anothertest1; cd /dev/shm/anothertest1; /ffmpeg -nostats -nostdin -user-agent "streamingriveriptv/1.0" -i "http://anothertest1" -codec copy -map 0:0 -map 0:1 -map_metadata 0  -f hls -hls_list_size 3 -hls_flags delete_segments -hls_time 5 -segment_list_size 3 -hls_segment_filename file%%07d.ts stream.m3u8"
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

var programsDataCacheExpected = `[program:proxy-test]
command=/opt/tools/hls-proxy_linux_amd64 --url "http://localhost:9005/test/stream.m3u8" --name "test" --frontend http://127.0.0.1:8085
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

[program:proxy-anothertest1]
command=/opt/tools/hls-proxy_linux_amd64 --url "http://localhost:9005/anothertest1/stream.m3u8" --name "anothertest1" --frontend http://127.0.0.1:8085
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

type MyWriter struct {
	t    *testing.T
	buff *bytes.Buffer
}

func (mw MyWriter) Write(b []byte) (int, error) {
	mw.buff.Write(b)
	return len(b), nil
}
func (mw MyWriter) Close() error { return nil }

func TestParserFFMPEG(t *testing.T) {
	writer := &MyWriter{t, bytes.NewBuffer(nil)}
	app := &App{}
	parser := &Parser{
		app,
		log.New(&FakeLogOutput{t}, "", 0),
		"ffmpeg",
		writer,
	}
	parser.Parse(programsDataFfmpeg)

	if writer.buff.String() != programsDataFfmpegExpected {
		t.Errorf("Unexpected output, %s", writer.buff.String())
	}
}

func TestParserCache(t *testing.T) {
	app := &App{}
	writer := &MyWriter{t, bytes.NewBuffer(nil)}
	parser := &Parser{
		app,
		log.New(&FakeLogOutput{t}, "", 0),
		"cache",
		writer,
	}
	parser.Parse(programsDataCache)

	if writer.buff.String() != programsDataCacheExpected {
		t.Errorf("Unexpected output,")
	}
}
