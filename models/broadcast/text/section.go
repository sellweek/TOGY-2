package text

import (
	"appengine"
	"appengine/datastore"
	"github.com/sellweek/gaemodel"
	"reflect"
	"time"
)

const kind = "TextSection"

var typ = reflect.TypeOf(TextSection{})

type TextSection struct {
	Title             string
	ShowTitle         bool
	Start, End        time.Time
	Created, Modified time.Time
	Delays            uint
	Published         bool
	key               *datastore.Key `datastore:"-"`
}

func (ts *TextSection) Key() *datastore.Key {
	return ts.key
}

func (ts *TextSection) SetKey(k *datastore.Key) {
	ts.key = k
}

func (ts *TextSection) Kind() string {
	return kind
}

func (ts *TextSection) Ancestor() *datastore.Key {
	return nil
}

func (ts *TextSection) Save(c appengine.Context) (err error) {
	//Why do I have to put this into parentheses?
	//Go can't parse it without them
	if (ts.Created == time.Time{}) {
		ts.Created = time.Now()
	}

	ts.Modified = time.Now()

	err = gaemodel.Save(c, ts)
	return
}

func (ts *TextSection) Delete(c appengine.Context) (err error) {
	err = gaemodel.Delete(c, ts)
	return
}

func Get(c appengine.Context, key *datastore.Key) (*TextSection, error) {
	m, err := gaemodel.GetByKey(c, typ, key)
	if err != nil {
		return err
	}
	return m.(*TextSection)
}

func GetAll(c appengine.Context) ([]*TextSection, error) {
	ms, err := gaemodel.GetAll(c, typ, kind, 0, 0)
	if err != nil {
		return err
	}
	return ms.([]*TextSection)
}

func GetActive(c appengine.Context) ([]*TextSection, error) {
	now := time.Now()
	q := datastore.NewQuery(kind).Filter("Published =", true).
		Filter("Start <", now).Filter("End >", now)
	ms, err := gaemodel.MultiQuery(c, typ, kind, q)
	if err != nil {
		return err
	}

	return ms.([]*TextSection)
}
