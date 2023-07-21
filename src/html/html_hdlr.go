package html

import (
	"gen8id-websocket/src/util"
	"github.com/kataras/blocks"
	"log"
	"net/http"
)

var Views *blocks.Blocks

func BaseHtmlWithMetaTag(w http.ResponseWriter, r *http.Request) {

	log.Println("[html] meta gen page called")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"OgUrl":         util.GetConfig().OgUrl,
		"OgSiteName":    util.GetConfig().OgSiteName,
		"OgType":        util.GetConfig().OgType,
		"OgTitle":       util.GetConfig().OgTitle,
		"OgDescription": util.GetConfig().OgDescription,
		"OgImage":       util.GetConfig().OgImage,
		"OgImageType":   util.GetConfig().OgImageType,
		"OgImageWidth":  util.GetConfig().OgImageWidth,
		"OgImageHeight": util.GetConfig().OgImageHeight,
		"OgLocale":      util.GetConfig().OgLocale,

		"Keywords":  util.GetConfig().Keywords,
		"Author":    util.GetConfig().Author,
		"Copyright": util.GetConfig().Copyright,
	}

	err := Views.ExecuteTemplate(w, "index", "", data)
	if err != nil {
		log.Println(err.Error())
	}
}
