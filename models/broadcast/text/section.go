package text

import (
	"appengine"
	"appengine/datastore"
	"github.com/sellweek/gaemodel"
	"reflect"
	"time"
)

const sectionKind = "TextSection"

var sectionTyp = reflect.TypeOf(TextSection{})

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
	return sectionKind
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
	slides, err := GetSlidesBySection(c, ts.Key())
	if err != nil {
		return
	}

	keys := make([]*datastore.Key, len(slides), len(slides))
	for i, s := range slides {
		keys[i] = s.Key()
	}

	if err = datastore.DeleteMulti(c, keys); err != nil {
		return
	}

	err = gaemodel.Delete(c, ts)
	return
}

func GetSection(c appengine.Context, key *datastore.Key) (*TextSection, error) {
	m, err := gaemodel.GetByKey(c, sectionTyp, key)
	if err != nil {
		return nil, err
	}
	return m.(*TextSection), nil
}

func GetAllSections(c appengine.Context) ([]*TextSection, error) {
	ms, err := gaemodel.GetAll(c, sectionTyp, sectionKind, 0, 0)
	if err != nil {
		return nil, err
	}
	return ms.([]*TextSection), nil
}

func GetActiveSections(c appengine.Context) ([]*TextSection, error) {
	now := time.Now()
	q := datastore.NewQuery(sectionKind).Filter("Published =", true).
		Filter("Start <", now).Filter("End >", now)
	ms, err := gaemodel.MultiQuery(c, sectionTyp, sectionKind, q)
	if err != nil {
		return nil, err
	}

	return ms.([]*TextSection), nil
}
