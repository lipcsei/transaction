package db

import (
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var (
	once    sync.Once
	session *mgo.Session
	info    *mgo.DialInfo
)

// Mongo is an mgo adapter for the application
type Mongo struct{}

// GetCollection is return with an *mgo.Collection reference and an
// *mgo.Session reference (because can close the session after finish)
func (m *Mongo) GetCollection(collectionName string) (*mgo.Collection, *mgo.Session) {
	once.Do(func() {
		var err error
		if session, info, err = m.connect(); err != nil {
			panic(err)
		}
	})

	s := session.Copy()

	return s.DB(info.Database).C(collectionName), s
}

func (m *Mongo) connect() (s *mgo.Session, i *mgo.DialInfo, err error) {
	dialURI := "mongodb://127.0.0.1:27017/test"

	i, err = mgo.ParseURL(dialURI)
	s, err = mgo.Dial(dialURI)

	if err != nil {
		return
	}

	s.SetMode(mgo.Monotonic, true)
	s.SetSafe(&mgo.Safe{})

	return
}
