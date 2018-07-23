package store

import (
	"github.com/globalsign/mgo"
	"github.com/vtfr/bossanova/common"
)

// mgoErr converts mgo errors to a better error type
func mgoErr(err error) error {
	switch {
	case err == nil:
		return nil
	case err == mgo.ErrNotFound:
		return common.ErrNotFound
	case mgo.IsDup(err):
		return common.ErrConflict
	default:
		//logrus.Panicln("Unknown error in MongoDB:", err)
		return err
	}
}
