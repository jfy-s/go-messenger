package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var (
	authHttpBackendURL = "http://127.0.0.1:52521"
	chatWSBackendURL   = "ws://127.0.0.1:52522/ws"
)

func newHTTPReverseProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	rp := httputil.NewSingleHostReverseProxy(url)
	origDirector := rp.Director
	rp.Director = func(r *http.Request) {
		origDirector(r)
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			r.Header.Set("X-Real-IP", clientIP)
		}
		r.Header.Set("X-Forwarded-Proto", "http")
	}

	rp.Transport = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		DialContext:     (&net.Dialer{Timeout: 30 * time.Second}).DialContext,
		IdleConnTimeout: 30 * time.Second,
	}

	return rp
}

func HandleWebSocketProxy(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	backendUrl, _ := url.Parse(chatWSBackendURL)
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 30 * time.Second,
	}
	requestHeader := http.Header{}
	if token := r.Header.Get("Authorization"); token != "" {
		requestHeader.Set("Authorization", token)
	}

	backendConn, _, err := dialer.Dial(backendUrl.String(), requestHeader)
	if err != nil {
		clientConn.Close()
	}

	proxy := func(src, dst *websocket.Conn) {
		defer src.Close()
		defer dst.Close()
		for {
			mt, message, err := src.ReadMessage()
			if err != nil {
				return
			}
			err = dst.WriteMessage(mt, message)
			if err != nil {
				return
			}
		}
	}

	go proxy(clientConn, backendConn)
	go proxy(backendConn, clientConn)
}

func main() {
	httpProxy := newHTTPReverseProxy(authHttpBackendURL)
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", HandleWebSocketProxy)
	mux.Handle("/", httpProxy)
	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server.ListenAndServe()
}
