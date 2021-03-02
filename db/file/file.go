package file

import (
	"io/ioutil"
	"os"
)

var DefaultDB = DB{
	Dir: "data/",
}

type DB struct {
	Dir string
}

func (db *DB) SaveFile(path string, b []byte) error {
	err := os.MkdirAll(db.Dir, 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(db.Dir+path, b, 0755)
}
