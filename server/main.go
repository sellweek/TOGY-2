package server

import (
	"github.com/gorilla/mux"
	"github.com/sellweek/TOGY-2/controllers"
	"github.com/sellweek/TOGY-2/util"
	"net/http"
)

func init() {
	r := mux.NewRouter()

	r.Handle("/", util.Handler(controllers.Home))

	http.Handle("/", r)
}
