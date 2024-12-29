package pasetobackend

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func InsertAdmindata(MongoConn *mongo.Database, email, role, password string) (InsertedID interface{}) {
	req := new(Admin)
	req.Email = email
	req.Password = password
	req.Role = role
	return InsertOneDoc(MongoConn, "admin", req)
}

func InsertUserData(conn *mongo.Client, username, email, password, name string, height, weight float64, age int) error {
	// Tentukan koleksi
	collection := conn.Database("yourdbname").Collection("user")

	// Masukkan data ke koleksi MongoDB
	_, err := collection.InsertOne(context.Background(), bson.M{
		"username": username,
		"email":    email,
		"password": password, // Pastikan password sudah di-hash sebelum dipanggil
		"name":     name,
		"height":   height,
		"weight":   weight,
		"age":      age,
	})

	// Periksa apakah ada error saat insert
	if err != nil {
		return fmt.Errorf("failed to insert user data: %v", err)
	}

	return nil // Tidak ada error
}

func DeleteAdmin(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	filter := bson.M{"email": admindata.Email}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func FindAdmin(mongoconn *mongo.Database, collection string, admindata Admin) Admin {
	filter := bson.M{"email": admindata.Email}
	return atdb.GetOneDoc[Admin](mongoconn, collection, filter)
}

func IsExist(Tokenstr, PublicKey string) bool {
	id := watoken.DecodeGetId(PublicKey, Tokenstr)
	return id != ""
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, admindata Admin) bool {
	filter := bson.M{"email": admindata.Email}
	res := atdb.GetOneDoc[Admin](mongoconn, collection, filter)
	return CompareHashPass(admindata.Password, res.Password)
}

func MongoCreateConnection(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	// Insert dokumen ke koleksi MongoDB
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc Error: Failed to insert document into collection '%s': %v\n", collection, err)
		return nil // Mengembalikan nil jika ada error
	}

	// Logging sukses
	fmt.Printf("InsertOneDoc Success: Inserted document with ID %v into collection '%s'\n", insertResult.InsertedID, collection)
	return insertResult.InsertedID
}

func GetOneAdmin(MongoConn *mongo.Database, colname string, admindata Admin) Admin {
	filter := bson.M{"email": admindata.Email}
	data := atdb.GetOneDoc[Admin](MongoConn, colname, filter)
	return data
}
