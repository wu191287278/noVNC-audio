package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	staticDir := flag.String("static", "./static", "Html static directory")
	downloadsDir := flag.String("downloads", "./downloads", "Html static directory")
	address := flag.String("address", ":8888", "Server port")
	udpAddress := flag.String("udpAddress", ":1234", "Jsmpeg Udp Server port")
	vncAddress := flag.String("vncAddress", "localhost:5900", "Vnc Server port")
	flag.Parse()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Location", "/static/vnc")
		writer.WriteHeader(302)
	})
	http.Handle("/websockify", websocket.Handler(func(wsconn *websocket.Conn) {
		defer wsconn.Close()
		var d net.Dialer
		var address = *vncAddress
		conn, err := d.DialContext(wsconn.Request().Context(), "tcp", address)
		if err != nil {
			log.Printf("[%s] [VNC_ERROR] [%v]", address, err)
			return
		}
		defer conn.Close()
		wsconn.PayloadType = websocket.BinaryFrame
		go func() {
			io.Copy(wsconn, conn)
			wsconn.Close()
			log.Printf("[%s] [VNC_SESSION_CLOSED]", address)
		}()
		io.Copy(conn, wsconn)
		log.Printf("[%s] [VNC_CLIENT_DISCONNECTED]", address)
	}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))
	http.Handle("/downloads", http.StripPrefix("/downloads/", http.FileServer(http.Dir(*downloadsDir))))
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	var writers = new(WsMultiWriter)
	writers.writers = map[*websocket.Conn]chan *[]byte{}
	go func() {
		RunJsmpegUDP(*udpAddress, writers)
	}()
	http.Handle("/audio", websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()
		conn.PayloadType = websocket.BinaryFrame
		ch := make(chan *[]byte, 1)
		writers.writers[conn] = ch
		for {
			select {
			case <-ch:
			}
			break
		}
	}))
	log.Printf("Http listening os %s \n", *address)

	log.Fatal(http.ListenAndServe(*address, nil))

}

func RunJsmpegUDP(address string, writer io.Writer) {
	udpAddr, _ := net.ResolveUDPAddr("udp4", address)

	//监听端口
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer udpConn.Close()

	fmt.Printf("Jsmpeg udp listening on %s \n", address)
	io.Copy(writer, udpConn)

}

type WsMultiWriter struct {
	writers map[*websocket.Conn]chan *[]byte
}

func (t *WsMultiWriter) Write(p []byte) (n int, err error) {
	for k, v := range t.writers {
		if v == nil {
			continue
		}
		n, err = k.Write(p)
		if err != nil {
			v <- nil
			t.writers[k] = nil
			continue
		}
		if n != len(p) {
			return
		}
	}
	return len(p), nil
}
