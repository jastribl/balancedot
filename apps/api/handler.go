package api

import (
	"log"
	"net/http"
)

// Handler is a a generic handler for the api
type Handler func(ResponseWriter, *Request) WriterResponse

func (fn Handler) ServeHTTP(httpW http.ResponseWriter, httpR *http.Request) {
	w := NewResponseWriter(httpW)
	r := NewRequest(httpR)
	defer func() {
		if ree := recover(); ree != nil {
			log.Printf("recovering from fatal: %+v", ree)
			w.SendUnexpectedError(ree)
		}
	}()
	err := fn(w, r)
	if err != WriterResponseSuccess {
		w.SendUnexpectedError(err, "Handler response was something other than WriterResponseSuccess")
	}
}
