package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	core_config "github.com/AlvinTendio/minder/config"
)

type ctxKey struct{}

// UserName ...
var UserName string

// UserID ...
var UserID uint64

// Action ...
var Action string

// Endpoint ...
var Endpoint string

// IPAddress ...
var IPAddress string

// BodyRequest ...
var BodyRequest string

type route struct {
	method     string
	regex      *regexp.Regexp
	handler    http.HandlerFunc
	actionName string
}

var routes = []route{}

type RestHandlerService interface {
	Serve(w http.ResponseWriter, r *http.Request)
}
type restHandlerService struct {
	config core_config.Config
}

func NewRestHandlerService(c core_config.Config) RestHandlerService {
	return &restHandlerService{c}
}

func Route(method, pattern string, handler http.HandlerFunc, actionName string) {
	pattern = "/minder" + pattern
	routes = append(routes, route{method, regexp.MustCompile("^" + pattern + "$"), handler, actionName})
}

func Param(r *http.Request, index int) string {
	val := r.Context().Value(ctxKey{})
	if val == nil {
		return ""
	}
	params, b := val.([]string)
	if !b {
		log.Println(" Invalid ")
	}
	return params[index]
}

func (t *restHandlerService) Serve(w http.ResponseWriter, r *http.Request) {
	IPAddress = GetIPAddress(r)
	// Capture Request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error", err)
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	BodyRequest = string(body)

	var allow []string
	// auth := false
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			Action = route.actionName
			Endpoint = r.URL.Path

			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func GetIPAddress(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func ResponseWrite(r *http.Request, rw http.ResponseWriter, response interface{}, statusCode int) {
	// Capture return  responce
	res, err := json.Marshal(response)
	if err != nil {
		log.Println("Error", err)
	}
	responses := string(res)

	// Capture GET method parameters
	getReq, err := json.Marshal(r.URL.Query())
	if err != nil {
		log.Println("Error", err)
	}
	getRequest := string(getReq)
	log.Println(getRequest)
	log.Println(responses)
	if rw.Header().Get("Content-Disposition") != "" {
		rw.Header().Set("Content-Type", "application/octet-stream")
	} else {
		rw.Header().Set("Content-type", "application/json")
		rw.WriteHeader(statusCode)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			log.Println("ERROR")
		}
	}
}
