package main

type Program struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Programs struct {
	VideoFfmpeg []Program `json:"videoffmpeg"`
	VideoCache  []Program `json:"videocache"`
}
