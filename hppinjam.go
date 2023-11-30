package pasetobackend

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
)

func GFCPostHandlerUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false

	// Mendapatkan data yang diterima dari permintaan HTTP POST
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		// Menggunakan variabel MONGOCONNSTRINGENV untuk string koneksi MongoDB
		mongoConnStringEnv := MONGOCONNSTRINGENV

		mconn := SetConnection(mongoConnStringEnv, dbname)

		// Lakukan pemeriksaan kata sandi menggunakan bcrypt
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			Response.Message = "Selamat Datang"
		} else {
			Response.Message = "Password Salah"
		}
	}

	// Mengirimkan respons sebagai JSON
	responseJSON, _ := json.Marshal(Response)
	return string(responseJSON)
}

// Login User NPM
// func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	var Response Credential
// 	Response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datauser User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 	} else {
// Assuming either email or npm is provided in the request
// if IsPasswordValid(mconn, collectionname, datauser) {
// 	Response.Status = true
// Using NPM as identifier, you can modify this as needed
// 			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
// 			if err != nil {
// 				Response.Message = "Gagal Encode Token : " + err.Error()
// 			} else {
// 				Response.Message = "Selamat Datang"
// 				Response.Token = tokenstring
// 			}
// 		} else {
// 			Response.Message = "Username atau Password Salah"
// 		}
// 	}

// 	return GCFReturnStruct(Response)
// }

// Login User Email
// func GCFPostHandlerEmail(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	var Response Credential
// 	Response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datauser User
// 	err := json.NewDecoder(r.Body).Decode(&datauser)
// 	if err != nil {
// 		Response.Message = "error parsing application/json: " + err.Error()
// 	} else {
// Assuming either email or npm is provided in the request
// if IsPasswordValidEmail(mconn, collectionname, datauser) {
// 	Response.Status = true
// Using NPM as identifier, you can modify this as needed
// 			tokenstring, err := watoken.Encode(datauser.Email, os.Getenv(PASETOPRIVATEKEYENV))
// 			if err != nil {
// 				Response.Message = "Gagal Encode Token : " + err.Error()
// 			} else {
// 				Response.Message = "Selamat Datang"
// 				Response.Token = tokenstring
// 			}
// 		} else {
// 			Response.Message = "Email atau Password Salah"
// 		}
// 	}

// 	return GCFReturnStruct(Response)
// }

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// Login Admin
func LoginAdmin(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var dataadmin Admin
	err := json.NewDecoder(r.Body).Decode(&dataadmin)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValidAdmin(mconn, Colname, dataadmin) {
			tokenstring, err := watoken.Encode(dataadmin.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

// Register User
func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	userdata := new(User)
	resp.Status = false
	conn := GetConnectionMongo(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := HashPassword(userdata.PasswordHash)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertUserdata(conn, userdata.Username, userdata.Password, hash)
		resp.Message = "Berhasil Input data"
	}
	response := ReturnStringStruct(resp)
	return response
}

// Post Parkiran
// func GCFCreateParkiran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var dataparkiran Parkiran
// 	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	if err := CreateParkiran(mconn, collectionname, dataparkiran); err != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Success Create Parkiran", dataparkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Create Parkiran", dataparkiran))
// 	}
// }

// Delete Parkiran
// func GCFDeleteParkiran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	var dataparkiran Parkiran
// 	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	if err := DeleteParkiran(mconn, collectionname, dataparkiran); err != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Success Delete Parkiran", dataparkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Delete Parkiran", dataparkiran))
// 	}
// }

// Update Parkiran
// func GCFUpdateParkiran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	var dataparkiran Parkiran
// 	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	if err := UpdatedParkiran(mconn, collectionname, bson.M{"id": dataparkiran.ID}, dataparkiran); err != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Success Update Parkiran", dataparkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Update Parkiran", dataparkiran))
// 	}
// }

// Get All Parkiran
// func GCFGetAllParkiran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	dataparkiran := GetAllParkiran(mconn, collectionname)
// 	if dataparkiran != nil {
// 		return GCFReturnStruct(CreateResponse(true, "success Get All Parkiran", dataparkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Get All Parkiran", dataparkiran))
// 	}
// }

// Get All Parkiran By Id
// func GCFGetAllParkiranID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	var dataparkiran Parkiran
// 	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	parkiran := GetAllParkiranID(mconn, collectionname, dataparkiran)
// 	if parkiran != (Parkiran{}) {
// 		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Parkiran", dataparkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Parkiran", dataparkiran))
// 	}
// }
