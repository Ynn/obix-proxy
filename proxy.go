package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/xml"
	"time"
	"os"
	"crypto/tls"
)

type ObixValue struct {
	XMLName xml.Name
	Value string `xml:"val,attr"`
	Display string `xml:"display,attr"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("OBIX_SERVER_URL") + r.URL.Path[1:]
	fmt.Printf("retrieve %s \n", url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Host",r.Header.Get("Host"))
	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extractName :=r.URL.Query().Get("extract")
	format := r.URL.Query().Get("format")
	serverName:=os.Getenv("OBIX_NAME");

	if format == "" {
		format = r.Header.Get("Accept")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	v := ObixValue{}
	err = xml.Unmarshal(body, &v)

	if resp.StatusCode!=http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, "%s", body)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	if v.XMLName.Local == "err" {
		http.Error(w, v.Display, http.StatusInternalServerError)
		return
	}

	switch format {
	case "line", "application/line", "application/influxdb-line":
		w.Header().Set("Content-type", "text/plain")
		timestamp := time.Now().UnixNano()
		fmt.Fprintf(w, "%s %s=%s %d",serverName, extractName, v.Value, timestamp)
	case "text","text/plain":
		w.Header().Set("Content-type", "text/plain")
		fmt.Fprintf(w, "%s", v.Value)
	case "json", "application/json":
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, "{\"%s\" : %s}", extractName, v.Value)
	default:
		w.Header().Set("Content-type", "application/xml")
		fmt.Fprintf(w, "%s", body)
	}


}

func main() {
	fmt.Printf("Starting proxy for %s \n", os.Getenv("OBIX_SERVER_URL"))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
