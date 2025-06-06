package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	serverRunning atomic.Bool
	listener      net.Listener
}

func Serve(port int) (*Server, error) {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("unable to create listener: %v", err)
	}

	s := &Server{
		listener: lsn,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.serverRunning.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.serverRunning.Load() {
				log.Printf("unable to accept new connection: %v", err)
				continue
			} else {
				log.Printf("server is closed")
				return
			}
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello World!\n"

	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Printf("unable to write response: %v", err)
	}
	return
}
