package main

type AppImage struct {
	Thumbnail string `json:"thumbnail"`
	Original  string `json:"original"`
}


func main() {

	cfg := InitConfig()

	server := NewAppServer(cfg)
	server.Run()
}
