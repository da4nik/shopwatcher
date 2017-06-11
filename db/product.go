package db

import (
	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/da4nik/shopwatcher/types"
)

// LoadProduct - reads product from DB
func LoadProduct(id int) (types.Product, error) {
	var product types.Product

	dbase := Connection()
	err := dbase.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(ProductsBucket))

		dbID := strconv.Itoa(id)
		rawProduct := b.Get([]byte(dbID))
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

		// Fill ID with next sequence value
		if product.ID == 0 {
			id, _ := b.NextSequence()
			product.ID = int(id)
		}

		data, errn := json.Marshal(product)
		if errn != nil {
			return errn
		}

		key := strconv.Itoa(product.ID)
		errn = b.Put([]byte(key), data)
		if errn != nil {
			return errn
		}

		return nil
	})
	return err
}

// DeleteProduct deletes product from list
func DeleteProduct(product types.Product) error {
	log := productsLogger()

	log.Debugf("Deleting product %s\n", product.Name)

	dbase := Connection()
	err := dbase.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProductsBucket))

		key := strconv.Itoa(product.ID)
		errn := b.Delete([]byte(key))
		if errn != nil {
			return errn
		}

		return nil
	})
	return err
}

// DeleteProductByID deletes product by it's ID
func DeleteProductByID(ID string) error {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	product, err := LoadProduct(intID)
	if err != nil {
		return err
	}

	return DeleteProduct(product)
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
