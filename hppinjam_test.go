package pasetobackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.Email = "farhan@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	CreateNewUserRole(mconn, "user", userdata)
}

func TestCreateNewAdminRole(t *testing.T) {
	var admindata Admin

	admindata.Email = "farhan@gmail.com"
	admindata.Password = "riziq"
	admindata.Role = "admin"
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	CreateNewAdminRole(mconn, "admin", admindata)
}

func CreateNewUserToken(t *testing.T) {
	var userdata User
	userdata.Email = "farhan@gmail.com"
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

func CreateNewAdminToken(t *testing.T) {
	var admindata User
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "riziq"
	admindata.Role = "admin"

	// Create a MongoDB connection
	mconn := SetConnection("MONGOSTRING", "hppinjam")

	// Call the function to create a user and generate a token
	err := CreateUserAndAddToken("your_private_key_env", mconn, "admin", admindata)

	if err != nil {
		t.Errorf("Error creating admin and token: %v", err)
	}
}

func TestGFCPostHandlerUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Email = "farhan@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"
	CreateNewUserRole(mconn, "user", userdata)
}

// Test Password Hash
func TestGeneratePasswordHash(t *testing.T) {
	passwordhash := "riziq"
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
	hasil, err := watoken.Encode("riziq", privateKey)
	fmt.Println(hasil, err)
}

func TestGenerateAdminPasswordHash(t *testing.T) {
	password := "riziq"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)
	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Email = "farhan@gmail.com"
	userdata.Password = "riziq"

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
	userdata.Email = "farhan@gmail.com"
	userdata.Password = "riziq"

	anu := IsPasswordValid(mconn, "admin", userdata)
	fmt.Println(anu)
}

func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Email = "farhan@gmail.com"
	userdata.Password = "riziq"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Email = "farhan@gmail.com"
	userdata.Password = "riziq"
	userdata.Role = "user"
	CreateUser(mconn, "user", userdata)
}

// Admin
func TestInsertUserAdmin(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var userdata User
	userdata.Username = "farhan@gmail.com"
	userdata.Password = "riziq"

	nama := InsertUser(mconn, "admin", userdata)
	fmt.Println(nama)
}

func TestAdminFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata User
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "riziq"
	admindata.Role = "admin"
	CreateUser(mconn, "user", admindata)
}

func TestIsAdminPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata User
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "riziq"

	anu := IsPasswordValid(mconn, "user", admindata)
	fmt.Println(anu)
}

func TestHashAdminFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var admindata Admin
	admindata.Email = "farhan@gmail.com"
	admindata.Password = "riziq"

	filter := bson.M{"email": admindata.Email}
	res := atdb.GetOneDoc[User](mconn, "admin", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(admindata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(admindata.Password, res.Password)
	fmt.Println("Match:   ", match)

}
