package pasetobackend

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb
func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func GetConnectionMongo(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func CreateUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.PasswordHash)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	userid := userdata.Username
	tokenstring, err := watoken.Encode(userid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get userid
	useridstring := watoken.DecodeGetId(publicKey, tokenstring)
	if useridstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(useridstring)
	userdata.Private = privateKey
	userdata.Public = publicKey
	userdata.PasswordHash = hashedPassword

	// Insert the user data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

func InsertUserdata(MongoConn *mongo.Database, username, password, passwordhash, role string) (InsertedID interface{}) {
	req := new(User)
	req.Username = username
	// req.NPM = npm
	req.Password = password
	req.PasswordHash = passwordhash
	// req.Email = email
	req.Role = role
	return InsertSatuDoc(MongoConn, "user", req)
}

func CreateAdmin(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(admindata.PasswordHash)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	adminid := admindata.Username
	tokenstring, err := watoken.Encode(adminid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get userid
	adminidstring := watoken.DecodeGetId(publicKey, tokenstring)
	if adminidstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(adminidstring)
	admindata.Private = privateKey
	admindata.Public = publicKey
	admindata.PasswordHash = hashedPassword

	// Insert the user data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

// Cek Password Username
func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{
		"$or": []bson.M{
			{"username": userdata.Username},
			// {"email": userdata.Email},
		},
	}

	var res User
	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)

	if err == nil {
		// Mengasumsikan res.PasswordHash adalah password terenkripsi yang tersimpan di database
		return CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
	}
	return false
}

// Cek Password Email
// func IsPasswordValidEmail(mongoconn *mongo.Database, collection string, userdata User) bool {
// 	filter := bson.M{
// 		"$or": []bson.M{
// 			{"email": userdata.Email},
// 			{"npm": userdata.NPM},
// 		},
// 	}

// 	var res User
// 	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)

// 	if err == nil {
// Mengasumsikan res.PasswordHash adalah password terenkripsi yang tersimpan di database
// 		return CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
// 	}
// 	return false
// }

// Cek Password Admin
func IsPasswordValidAdmin(mongoconn *mongo.Database, collection string, userdata Admin) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[Admin](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

// FUNCTION CRUD
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

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
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

func InsertSatuDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

// Parkiran
// func CreateNewParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
// 	return atdb.InsertOneDoc(mongoconn, collection, parkirandata)
// }

// Function Parkiran
// func CreateParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
// 	return atdb.InsertOneDoc(mongoconn, collection, parkirandata)
// }

// func DeleteParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
// 	filter := bson.M{"parkiranid": parkirandata.ParkiranId}
// 	return atdb.DeleteOneDoc(mongoconn, collection, filter)
// }

// func UpdatedParkiran(mongoconn *mongo.Database, collection string, filter bson.M, parkirandata Parkiran) interface{} {
// 	filter = bson.M{"parkiranid": parkirandata.ParkiranId}
// 	return atdb.ReplaceOneDoc(mongoconn, collection, filter, parkirandata)
// }

// func GetAllParkiran(mongoconn *mongo.Database, collection string) []Parkiran {
// 	parkiran := atdb.GetAllDoc[[]Parkiran](mongoconn, collection)
// 	return parkiran
// }

// func GetAllParkiranID(mongoconn *mongo.Database, collection string, parkirandata Parkiran) Parkiran {
// 	filter := bson.M{
// 		"parkiranid":     parkirandata.ParkiranId,
// 		"nama":           parkirandata.Nama,
// 		"npm":            parkirandata.NPM,
// 		"jurusan":        parkirandata.Jurusan,
// 		"namakendaraan":  parkirandata.NamaKendaraan,
// 		"nomorkendaraan": parkirandata.NomorKendaraan,
// 		"jeniskendaraan": parkirandata.JenisKendaraan,
// 	}
// 	parkiranID := atdb.GetOneDoc[Parkiran](mongoconn, collection, filter)
// 	return parkiranID
// }
