package db

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

const (
	dbFileName = "shopwatcher.db"

	// ProductsBucket bucket name to store products
	ProductsBucket = "products"
)

var dbConnection *bolt.DB

func init() {
	var err error
	log := logger()

	dbConnection, err = bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		log.Fatalf("Error openning db %s", dbFileName)
	}

	dbConnection.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ProductsBucket))
		if err != nil {
			log.Fatalf("Error creating bucket %s", ProductsBucket)
		}
		return nil
	})
}

// Connection returns db connection
func Connection() *bolt.DB {
	return dbConnection
}

// Close closes db connection
func Close() {
	dbConnection.Close()
}

func logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"module": "db",
	})
}
