package main

import (
	"flag"
	"gen8id-websocket/src/html"
	"gen8id-websocket/src/util"
	"gen8id-websocket/src/ws"
	"github.com/gorilla/mux"
	"github.com/kataras/blocks"
	"log"
	"net/http"
)

func main() {

	env := flag.String("env", "", "")
	flag.Parse()

	log.Println("starting", *env, "gorilla websocket/mux server")

	var conf = util.LoadConfig("config.yml")
	html.Views = blocks.New(util.GetHtmlTemplateDir(*env)).Reload(true)
	err := html.Views.Load()
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/ws", ws.StreamUpload)
	router.PathPrefix("/img").HandlerFunc(html.BaseHtmlWithMetaTag)
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
