package lms

import (
	"errors"

	"github.com/gocroot/helper/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetRekapPendaftaranUsers(db *mongo.Database) (rkp RekapitulasiUser, err error) {
	//copy data user
	users, err := GetAllUser(db)
	if err != nil {
		err = errors.New("GetAllUser:" + err.Error())
		return
	}
	_, err = atdb.InsertManyDocs[User](db, "lmsusers", users)
	if err != nil {
		err = errors.New("InsertManyDocs:" + err.Error())
		return
	}
	//hitung setiap kelompok status user
	filter := bson.M{
		"profileapproved": 1,
		"roles":           "Sebagai Pengguna LMS",
	}
	count1, err := atdb.GetCountDoc(db, "lmsusers", filter)
	if err != nil {
		return
	}
	filter = bson.M{
		"profileapproved": 2,
		"roles":           "Sebagai Pengguna LMS",
	}
	count2, err := atdb.GetCountDoc(db, "lmsusers", filter)
	if err != nil {
		return
	}
	filter = bson.M{
		"profileapproved": 3,
		"roles":           "Sebagai Pengguna LMS",
	}
	count3, err := atdb.GetCountDoc(db, "lmsusers", filter)
	if err != nil {
		return
	}
	filter = bson.M{
		"profileapproved": 4,
		"roles":           "Sebagai Pengguna LMS",
	}
	count4, err := atdb.GetCountDoc(db, "lmsusers", filter)
	if err != nil {
		return
	}
	count5, err := atdb.GetCountDoc(db, "lmsusers", bson.M{"roles": "Sebagai Pengguna LMS"})
	if err != nil {
		return
	}
	// 1. Belum Lengkap
	// 2. Menunggu Persetujuan
	// 3. Disetujui
	// 4. Ditolak
	rkp = RekapitulasiUser{
		BelumLengkap:        count1,
		MenungguPersetujuan: count2,
		Disetujui:           count3,
		Ditolak:             count4,
		Total:               count5,
	}
	//drop collection user
	err = atdb.DropCollection(db, "lmsusers")
	if err != nil {
		return
	}
	return

}
