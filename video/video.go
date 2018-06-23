package video

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"github.com/dhowden/raspicam"
	"encoding/base64"
	"time"
	"fmt"
	"os"
	"strconv"
)

var (
	cmd = raspicam.NewStill()
	upgrader = websocket.Upgrader {
		WriteBufferSize: 1024*1024,
	}
)

type webSocketVideoWriter struct {
	conn *websocket.Conn
}

func VideoHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsWriter := webSocketVideoWriter{
		conn: conn,
	}
	for {
		errCh := make(chan error)
		go writeErrors(errCh)
		raspicam.Capture(cmd, wsWriter, errCh)
		time.Sleep(time.Second)
	}
}

func (writer webSocketVideoWriter) Write(p []byte) (n int, err error) {
	f, _ := os.Create("images/" + strconv.Itoa(int(time.Now().Unix())) + ".jpg")
	f.Write(p)
	data := base64.StdEncoding.EncodeToString(p)
	err = writer.conn.WriteMessage(websocket.TextMessage, []byte(data))
	return len(p), err
}

func writeErrors(errCh chan error) {
	for x := range errCh {
		fmt.Fprintf(os.Stderr, "%v\n", x)
	}
}
