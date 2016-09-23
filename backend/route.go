package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Func    http.HandlerFunc
}

type apiFunc func(req *http.Request, params interface{}) (interface{}, error)

var AcornUrlPrefix string

var routes = []Route{
	Route{"index", "GET", "/", handleIndex},
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404 page not found"))
}

func wrapJson(inner apiFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var params interface{}
		var result interface{}
		var err error
		if req.ContentLength != 0 {
			if err = json.NewDecoder(req.Body).Decode(&params); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		result, err = inner(req, params)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		if err = json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
}

func wrapRequest(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		r.URL.Host = r.Host
		r.URL.Scheme = "http"
		AcornUrlPrefix = r.URL.Scheme + "://" + r.URL.Host
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8002")
		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = wrapRequest(http.HandlerFunc(handle404), "err-404")
	for _, route := range routes {
		handler := wrapRequest(route.Func, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.PathPrefix("/lbapi-mock/").Handler(wrapRequest(wrapJson(lbapiMockDashboard), "lbapi-mock"))

	return router
}
