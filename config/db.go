package config

import (
	"github.com/gocroot/helper/atdb"
)

var MongoString string = "MONGOSTRING"

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "petapedia",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
