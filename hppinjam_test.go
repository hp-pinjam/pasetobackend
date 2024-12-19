package pasetobackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNewAdminRole(t *testing.T) {
	var admindata Admin
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "akuriziq"
	admindata.Role = "admin"
	mconn := SetConnection("MONGOSTRING", "Fitness")
	CreateNewAdminRole(mconn, "admin", admindata)
}

func TestDeleteAdmin(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "musa@gmail.com"
	DeleteAdmin(mconn, "admin", admindata)
}

func CreateNewAdminToken(t *testing.T) {
	var admindata Admin
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "akurijik"
	admindata.Role = "admin"

	// Create a MongoDB connection
	mconn := SetConnection("MONGOSTRING", "hppinjam")

	// Call the function to create a admin and generate a token
	err := CreateAdminAndAddToken("your_private_key_env", mconn, "admin", admindata)

	if err != nil {
		t.Errorf("Error creating admin and token: %v", err)
	}
}

func TestGFCPostHandlerAdmin(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "akurijik"
	admindata.Role = "admin"
	CreateNewAdminRole(mconn, "admin", admindata)
}

func TestHp(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var hpdata Hp
	hpdata.Nomorid = 1
	hpdata.Title = "iphone"
	hpdata.Description = "hp sultan"
	hpdata.Image = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"
	CreateNewHp(mconn, "hp", hpdata)
}

func TestAllHp(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	hp := GetAllHp(mconn, "hp")
	fmt.Println(hp)
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "ganteng"
	hash, _ := HashPass(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)
	match := CompareHashPass(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("alagaday", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "edi@gmail.com"
	admindata.Password = "pecin"

	filter := bson.M{"email": admindata.Email}
	res := atdb.GetOneDoc[Admin](mconn, "admin", filter)
	fmt.Println("Mongo Admin Result: ", res)
	hash, _ := HashPass(admindata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CompareHashPass(admindata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "bangsat"
	admindata.Password = "ganteng"

	anu := IsPasswordValid(mconn, "admin", admindata)
	fmt.Println(anu)
}

func TestAdminFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "hppinjam@gmail.com"
	admindata.Password = "hebat"
	admindata.Role = "admin"
	CreateAdmin(mconn, "admin", admindata)
}

func TestLoginn(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "hppinjam@gmail.com"
	admindata.Password = "hebat"
	IsPasswordValid(mconn, "admin", admindata)
	fmt.Println(admindata)
}

func TestInsertWorkout(t *testing.T) {
	// Mock koneksi database
	mconn := SetConnection("MONGOSTRING", "Fitness") // Pastikan string koneksi sesuai konfigurasi Anda

	// Data workout untuk di-insert
	workoutData := Workout{
		Name:       "Push Up",
		Gif:        "https://example.com/pushup.gif",
		Repetition: "15",
		Calories:   50,
	}

	// Panggil fungsi insertWorkout
	insertedID := insertWorkout(mconn, "workout", workoutData)

	// Validasi hasil insert
	if insertedID == primitive.NilObjectID {
		t.Errorf("Insert workout gagal: ID tidak ditemukan")
	} else {
		t.Logf("Insert workout sukses: ID %v, Data %+v", insertedID, workoutData)
	}
}
