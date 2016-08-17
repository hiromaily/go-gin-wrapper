package mongodb

import (
	"errors"
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"time"
)

//MongoDB Ver.3.x
//https://gist.github.com/border/3489566

//Query and Projection Operators
//https://docs.mongodb.com/manual/reference/operator/query/

type MongoInfo struct {
	Session *mgo.Session
	Db      *mgo.Database
	C       *mgo.Collection
}

var mgInfo MongoInfo

//-----------------------------------------------------------------------------
// Settings
//-----------------------------------------------------------------------------
// create session object
func New(host, db, user, pass string, port uint16) {
	var err error
	if mgInfo.Session == nil {
		//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
		//mgInfo.session, _ = mgo.Dial("mongodb://user:pass@localhost:port/test")
		mongoUrl := ""
		if db == "" {
			//session, err := mgo.Dial("localhost:40001")
			mongoUrl = fmt.Sprintf("mongodb://%s:%d", host, port)
		} else {
			if user != "" && pass != "" {
				mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, pass, host, port, db)
			} else {
				mongoUrl = fmt.Sprintf("mongodb://%s:%d/%s", host, port, db)
			}
		}
		fmt.Printf("mongo url: %s\n", mongoUrl)
		mgInfo.Session, err = mgo.Dial(mongoUrl)
		//fmt.Println(mgInfo.Session)

		if err != nil {
			panic(err)
		}
		//mgInfo.Session.SetMode(mgo.Monotonic, true)
	}
}

func NewAdvanced(host, username, password, database string) {
	var err error
	if mgInfo.Session == nil {
		mgInfo.Session, err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:    []string{host},
			Username: username,
			Password: password,
			Database: database,
		})
		if err != nil {
			panic(err)
		}
	}
}

// singleton architecture
func GetMongo() *MongoInfo {
	if mgInfo.Session == nil {
		panic("Before call this, call New in addtion to arguments")
	}
	return &mgInfo
}

// close
func (mi *MongoInfo) Close() {
	mi.Session.Close()
}

//-----------------------------------------------------------------------------
// Database
//-----------------------------------------------------------------------------
// reset session.DB object to mi.db
func (mi *MongoInfo) GetDB(dbName string) *mgo.Database {
	//mi.db = mi.session.DB("test")
	mi.Db = mi.Session.DB(dbName)
	return mi.Db
}

func (mi *MongoInfo) DropDB(dbName string) error {
	err := mi.Session.DB(dbName).DropDatabase()
	return err
}

//-----------------------------------------------------------------------------
// Collection
//-----------------------------------------------------------------------------
// Set index
func (mi *MongoInfo) SetExpireOnCollection(sessionExpire time.Duration) error {
	sessionTTL := mgo.Index{
		Key:         []string{"createdAt"},
		Unique:      false,
		DropDups:    false,
		Background:  true,
		ExpireAfter: sessionExpire,
	} // sessionExpire is a time.Duration
	fmt.Println(sessionExpire)

	err := mi.C.EnsureIndex(sessionTTL)

	return err
}

// create collection
func (mi *MongoInfo) CreateCol(colName string) error {
	err := mi.Session.Run(bson.D{{"create", colName}}, nil)
	if err == nil {
		mi.C = mi.Db.C(colName)
	}
	return err
}

// get and set collection
func (mi *MongoInfo) GetCol(colName string) *mgo.Collection {
	mi.C = mi.Db.C(colName)
	return mi.C
}

// drop collection
func (mi *MongoInfo) DropCol(colName string) (err error) {
	err = mi.Db.C(colName).DropCollection()
	mi.C = nil
	return
}

//-----------------------------------------------------------------------------
// Document
//-----------------------------------------------------------------------------
//Get Count
func (mi *MongoInfo) GetCount() int {
	cnt, _ := mi.C.Count()
	return cnt
}

//Query One
func (mi *MongoInfo) FindOne(bd bson.M, data interface{}) error {

	//p := new(Person) //return is address of Person??
	//if name == "" {
	//	mi.C.Find(bson.M{}).One(data)
	//} else {
	//	mi.C.Find(bson.M{"name": name}).One(data)
	//}
	return mi.C.Find(bd).One(data)
}

// delete all documents record from collection. Version3.x
func (mi *MongoInfo) DelAllDocs(colName string) (err error) {
	if colName != "" {
		//mi.Db.C(colName).Remove(bson.M{})
		_, err = mi.Db.C(colName).RemoveAll(bson.M{})
	} else {
		//mi.C.Remove(bson.M{})
		_, err = mi.C.RemoveAll(bson.M{})
	}
	return
}

//-----------------------------------------------------------------------------
// Util
//-----------------------------------------------------------------------------
// convert datetime GMT
// MongoDB stores times in UTC by default
func ConvertDateTime() {
	//user.CreatedAt.Local()
}

//-----------------------------------------------------------------------------
// Load Json
//-----------------------------------------------------------------------------
func LoadJsonFile(filePath string) ([]byte, error) {
	// Loading jsonfile
	if filePath == "" {
		err := errors.New("Nothing Json File")
		return nil, err
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
