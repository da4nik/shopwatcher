package db

import (
	"encoding/json"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/da4nik/shopwatcher/types"
)

// LoadProduct - reads product from DB
func LoadProduct(url string) (types.Product, error) {
	url = strings.Trim(url, " \n")

	var product types.Product

	dbase := Connection()
	err := dbase.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(ProductsBucket))

		rawProduct := b.Get([]byte(url))
		if product, err = UnmarshalProduct(rawProduct); err != nil {
			return err
		}

		return nil
	})
	return product, err
}

// AllProducts returns all products
func AllProducts() (products []types.Product, err error) {
	dbase := Connection()
	err = dbase.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProductsBucket))
		c := b.Cursor()
		for k, data := c.First(); k != nil; k, data = c.Next() {
			product, err := UnmarshalProduct(data)
			if err == nil {
				products = append(products, product)
			}
		}
		return nil
	})
	return
}

// SaveProduct saves product
func SaveProduct(product types.Product) error {
	log := productsLogger()

	log.Debugf("Saving product %s\n", product.Name)

	dbase := Connection()
	err := dbase.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProductsBucket))

		data, errn := json.Marshal(product)
		if errn != nil {
			return errn
		}

		errn = b.Put([]byte(product.URL), data)
		if errn != nil {
			return errn
		}

		return nil
	})
	return err
}

// UnmarshalProduct - unmarshal product from db
func UnmarshalProduct(data []byte) (product types.Product, err error) {
	err = json.Unmarshal(data, &product)
	return
}

func productsLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"module": "db.product",
	})
}
