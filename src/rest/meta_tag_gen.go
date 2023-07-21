package rest

import (
	"gen8id-websocket/src/utils"
	"github.com/kataras/blocks"
	"log"
	"net/http"
)

var conf = utils.GetConfig()
var Views *blocks.Blocks

func BaseHtmlWithMetaTag(w http.ResponseWriter, r *http.Request) {

	log.Println("[html] meta gen page called")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"og:url":          conf.OgUrl,
		"og:site_name":    conf.OgSiteName,
		"og:type":         conf.OgType,
		"og:title":        conf.OgTitle,
		"og:description":  conf.OgDescription,
		"og:image":        conf.OgImage,
		"og:image:type":   conf.OgImageType,
		"og:image:width":  conf.OgImageWidth,
		"og:image:height": conf.OgImageHeight,
		"og:locale":       conf.OgLocale,

		"keywords":  conf.Keywords,
		"author":    conf.Author,
		"copyright": conf.Copyright,
	}

	err := Views.ExecuteTemplate(w, "index", "", data)
	if err != nil {
		log.Println(err.Error())
	}
}
