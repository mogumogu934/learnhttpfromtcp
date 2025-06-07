package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/mogumogu934/learnhttpfromtcp/internal/request"
	"github.com/mogumogu934/learnhttpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode int
	Message    string
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(msgBytes))
	response.WriteHeaders(w, headers)
	w.Write(msgBytes)
}

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusCodeOK)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}
