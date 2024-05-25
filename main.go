package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"time"
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

		s.broadcast(msg)

	}

}

func (s *Server) imitateBookOrder(ws *websocket.Conn) {
	log.Println("Client from ::: ", ws.RemoteAddr())

	for {
		payload := []byte(fmt.Sprintf("Order from client ::: %d", time.Now().UnixNano()))

		ws.Write(payload)
		time.Sleep(time.Second * 3)

	}
}

func (s *Server) broadcast(data []byte) {

	for ws := range s.conns {
		go func(ws *websocket.Conn) {

			_, err := ws.Write(data)
			if err != nil {
				log.Println("Error while ")
			}

		}(ws)
	}

}
func main() {

	server := newServer()

	http.Handle("/addr", websocket.Handler(server.handleWs))
	http.Handle("/book", websocket.Handler(server.imitateBookOrder))

	log.Fatal(http.ListenAndServe(":3000", nil))

}
