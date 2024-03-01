package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ksrinimba/ssd-jwt-auth/ssdjwtauth"
)

var GATE_URL = "http://localhost:8001/internalToken"

func backEnd(w http.ResponseWriter, r *http.Request) {
	log.Println("Reached Backend")
	tokenStr := ssdjwtauth.GetTokenStrFromHeader(r)
	ssdTok, err := ssdjwtauth.DecodeToken(tokenStr)
	if err != nil {
		w.Write([]byte("Something really bad happend with the token, sorry"))
		return
	}
	if ssdTok.GetTokenType() == ssdjwtauth.SSDTokenTypeInternal {
		w.Write([]byte(fmt.Sprintf("Reached Authenticated Backend using a valid Internal Token:%+v", ssdTok)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Reached Authenticated Backend, some other token:%+v", ssdTok)))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called from:", r.Host)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Not found"))
}

func main() {
	//Get a Token from ssd-gate
	token, err := getInternalToken()
	if err != nil {
		log.Println("Unable to get token, exiting")
		os.Exit(1)
	}
	log.Println("Use This token for Auth testing:")
	log.Printf("ITOK=%s", token)
	log.Println("curl -vvv -H \"X-OpsMx-Auth: $ITOK\" http://localhost:8005/hello")
	// This part will change. Use this for now.
	ssdjwtauth.InitJWTSecret("myownsecret", []string{"x", "b"}, 3600, 3600*24*30, 300)

	r := mux.NewRouter()

	r.PathPrefix("/hello").HandlerFunc(backEnd)
	r.Use(ssdjwtauth.JWTAuthMiddleware)
	r.PathPrefix("/").HandlerFunc(NotFoundHandler)
	log.Println("Started server, listing on 8005")
	log.Fatal(http.ListenAndServe(":8005", r))
}

func getInternalToken() (string, error) {
	req, err := http.NewRequest("GET", GATE_URL, nil)
	if err != nil {
		log.Printf("A Request could not be created:%v", err)
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error calling:%v", err)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Not able to read")
		return "", err
	}
	return string(body), nil
}
