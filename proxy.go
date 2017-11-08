package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/xml"
	"strconv"
	"time"
	"os"
)

type ObixValue struct {
	XMLName xml.Name
	Value string `xml:"val,attr"`
	Display string `xml:"display,attr"`
}


func handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("OBIX_SERVER_URL") + r.URL.Path[1:]
	resp, err := http.Get(url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if(resp.StatusCode!=http.StatusOK){
		w.WriteHeader(resp.StatusCode)
		return
	}

	extractName :=r.URL.Query().Get("extract")
	format := r.URL.Query().Get("format")
	serverName:=os.Getenv("OBIX_NAME");

	if(format == ""){
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

	if err != nil {
		return
	}

	if(v.XMLName.Local == "err"){
		http.Error(w, v.Display, http.StatusInternalServerError)
		return
	}

	if extractName != "" || format == "text" {
		switch format {
		case "line", "application/line", "application/influxdb-line":
			w.Header().Set("Content-type", "text/plain")
			timestamp := time.Now().UnixNano()

			switch v.XMLName.Local {
			case "int":
				if value, err := strconv.Atoi(v.Value); err == nil {
					fmt.Fprintf(w, "%s %s=%d %d",serverName, extractName, value, timestamp)
				}
			case "real":
				if value, err := strconv.ParseFloat(v.Value, 64); err == nil {
					fmt.Fprintf(w, "%s %s=%f %d",serverName, extractName, value, timestamp)
				}
			default:
				fmt.Fprintf(w, "%s %s=%s %d",serverName, extractName, v.Value, timestamp)
			}
		case "text","text/plain":
			w.Header().Set("Content-type", "text/plain")
			fmt.Fprintf(w, "%s", v.Value)
		case "json", "application/json":
			w.Header().Set("Content-type", "application/json")
			switch v.XMLName.Local {
			case "int":
				if value, err := strconv.Atoi(v.Value); err == nil {
					fmt.Fprintf(w, "{\"%s\" : %d}", extractName, value)
				}
			case "real":
				if value, err := strconv.ParseFloat(v.Value, 64); err == nil {
					fmt.Fprintf(w, "{\"%s\" : %f}", extractName, value)
				}
			default:
				fmt.Fprintf(w, "{\"%s\" : \"%s\"}", extractName, v.Value)
			}
		default:
			w.Header().Set("Content-type", "application/xml")
			fmt.Fprintf(w, "%s", body)
		}
	}else{
		w.Header().Set("Content-type", "application/xml")
		fmt.Fprintf(w, "%s", body)
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
