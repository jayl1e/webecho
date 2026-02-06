package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var listenAddr = ":8080"

func echoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(`{"severity": "INFO", "method":"%s", "proto": "%s", "msg": "access %s from %s"}`, r.Method, r.Proto, r.RequestURI, r.RemoteAddr)
	if r.ContentLength == 0 || r.Header.Get("Skip-Body") != "" || r.URL.Query().Get("skip_body") != "" {
		echoNoBodyHandler(w, r)
	} else {
		echoBodyHandler(w, r)
	}
}

func echoNoBodyHandler(w http.ResponseWriter, r *http.Request) {
	io.Copy(io.Discard, r.Body)
	defer r.Body.Close()
	pre_sleep(r)
	req_header, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Remote-Address", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	post_sleep(w, r)
	w.Write(req_header)
}

func pre_sleep(r *http.Request) {
	if pre_sleep, err := time.ParseDuration(r.URL.Query().Get("pre_sleep")); err == nil {
		log.Print("pre_sleep duration: ", pre_sleep)
		time.Sleep(pre_sleep)
	}
}
func post_sleep(w http.ResponseWriter, r *http.Request) {
	if post_sleep, err := time.ParseDuration(r.URL.Query().Get("post_sleep")); err == nil {
		log.Print("post_sleep duration: ", post_sleep)
		switch w := w.(type) {
		case http.Flusher:
			w.Flush()
		}
		time.Sleep(post_sleep)
	}
}

func echoBodyHandler(w http.ResponseWriter, r *http.Request) {
	pre_sleep(r)
	mw := multipart.NewWriter(w)
	w.Header().Set("Remote-Address", r.RemoteAddr)
	req_header, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Print(err)
	}
	header_part, err := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain"},
	})
	if err != nil {
		log.Print(err)
	}
	header_part.Write(req_header)
	post_sleep(w, r)
	body_mw, err := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{r.Header.Get("Content-Type")},
	})
	if err != nil {
		log.Print(err)
	}
	io.Copy(body_mw, r.Body)
	defer r.Body.Close()
	mw.Close()
}
func metricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	content := `up{svc="webecho"} 1`
	w.Write([]byte(content))
}

func init() {
	flag.StringVar(&listenAddr, "l", ":8080", "listen address")
}

func main() {
	flag.Parse()
	http.HandleFunc("/metrics", metricHandler)
	http.HandleFunc("/", echoHandler)

	log.SetFlags(0)
	h2h := h2c.NewHandler(http.DefaultServeMux, &http2.Server{})

	fmt.Println("started")
	log.Printf("Server starting on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, h2h); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
