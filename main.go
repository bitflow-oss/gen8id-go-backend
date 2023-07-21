package main

import (
	"flag"
	"gen8id-websocket/src/rest"
	"gen8id-websocket/src/utils"
	"gen8id-websocket/src/ws"
	"github.com/kataras/blocks"
	"log"
	"net/http"
)

func main() {
	log.Println("starting gorilla websocket server")

	var conf = utils.LoadConfig("config.yml")
	rest.Views = blocks.New(conf.HtmlRootDir).Reload(true)
	err := rest.Views.Load()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/ws", ws.StreamUpload)
	http.HandleFunc("/meta", rest.BaseHtmlWithMetaTag)
	// ListenAndServeTLS
	var addr = flag.String("addr", ":"+conf.ServerPort, "http service address")
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	// todo: startup ASCII art needed
}

/*
func main() {
	utils.UpscaleTest()
}
*/
