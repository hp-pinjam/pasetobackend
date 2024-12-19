package pasetobackend

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Private  string `json:"private,omitempty" bson:"private,omitempty"`
	Public   string `json:"public,omitempty" bson:"public,omitempty"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
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
	Admin string    `json:"admin"`
	Hp    string    `json:"hp"`
	Role  string    `json:"role"`
	Exp   time.Time `json:"exp"`
	Iat   time.Time `json:"iat"`
	Nbf   time.Time `json:"nbf"`
}

// type Crawling struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" `
// 	Created_at string             `json:"created_at" bson:"created_at"`
// 	Full_text  string             `json:"full_text" bson:"full_text"`
// 	Username   string             `json:"username" bson:"username"`
// 	Location   string             `json:"location" bson:"location"`
// }

type Hp struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" `
	Nomorid     int                `json:"nomorid" bson:"nomorid"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Status      bool               `json:"status" bson:"status"`
}

type About struct {
	ID          int    `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Image       string `json:"image" bson:"image"`
	Status      bool   `json:"status" bson:"status"`
}

// type Contact struct {
// 	ID       int    `json:"id" bson:"id"`
// 	FullName string `json:"fullname" bson:"fullname"`
// 	Email    string `json:"email" bson:"email"`
// 	Phone    string `json:"phone" bson:"phone"`
// 	Message  string `json:"image" bson:"image"`
// 	Status   bool   `json:"status" bson:"status"`
// }

type Workout struct {
	Name       string `json:"name" bson:"name"`
	Gif        string `json:"gif" bson:"gif"`
	Repetition string `json:"repetition" bson:"repetition"`
	Calories   int    `json:"calories" bson:"calories"`
}
