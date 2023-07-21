package main

import (
	"gen8id-websocket/src/rest"
	"gen8id-websocket/src/utils"
	"gen8id-websocket/src/ws"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/ws", ws.StreamUpload)
	router.PathPrefix("/img").HandlerFunc(rest.BaseHtmlWithMetaTag)
	err = http.ListenAndServe(conf.ServerPort, router)
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
