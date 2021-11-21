package main

//Program struct
type Program struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

//Programs struct
type Programs struct {
	VideoFfmpeg []Program `json:"videoffmpeg"`
	VideoCache  []Program `json:"videocache"`
}
