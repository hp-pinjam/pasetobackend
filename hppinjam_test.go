package pasetobackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// Test Password Hash
func TestGeneratePasswordHash(t *testing.T) {
	passwordhash := "riziqpass"
	hash, _ := HashPassword(passwordhash) // ignore error for the sake of simplicity

	fmt.Println("Password:", passwordhash)
	fmt.Println("Hash:    ", hash)
	match := CheckPasswordHash(passwordhash, hash)
	fmt.Println("Match:   ", match)
}

// Generate Private & Public Key
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("riziqpass", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "riziq"
	userdata.PasswordHash = "riziqpass"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.PasswordHash)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "riziq"
	userdata.PasswordHash = "riziqpass"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "riziq"
	// userdata.NPM = "1214000"
	userdata.Password = "riziqpass"
	userdata.PasswordHash = "riziqpass"
	// userdata.Email = "pakarbi2023@gmail.com"
	userdata.Role = "user"
	CreateUser(mconn, "user", userdata)
}

func TestAdminFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Username = "adminhp"
	admindata.Password = "adminhppass"
	admindata.PasswordHash = "adminhppass"
	// admindata.Email = "PakArbi2023@gmail.com"
	admindata.Role = "admin"
	CreateAdmin(mconn, "admin", admindata)
}

// func TestParkiran(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "PakarbiDB")
// 	var parkirandata Parkiran
// 	parkirandata.ParkiranId = "1"
// 	parkirandata.Nama = "Farhan Rizki Maulana"
// 	parkirandata.NPM = "1214020"
// 	parkirandata.Jurusan = "D4 Teknik Informatika"
// 	parkirandata.NamaKendaraan = "Supra X 125"
// 	parkirandata.NomorKendaraan = "F 1234 NR"
// 	parkirandata.JenisKendaraan = "Motor"
// 	CreateNewParkiran(mconn, "parkiran", parkirandata)
// }

// func TestAllParkiran(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "PakarbiDB")
// 	parkiran := GetAllParkiran(mconn, "parkiran")
// 	fmt.Println(parkiran)
// }
