package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Query bson.M

func NewDB(host string) (*DB, error) {
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	return &DB{session: session}, nil
}

type DB struct {
	name    string
	session *mgo.Session
}

func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) DB() *mgo.Database {
	return db.session.DB(db.name)
}

func (db *DB) C(name string) *mgo.Collection {
	return db.DB().C(name)
}

func (db *DB) Upsert(name string, q Query, v interface{}) (updated bool, err error) {
	info, err := db.C(name).Upsert(q, v)
	updated = info != nil && info.UpsertedId != nil
	return
}

func (db *DB) EnsureIndex() error {
	indexes := []mgo.Index{
		{
			Key:        []string{"lastbuildid"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		},
	}

	for _, index := range indexes {
		err := db.C("new_builds").EnsureIndex(index)
		if err != nil {
			return err
		}
	}

	return nil
}
