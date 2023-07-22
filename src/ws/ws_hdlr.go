package ws

import (
	"fmt"
	"gen8id-websocket/src/extn"
	"gen8id-websocket/src/util"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

/*
TextMessage = 1
// BinaryMessage denotes a binary data message.
BinaryMessage = 2
// CloseMessage denotes a close control message. The optional message
// payload contains a numeric code and text. Use the FormatCloseMessage
// function to format a close message payload.
CloseMessage = 8
// PingMessage denotes a ping control message. The optional message payload
// is UTF-8 encoded text.
PingMessage = 9
// PongMessage denotes a pong control message. The optional message payload
// is UTF-8 encoded text.
PongMessage = 10
*/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		return true // origin == "::1"
	},
}

// StreamUpload github.com/chai2010/webp
func StreamUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("upload called")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(conn)

	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}

		if messageType == websocket.BinaryMessage {
			start := time.Now()
			imgUrl, err := saveBinaryMessage(reader)
			if err != nil {
				log.Println(err)
				return
			}
			elapsed := time.Since(start)
			secs := elapsed.Seconds()
			log.Printf("Frame rate = %f\n", 256.0/secs)
			// return : image file URL, width x height, file type : png, jpg, webp, 사이즈 적합여부
			err = conn.WriteMessage(websocket.TextMessage, []byte(imgUrl))

		} else if messageType == websocket.TextMessage {
			// echo message
			var strMsg = util.StreamToByte(reader)
			log.Println("got msg:", string(strMsg))
			err = conn.WriteMessage(messageType, strMsg)
			if err != nil {
				log.Println("err:", err)
				break
			}

		}
	}
}

func saveBinaryMessage(reader io.Reader) (string, error) {

	var conf = util.GetConfig()

	var orgFilePath = filepath.Join(conf.UpldRltvPath,
		fmt.Sprintf(conf.OrgImgFileNm, time.Now().UnixMilli()))

	initFile, err := os.Create(orgFilePath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	_, err = io.Copy(initFile, reader)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = initFile.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	var fileHash, _ = util.ExtractFileHash(initFile.Name())
	var hashedFilename = fmt.Sprintf(conf.HashImgFileNm, fileHash)
	var hashedFilePath = filepath.Join(conf.UpldRltvPath, hashedFilename)
	err = os.Rename(orgFilePath, hashedFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var imgUrl = extn.ObjectPrivateUpload(conf.UpldRltvPath, hashedFilename)
	return imgUrl, nil
}
