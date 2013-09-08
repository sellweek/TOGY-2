package text

import (
	"appengine"
	"appengine/datastore"
	"github.com/sellweek/gaemodel"
	"reflect"
)

const slideKind = "TextSlide"

var slideTyp = reflect.TypeOf(TextSlide{})

type TextSlide struct {
	Title    string
	Contents []byte
	Section  *datastore.Key
	key      *datastore.Key `datastore:"-"`
}

func (ts *TextSlide) Key() *datastore.Key {
	return ts.key
}

func (ts *TextSlide) SetKey(k *datastore.Key) {
	ts.key = k
}

func (ts *TextSlide) Kind() string {
	return slideKind
}

func (ts *TextSlide) Ancestor() *datastore.Key {
	return ts.Section
}

func (ts *TextSlide) Save(c appengine.Context) error {
	return gaemodel.Save(c, ts)
}

func (ts *TextSlide) Delete(c appengine.Context) error {
	return gaemodel.Delete(c, ts)
}

func GetSlide(c appengine.Context, key *datastore.Key) (*TextSlide, error) {
	m, err := gaemodel.GetByKey(c, slideTyp, key)
	if err != nil {
		return nil, err
	}
	return m.(*TextSlide), nil
}

func GetSlidesBySection(c appengine.Context, sec *datastore.Key) ([]*TextSlide, error) {
	ms, err := gaemodel.GetByAncestor(c, slideTyp, slideKind, sec)
	if err != nil {
		return nil, err
	}
	return ms.([]*TextSlide), nil
}
