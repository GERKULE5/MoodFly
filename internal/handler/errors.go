package handler

import (
	apperror "MoodFly/pkg/error"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	code, msg := apperror.ToHTTP(err)
	http.Error(w, msg, code)
}
