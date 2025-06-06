package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/mogumogu934/learnhttpfromtcp/internal/response"
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
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Print(err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Print(err)
	}
	return
}
