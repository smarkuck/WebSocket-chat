package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/smarkuck/WebSocket-chat/bot"
	"golang.org/x/net/websocket"
)

const websocketAddr = "localhost:80"
const tcpAddr = "localhost:4000"

func main() {
	go listenTCP()
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(websocketAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func listenTCP() {
	l, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(c)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, websocketAddr)
}

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go match(s)
	<-s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by other goroutine
	case p := <-partner:
		chat(p, c)
	case <-time.After(5 * time.Second):
		chat(bot.New(), c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(io.MultiWriter(w, bot.Chain), r)
	errc <- err
}
