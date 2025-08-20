package server

import (
	"encoding/json"
	"net/http"
)

// Реализуй меня.
type jsonWriter struct {
	log ILogger
	w   http.ResponseWriter
}

func (jw jsonWriter) Write(v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		jw.log.Error(err.Error())

		jw.w.WriteHeader(http.StatusInternalServerError)
		jw.w.Write([]byte(err.Error()))

		return
	}

	jw.w.Header().Set("content-type", "application/json")
	jw.w.Write(bytes)
}
