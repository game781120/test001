package utils

import (
	"net/http"
)

type StringRender struct {
	Data string
}

func (r StringRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	_, err := w.Write([]byte(r.Data))
	return err
}

func (r StringRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = []string{"text/event-stream; charset=utf-8"}

	if _, exist := header["Cache-Control"]; !exist {
		header["Cache-Control"] = []string{"no-cache"}
	}
}
