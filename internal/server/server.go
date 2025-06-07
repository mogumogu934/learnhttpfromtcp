package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/mogumogu934/learnhttpfromtcp/internal/request"
	"github.com/mogumogu934/learnhttpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	serverRunning atomic.Bool
	handler       Handler
	listener      net.Listener
}

func Serve(port int, handler Handler) (*Server, error) {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("unable to create listener: %v", err)
	}

	s := &Server{
		handler:  handler,
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
			if s.serverRunning.Load() == false {
				log.Print("server is closed")
				return
			}
			log.Printf("unable to accept new connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		b := fmt.Sprintf("unable to parse request: %v", err)
		body := []byte(b)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
	return
}
