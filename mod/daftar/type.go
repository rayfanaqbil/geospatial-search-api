package daftar

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lokasi struct { //lokasi yang bisa melakukan presensi
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Nama     string             `bson:"nama,omitempty"`
	Batas    Geometry           `bson:"batas,omitempty"`
	Kategori string             `bson:"kategori,omitempty"`
}

type Geometry struct { //data geometry untuk lokasi presensi
	Type        string      `json:"type" bson:"type"`
	Coordinates interface{} `json:"coordinates" bson:"coordinates"`
}

type PresensiLokasi struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PhoneNumber string             `bson:"phonenumber,omitempty"`
	Lokasi      Lokasi             `bson:"lokasi,omitempty"`
	Selfie      bool               `bson:"selfie,omitempty"`
	IsMasuk     bool               `bson:"ismasuk,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

type PresensiSelfie struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IDUser      string             `json:"iduser,omitempty" bson:"iduser,omitempty"`
	CekInLokasi PresensiLokasi     `json:"cekinlokasi,omitempty" bson:"cekinlokasi,omitempty"`
	IsMasuk     bool               `json:"ismasuk,omitempty" bson:"ismasuk,omitempty"`
	Commit      string             `json:"commit,omitempty" bson:"commit,omitempty"`
	Remaining   int                `json:"remaining,omitempty" bson:"remaining,omitempty"`
	Filehash    string             `json:"filehash,omitempty" bson:"filehash,omitempty"`
}

type FaceDetect struct {
	IDUser    string `json:"iduser,omitempty" bson:"iduser,omitempty"`
	IDFile    string `json:"idfile,omitempty" bson:"idfile,omitempty"`
	Nfaces    int    `json:"nfaces,omitempty" bson:"nfaces,omitempty"`
	Base64Str string `json:"base64str,omitempty" bson:"base64str,omitempty"`
}

type KTPProps struct {
	IDUser    string `json:"iduser,omitempty" bson:"iduser,omitempty"`
	IDFile    string `json:"idfile,omitempty" bson:"idfile,omitempty"`
	Ncard     int    `json:"ncard,omitempty" bson:"ncard,omitempty"`
	Base64Str string `json:"base64str,omitempty" bson:"base64str,omitempty"`
}

type Config struct {
	PhoneNumber          string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	LeaflyURL            string `json:"leaflyurl,omitempty" bson:"leaflyurl,omitempty"`
	LeaflyURLKTP         string `json:"leaflyurlktp,omitempty" bson:"leaflyurlktp,omitempty"`
	LeaflySecret         string `json:"leaflysecret,omitempty" bson:"leaflysecret,omitempty"`
	DomyikadoPresensiURL string `json:"domyikadopresensiurl,omitempty" bson:"domyikadopresensiurl,omitempty"`
	DomyikadoUserURL     string `json:"domyikadouserurl,omitempty" bson:"domyikadouserurl,omitempty"`
	DomyikadoSecret      string `json:"domyikadosecret,omitempty" bson:"domyikadosecret,omitempty"`
}

type FaceInfo struct {
	PhoneNumber string `phonenumber:"secret,omitempty" bson:"phonenumber,omitempty"`
	Commit      string `json:"commit,omitempty" bson:"commit,omitempty"`
	Remaining   int    `json:"remaining,omitempty" bson:"remaining,omitempty"`
	FileHash    string `json:"filehash,omitempty" bson:"filehash,omitempty"`
	Error       string `json:"error,omitempty" bson:"error,omitempty"`
}

type Userdomyikado struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string             `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber          string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email                string             `bson:"email,omitempty" json:"email,omitempty"`
	GithubUsername       string             `bson:"githubusername,omitempty" json:"githubusername,omitempty"`
	GitlabUsername       string             `bson:"gitlabusername,omitempty" json:"gitlabusername,omitempty"`
	GitHostUsername      string             `bson:"githostusername,omitempty" json:"githostusername,omitempty"`
	Poin                 float64            `bson:"poin,omitempty" json:"poin,omitempty"`
	GoogleProfilePicture string             `bson:"googleprofilepicture,omitempty" json:"picture,omitempty"`
}
