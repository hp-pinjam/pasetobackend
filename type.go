package pasetobackend

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username     string `json:"username" bson:"username"`
	NPM          string `json:"npm" bson:"npm"`
	Password     string `json:"password" bson:"password"`
	PasswordHash string `json:"passwordhash" bson:"passwordhash"`
	Email        string `bson:"email,omitempty" json:"email,omitempty"`
	Role         string `json:"role,omitempty" bson:"role,omitempty"`
	Token        string `json:"token,omitempty" bson:"token,omitempty"`
	Private      string `json:"private,omitempty" bson:"private,omitempty"`
	Public       string `json:"public,omitempty" bson:"public,omitempty"`
}

type Admin struct {
	Username     string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
	PasswordHash string `json:"passwordhash" bson:"passwordhash"`
	Email        string `bson:"email,omitempty" json:"email,omitempty"`
	Role         string `json:"role,omitempty" bson:"role,omitempty"`
	Token        string `json:"token,omitempty" bson:"token,omitempty"`
	Private      string `json:"private,omitempty" bson:"private,omitempty"`
	Public       string `json:"public,omitempty" bson:"public,omitempty"`
}

type Parkiran struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" `
	ParkiranId     int                `json:"parkiranid" bson:"parkiranid"`
	Nama           string             `json:"nama" bson:"nama"`
	NPM            string             `json:"npm" bson:"npm"`
	Jurusan        string             `json:"jurusan" bson:"jurusan"`
	NamaKendaraan  string             `json:"namakendaraan" bson:"namakendaraan"`
	NomorKendaraan string             `bson:"nomorkendaraan,omitempty" json:"nomorkendaraan,omitempty"`
	JenisKendaraan string             `json:"jeniskendaraan,omitempty" bson:"jeniskendaraan,omitempty"`
	Status         bool               `json:"status" bson:"status"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  bool        `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

type Payload struct {
	Id   primitive.ObjectID `json:"id"`
	Role string             `json:"role"`
	Exp  time.Time          `json:"exp"`
	Iat  time.Time          `json:"iat"`
	Nbf  time.Time          `json:"nbf"`
}
