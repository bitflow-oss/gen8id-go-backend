package ws

import (
	"bytes"
	"fmt"
	"gen8id-websocket/src/cnst"
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
			var strMsg = streamToByte(reader)
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

	var orgFilePath = filepath.Join(cnst.UPLOAD_REL_PATH,
		fmt.Sprintf(cnst.ORG_IMG_FILENAME, time.Now().UnixMilli()))

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
	var hashedFilename = fmt.Sprintf(cnst.HASH_IMG_FILENAME, fileHash)
	var hashedFilePath = filepath.Join(cnst.UPLOAD_REL_PATH, hashedFilename)
	err = os.Rename(orgFilePath, hashedFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var imgUrl = util.ObjectPrivateUpload(cnst.UPLOAD_REL_PATH, hashedFilename)
	// imgUrl = utils.GenerateThumbnailWithWatermark(gloval_consts.ORG_IMG_FILENAME, fileHash)
	// log.Printf("image saved to %s, uploaded to %s\n", fileHash, imgUrl)
	return imgUrl, nil
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf.Bytes()
}
