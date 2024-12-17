package pasetobackend

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// admin
func CreateNewAdminRole(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.Password)
	if err != nil {
		return err
	}
	admindata.Password = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

func CreateAdminAndAddToken(privateKeyEnv string, mongoconn *mongo.Database, collection string, admindata Admin) error {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.Password)
	if err != nil {
		return err
	}
	admindata.Password = hashedPassword

	// Create a token for the admin
	tokenstring, err := watoken.Encode(admindata.Email, os.Getenv(privateKeyEnv))
	if err != nil {
		return err
	}

	admindata.Token = tokenstring

	// Insert the admin data into the MongoDB collection
	if err := atdb.InsertOneDoc(mongoconn, collection, admindata.Email); err != nil {
		return nil // Mengembalikan kesalahan yang dikembalikan oleh atdb.InsertOneDoc
	}

	// Return nil to indicate success
	return nil
}

func CreateResponse(status bool, message string, data interface{}) Response {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return response
}

func CreateAdmin(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.Password)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	adminid := admindata.Email
	tokenstring, err := watoken.Encode(adminid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get adminid
	adminidstring := watoken.DecodeGetId(publicKey, tokenstring)
	if adminidstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(adminidstring)
	admindata.Private = privateKey
	admindata.Public = publicKey
	admindata.Password = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

// hp
func CreateNewHp(mongoconn *mongo.Database, collection string, hpdata Hp) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, hpdata)
}

// hp function
func insertHp(mongoconn *mongo.Database, collection string, hpdata Hp) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, hpdata)
}

func DeleteHp(mongoconn *mongo.Database, collection string, hpdata Hp) interface{} {
	filter := bson.M{"nomorid": hpdata.Nomorid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedHp(mongoconn *mongo.Database, collection string, filter bson.M, hpdata Hp) interface{} {
	updatedFilter := bson.M{"nomorid": hpdata.Nomorid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, hpdata)
}

func GetAllHp(mongoconn *mongo.Database, collection string) []Hp {
	hp := atdb.GetAllDoc[[]Hp](mongoconn, collection)
	return hp
}
func GetAllHps(MongoConn *mongo.Database, colname string, email string) []Admin {
	data := atdb.GetAllDoc[[]Admin](MongoConn, colname)
	return data
}

func GetAllHpID(mongoconn *mongo.Database, collection string, hpdata Hp) Hp {
	filter := bson.M{
		"nomorid":     hpdata.Nomorid,
		"title":       hpdata.Title,
		"description": hpdata.Description,
		"image":       hpdata.Image,
	}
	hpID := atdb.GetOneDoc[Hp](mongoconn, collection, filter)
	return hpID
}

func insertWorkout(collection *mongo.Collection, workout Workout) {
	_, err := collection.InsertOne(context.TODO(), workout)
	if err != nil {
		log.Println("Error inserting workout:", err)
	}
}

// about function

// func InsertAbout(mongoconn *mongo.Database, collection string, aboutdata About) interface{} {
// 	return atdb.InsertOneDoc(mongoconn, collection, aboutdata)
// }

// func DeleteAbout(mongoconn *mongo.Database, collection string, aboutdata About) interface{} {
// 	filter := bson.M{"id": aboutdata.ID}
// 	return atdb.DeleteOneDoc(mongoconn, collection, filter)
// }

// func UpdatedAbout(mongoconn *mongo.Database, collection string, filter bson.M, aboutdata About) interface{} {
// 	updatedFilter := bson.M{"id": aboutdata.ID}
// 	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, aboutdata)
// }

// func GetAllAbout(mongoconn *mongo.Database, collection string) []About {
// 	about := atdb.GetAllDoc[[]About](mongoconn, collection)
// 	return about
// }

// contact function

// func InsertContact(mongoconn *mongo.Database, collection string, contactdata Contact) interface{} {
// 	return atdb.InsertOneDoc(mongoconn, collection, contactdata)
// }

// func GetAllContact(mongoconn *mongo.Database, collection string) []Contact {
// 	contact := atdb.GetAllDoc[[]Contact](mongoconn, collection)
// 	return contact
// }

// func GetIdContact(mongoconn *mongo.Database, collection string, contactdata Contact) Contact {
// 	filter := bson.M{"id": contactdata.ID}
// 	return atdb.GetOneDoc[Contact](mongoconn, collection, filter)
// }

//crawling function

// func GetAllCrawling(mongoconn *mongo.Database, collection string) []Crawling {
// 	crawling := atdb.GetAllDoc[[]Crawling](mongoconn, collection)
// 	return crawling
// }
