package pasetobackend

import (
	"context"
	"fmt"
	"os"
	"time"

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

// func insertWorkout(collection *mongo.Collection, workout Workout) {
// 	_, err := collection.InsertOne(context.TODO(), workout)
// 	if err != nil {
// 		log.Println("Error inserting workout:", err)
// 	}
// }

func insertWorkout(mongoconn *mongo.Database, collection string, workout Workout) interface{} {
	return InsertOneDoc(mongoconn, collection, workout)
}

func GetAllWorkout(conn *mongo.Database, colname string) ([]Workout, error) {
	collection := conn.Collection(colname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query untuk mendapatkan semua dokumen
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Menyimpan hasil ke slice
	var workouts []Workout
	for cursor.Next(ctx) {
		var workout Workout

		// Decode dokumen dengan penanganan khusus untuk tipe campuran
		var rawWorkout map[string]interface{}
		if err := cursor.Decode(&rawWorkout); err != nil {
			return nil, err
		}

		// Tangani tipe campuran untuk repetition
		if repetition, ok := rawWorkout["repetition"]; ok {
			switch v := repetition.(type) {
			case string:
				rawWorkout["repetition"] = v
			case int, int32, int64:
				rawWorkout["repetition"] = fmt.Sprintf("%d", v) // Konversi integer ke string
			default:
				rawWorkout["repetition"] = ""
			}
		}

		// Decode ulang ke struct Workout
		data, err := bson.Marshal(rawWorkout)
		if err != nil {
			return nil, err
		}

		if err := bson.Unmarshal(data, &workout); err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	// Jika terjadi error selama iterasi cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func UpdatedWorkout(conn *mongo.Database, colname string, filter bson.M, updateData Workout) {
	collection := conn.Collection(colname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateData})
	if err != nil {
		fmt.Println("Error updating workout:", err)
	}
}

func DeleteWorkout(conn *mongo.Database, colname string, workoutData Workout) {
	collection := conn.Collection(colname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"number_id": workoutData.NumberID})
	if err != nil {
		fmt.Println("Error deleting workout:", err)
	}
}

func GetWorkoutByNumberID(conn *mongo.Database, colname string, numberID int) (*Workout, error) {
	collection := conn.Collection(colname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Filter untuk mencari dokumen berdasarkan NumberID
	filter := bson.M{"number_id": numberID}

	var workout Workout

	// Mendapatkan data berdasarkan filter
	err := collection.FindOne(ctx, filter).Decode(&workout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Tidak ditemukan
		}
		return nil, err // Error lainnya
	}

	return &workout, nil
}

func GenerateNumberID(conn *mongo.Database, colname string) int {
	collection := conn.Collection(colname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error counting documents:", err)
		return 1 // Default ke 1 jika terjadi error
	}
	return int(count + 1)
}
