package html

import (
	"gen8id-websocket/src/utils"
	"github.com/kataras/blocks"
	"log"
	"net/http"
)

var Views *blocks.Blocks

func BaseHtmlWithMetaTag(w http.ResponseWriter, r *http.Request) {

	log.Println("[html] meta gen page called")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"OgUrl":         utils.GetConfig().OgUrl,
		"OgSiteName":    utils.GetConfig().OgSiteName,
		"OgType":        utils.GetConfig().OgType,
		"OgTitle":       utils.GetConfig().OgTitle,
		"OgDescription": utils.GetConfig().OgDescription,
		"OgImage":       utils.GetConfig().OgImage,
		"OgImageType":   utils.GetConfig().OgImageType,
		"OgImageWidth":  utils.GetConfig().OgImageWidth,
		"OgImageHeight": utils.GetConfig().OgImageHeight,
		"OgLocale":      utils.GetConfig().OgLocale,

		"Keywords":  utils.GetConfig().Keywords,
		"Author":    utils.GetConfig().Author,
		"Copyright": utils.GetConfig().Copyright,
	}

	err := Views.ExecuteTemplate(w, "index", "", data)
	if err != nil {
		log.Println(err.Error())
	}
}
