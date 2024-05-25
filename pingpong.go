package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Context struct {
	pathReplacerRegexp *regexp.Regexp
	pingDomain         string
	pongDomain         string
	showLog            bool
}

type handler func(http.ResponseWriter, *http.Request)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	ctx := Context{}
	ctx.pingDomain = os.Getenv("PING")
	if len(ctx.pingDomain) == 0 {
		ctx.pingDomain = "http://0.0.0.0:" + port
	}
	ctx.pongDomain = os.Getenv("PONG")
	if len(ctx.pongDomain) == 0 {
		ctx.pongDomain = "http://0.0.0.0:" + port
	}
	envLog := os.Getenv("LOG")
	if len(envLog) == 0 || envLog != "false" {
		ctx.showLog = true
	}

	ctx.pathReplacerRegexp = regexp.MustCompile(`/(p[io]ng)/\d+`)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping/{counter}", Handler(ctx))
	mux.HandleFunc("/pong/{counter}", Handler(ctx))

	if ctx.showLog {
		fmt.Printf("Context: %+v\n", ctx)
	}
	fmt.Printf("listening on port %s ...\n", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
		panic(err)
	}
}

type GetResponse struct {
	Origin  string `json:"origin,omitempty"`
	Counter int    `json:"counter,omitempty"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Handler(ctx Context) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()

		if ctx.showLog {
			fmt.Printf("\n> GET %s", r.RequestURI)
		}

		counter, err := strconv.Atoi(r.PathValue("counter"))
		if err != nil {
			if ctx.showLog {
				fmt.Printf(" ->    WARN: counter must be integer: %s\n", err.Error())
			}
			http.Error(w, "counter must be integer", http.StatusBadRequest)
			return
		}

		payload := GetResponse{}
		payload.Origin = ctx.pathReplacerRegexp.ReplaceAllString(r.RequestURI, "$1")
		payload.Counter = counter

		payloads := []*GetResponse{&payload}

		if counter > 1 {
			var url string
			switch payload.Origin {
			case "ping":
				url = fmt.Sprintf("%s/pong/%d", ctx.pongDomain, (counter - 1))
			case "pong":
				url = fmt.Sprintf("%s/ping/%d", ctx.pingDomain, (counter - 1))
			}
			responses := []*GetResponse{}
			if ctx.showLog {
				fmt.Printf("\n< GET %s", url)
			}
			res, err := http.Get(url)
			if err != nil {
				p := &GetResponse{}
				p.Error = err.Error()
				responses = append(responses, p)
			} else {
				defer res.Body.Close()
				if err := json.NewDecoder(res.Body).Decode(&responses); err != nil {
					p := &GetResponse{}
					p.Error = err.Error()
					responses = append(responses, p)
				}
			}
			payloads = append(payloads, responses...)
		}

		encoder := json.NewEncoder(w)

		payload.Latency = time.Since(startedAt).String()
		if err := encoder.Encode(payloads); err != nil {
			fmt.Printf(" ->    ERROR: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
