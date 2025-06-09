package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mogumogu934/learnhttpfromtcp/internal/headers"
	"github.com/mogumogu934/learnhttpfromtcp/internal/request"
	"github.com/mogumogu934/learnhttpfromtcp/internal/response"
	"github.com/mogumogu934/learnhttpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		videoHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeOK)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCodeOK)
	h := response.GetDefaultHeaders(0)
	delete(h, "Content-Length")
	h.Overwrite("Transfer-Encoding", "chunked")
	h.Overwrite("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)
	fullBody := make([]byte, 0)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Println("unable to write chunked body:", err)
				break
			}
			fullBody = append(fullBody, buf[:n]...)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("unable to read response body:", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Println("unable to write chunked body done:", err)
	}

	t := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	t.Overwrite("X-Content-SHA256", sha256)
	t.Overwrite("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	err = w.WriteTrailers(t)
	if err != nil {
		log.Println("unable to write trailers:", err)
	}
}

func videoHandler(w *response.Writer, _ *request.Request) {
	h := response.GetDefaultHeaders(0)
	h.Overwrite("Content-Type", "video/mp4")
	f, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		log.Println("unable to read contents from video:", err)
		return
	}
	h.Overwrite("Content-Length", fmt.Sprintf("%d", len(f)))

	w.WriteStatusLine(response.StatusCodeOK)
	w.WriteHeaders(h)
	w.WriteBody(f)
	return
}
