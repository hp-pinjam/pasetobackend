package pasetobackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// user
func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	CreateNewUserRole(mconn, "user", userdata)
}

// user
func TestDeleteUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Email = "farhanriziq@gmail.com"
	DeleteUser(mconn, "user", userdata)
}

// user
func CreateNewUserToken(t *testing.T) {
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"

	// Create a MongoDB connection
	mconn := SetConnection("MONGOSTRING", "hppinjam")

	// Call the function to create a user and generate a token
	err := CreateUserAndAddToken("your_private_key_env", mconn, "user", userdata)

	if err != nil {
		t.Errorf("Error creating user and token: %v", err)
	}
}

// user
func TestGFCPostHandlerUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"
	CreateNewUserRole(mconn, "user", userdata)
}

// Test Insert Hp
func TestInsertHp(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var hpdata Hp
	hpdata.Nomorid = 1
	hpdata.Title = "garut"
	hpdata.Description = "waw garut keren banget"
	hpdata.Image = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"
	CreateNewHp(mconn, "hp", hpdata)
}

// Test All Hp
func TestAllHp(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	hp := GetAllHp(mconn, "hp")
	fmt.Println(hp)
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "riziq"
	hash, _ := HashPass(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)
	match := CompareHashPass(password, hash)
	fmt.Println("Match:   ", match)
}

// pasetokey
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("hppinjam", privateKey)
	fmt.Println(hasil, err)
}

// user
func TestHashFunctionUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[Admin](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPass(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CompareHashPass(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)
}

// user
func TestUserIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"

	anu := UserIsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

// user
func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"
	CreateUser(mconn, "user", userdata)
}

// user
func TestLoginUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "farhanriziq"
	userdata.Email = "farhanriziq@gmail.com"
	userdata.Password = "riziq"
	UserIsPasswordValid(mconn, "user", userdata)
	fmt.Println(userdata)
}
