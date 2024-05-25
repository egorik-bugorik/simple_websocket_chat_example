package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func newServer() *Server {
	return &Server{conns: map[*websocket.Conn]bool{}}
}

func (s *Server) handleWs(ws *websocket.Conn) {

	log.Println("Inner connection from client :::: ", ws.RemoteAddr())
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {

	var buf = make([]byte, 1024)

	for {

		read, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {

				break

			}

			log.Println("Error on read --->>>>  ", err)
			continue
		}
		msg := buf[:read]
		log.Println("Message from client ::: ", string(msg))

		ws.Write([]byte("Thank you fro message!!!"))

	}

}
func main() {

	server := newServer()

	http.Handle("/addr", websocket.Handler(server.handleWs))

	log.Fatal(http.ListenAndServe(":3000", nil))

}
