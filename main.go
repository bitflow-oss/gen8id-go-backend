package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gen8id-websocket/utils"
	"github.com/gorilla/websocket"
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

var addr = flag.String("addr", ":8081", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		return true // origin == "::1"
	},
}

/**
 * convert to webp then upload file to cloud object storage
 * "github.com/chai2010/webp"
 */
func upload(w http.ResponseWriter, r *http.Request) {
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
			fmt.Printf("Frame rate = %f\n", 256.0/secs)
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

	file1, err := os.Create(utils.INIT_FILE_NAME)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(file1)

	_, err = io.Copy(file1, reader)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var fileHash, _ = utils.ExtractFileHash(file1.Name())
	var imgUrl = utils.GenerateThumbnailWithWatermark(utils.INIT_FILE_NAME, fileHash)

	fmt.Printf("image saved to %s, uploaded to %s\n", fileHash, imgUrl)

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

func main() {
	log.Println("starting gorilla websocket server")
	http.HandleFunc("/ws", upload)
	// ListenAndServeTLS
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	// todo: startup ASCII art needed
}
