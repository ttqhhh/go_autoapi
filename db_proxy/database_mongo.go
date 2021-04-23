package db_proxy

import (
	"github.com/beego/beego/v2/server/web"
	"gopkg.in/mgo.v2"
	"log"
)

var globalS *mgo.Session

func InitMongoDB() {
	host, _ := web.AppConfig.String("mongo_host")
	diaInfo := &mgo.DialInfo{
		Addrs: []string{host},
	}

	s, err := mgo.DialWithInfo(diaInfo)

	if err != nil {
		log.Fatalln("create session error", err)
	}

	globalS = s
}

func Connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB(db).C(collection)
	return s, c
}

func Insert(db, collection string, docs ...interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}

func Update(db, collection string, query, update interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()
	return c.Update(query, update)
}

func Remove(db, collection string, query interface{}) error {
	ms, c := Connect(db, collection)
	defer ms.Close()
	return c.Remove(query)
}
