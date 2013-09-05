package util

import (
	"appengine"
	"github.com/gorilla/mux"
	"github.com/mjibson/appstats"
	"net/http"
	"time"
)

//WARNING: This will not work in local SDK
//unless you add tzdata into google_appengine/goroot/lib/time
var Tz, _ = time.LoadLocation("Europe/Bratislava")

//Context is the type used for passing data to handlers
type Context struct {
	Ac   appengine.Context
	W    http.ResponseWriter
	R    *http.Request
	Vars map[string]string
}

//Handler maps standard net/http handlers to handlers accepting Context
func Handler(hand func(Context) error) http.Handler {
	return appstats.NewHandler(func(c appengine.Context, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := hand(Context{Ac: c, W: w, R: r, Vars: vars})
		if err != nil {
			c.Errorf("Error 500. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

//NormalizeDate strips the time part from time.Date leaving only
//year, month and day.
func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, Tz)
}
