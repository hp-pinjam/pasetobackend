package pasetobackend

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// <--- ini Login & Register Admin --->
func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var dataadmin Admin
	err := json.NewDecoder(r.Body).Decode(&dataadmin)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, Colname, dataadmin) {
			tokenstring, err := watoken.Encode(dataadmin.Email, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang SUPERADMIN"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	admindata := new(Admin)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&admindata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := HashPass(admindata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertAdmindata(conn, admindata.Email, admindata.Role, hash)
		resp.Message = "Berhasil Input data"
	}
	response := ReturnStringStruct(resp)
	return response
}

func RegisterUser(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname) // conn bertipe *mongo.Database
	userdata := new(User)

	// Decode request body ke dalam struct User
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "Error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		// Hash password sebelum menyimpan ke database
		hash, err := HashPass(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password: " + err.Error()
		} else {
			// Masukkan data user ke dalam koleksi MongoDB
			collection := conn.Collection("user") // Ambil koleksi "user" dari database
			_, err := collection.InsertOne(context.Background(), bson.M{
				"username": userdata.Username,
				"email":    userdata.Email,
				"password": hash, // Simpan password yang sudah di-hash
			})
			if err != nil {
				resp.Message = "Gagal Input data ke user: " + err.Error()
			} else {
				resp.Message = "Berhasil Input data ke user"
			}
		}
	}

	// Konversi respons ke format string
	response := ReturnStringStruct(resp)
	return response
}

// <--- ini hp --->

// hp post
func GCFInsertHp(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collhp string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datahp Hp
				err := json.NewDecoder(r.Body).Decode(&datahp)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertHp(mconn, collhp, Hp{
						Nomorid:     datahp.Nomorid,
						Title:       datahp.Title,
						Description: datahp.Description,
						Image:       datahp.Image,
						Status:      datahp.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Hp"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

// delete Hp
func GCFDeleteHp(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collhp string, r *http.Request) string {

	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datahp Hp
				err := json.NewDecoder(r.Body).Decode(&datahp)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteHp(mconn, collhp, datahp)
					respon.Status = true
					respon.Message = "Berhasil Delete Hp"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(respon)
}

// update Hp
func GCFUpdateHp(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collhp string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var admindata Admin

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var datahp Hp
				err := json.NewDecoder(r.Body).Decode(&datahp)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()

				} else {
					UpdatedHp(mconn, collhp, bson.M{"id": datahp.ID}, datahp)
					response.Status = true
					response.Message = "Berhasil Update Hp"
					GCFReturnStruct(CreateResponse(true, "Success Update Hp", datahp))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}

		}
	}
	return GCFReturnStruct(response)
}

// get all hp
func GCFGetAllHp(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datahp := GetAllHp(mconn, collectionname)
	if datahp != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Hp", datahp))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Hp", datahp))
	}
}

func GCFGetAllHpg(publickey, Mongostring, dbname, colname string, r *http.Request) string {
	resp := new(Credential)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Exist"
	} else {
		existing := IsExist(tokenlogin, os.Getenv(publickey))
		if !existing {
			resp.Status = false
			resp.Message = "Kamu kayaknya belum punya akun"
		} else {
			koneksyen := SetConnection(Mongostring, dbname)
			datahp := GetAllHp(koneksyen, colname)
			yas, _ := json.Marshal(datahp)
			resp.Status = true
			resp.Message = "Data Berhasil diambil"
			resp.Token = string(yas)
		}
	}
	return ReturnStringStruct(resp)
}

func GetAllDataHps(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetHp(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Tidak ada data  " + tokenlogin
		} else {
			// Langsung ambil data hp
			datahp := GetAllHp(conn, colname)
			if datahp == nil {
				req.Status = false
				req.Message = "Data hp tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Hp berhasil diambil"
				req.Data = datahp
			}
		}
	}
	return ReturnStringStruct(req)
}

// get all hp by id
func GCFGetAllHpID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datahp Hp
	err := json.NewDecoder(r.Body).Decode(&datahp)
	if err != nil {
		return err.Error()
	}

	hp := GetAllHpID(mconn, collectionname, datahp)
	if hp != (Hp{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Hp", datahp))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Hp", datahp))
	}
}

// <--- ini about --->

// about post
// func GCFInsertAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
// 	var response Credential
// 	response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var authdata Admin

// 	gettoken := r.Header.Get("Login")

// 	if gettoken == "" {
// 		response.Message = "Header Login Not Exist"
// 	} else {
// Process the request with the "Login" token
// 		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
// 		authdata.Email = checktoken
// 		if checktoken == "" {
// 			response.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			auth2 := FindAdmin(mconn, colladmin, authdata)
// 			if auth2.Role == "admin" {

// 				var dataabout About
// 				err := json.NewDecoder(r.Body).Decode(&dataabout)
// 				if err != nil {
// 					response.Message = "Error parsing application/json: " + err.Error()
// 				} else {
// 					InsertAbout(mconn, collabout, About{
// 						ID:          dataabout.ID,
// 						Title:       dataabout.Title,
// 						Description: dataabout.Description,
// 						Image:       dataabout.Image,
// 						Status:      dataabout.Status,
// 					})
// 					response.Status = true
// 					response.Message = "Berhasil Insert About"
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Insert data karena bukan admin"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)

// }

// delete about
// func GCFDeleteAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
// 	var response Credential
// 	response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var authdata Admin

// 	gettoken := r.Header.Get("Login")

// 	if gettoken == "" {
// 		response.Message = "Header Login Not Exist"
// 	} else {
// Process the request with the "Login" token
// 		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
// 		authdata.Email = checktoken
// 		if checktoken == "" {
// 			response.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			auth2 := FindAdmin(mconn, colladmin, authdata)
// 			if auth2.Role == "admin" {

// 				var dataabout About
// 				err := json.NewDecoder(r.Body).Decode(&dataabout)
// 				if err != nil {
// 					response.Message = "Error parsing application/json: " + err.Error()
// 				} else {
// 					DeleteAbout(mconn, collabout, dataabout)
// 					response.Status = true
// 					response.Message = "Berhasil Delete About"
// 					CreateResponse(true, "Success Delete About", dataabout)
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Delete data karena bukan admin"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

// update about
// func GCFUpdateAbout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collabout string, r *http.Request) string {
// 	var response Credential
// 	response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var authdata Admin

// 	gettoken := r.Header.Get("Login")

// 	if gettoken == "" {
// 		response.Message = "Header Login Not Exist"
// 	} else {
// Process the request with the "Login" token
// 		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
// 		authdata.Email = checktoken
// 		if checktoken == "" {
// 			response.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			auth2 := FindAdmin(mconn, colladmin, authdata)
// 			if auth2.Role == "admin" {
// 				var dataabout About
// 				err := json.NewDecoder(r.Body).Decode(&dataabout)
// 				if err != nil {
// 					response.Message = "Error parsing application/json: " + err.Error()
// 				} else {
// 					UpdatedAbout(mconn, collabout, bson.M{"id": dataabout.ID}, dataabout)
// 					response.Status = true
// 					response.Message = "Berhasil Update Hp"
// 					CreateResponse(true, "Success Update About", dataabout)
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Update data karena bukan admin"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

// get all about
// func GCFGetAllAbout(MONGOCONNSTRINGENV, dbname, collectionname string) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	dataabout := GetAllAbout(mconn, collectionname)
// 	if dataabout != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Berhasil Get All About", dataabout))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Gagal Get All About", dataabout))
// 	}
// }

// func GCFGetAllAboutt(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datahp := GetAllAbout(koneksyen, colname)
// 			yas, _ := json.Marshal(datahp)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

// <--- ini contact --->

// contact post
// func GCFInsertContact(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var datacontact Contact
// 	err := json.NewDecoder(r.Body).Decode(&datacontact)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	if err := InsertContact(mconn, collectionname, datacontact); err != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Success Create Contact", datacontact))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Create Contact", datacontact))
// 	}
// }

// get all contact
// func GCFGetAllContact(MONGOCONNSTRINGENV, dbname, collectionname string) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	datacontact := GetAllContact(mconn, collectionname)
// 	if datacontact != nil {
// 		return GCFReturnStruct(CreateResponse(true, "success Get All Contact", datacontact))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Get All Contact", datacontact))
// 	}
// }

// func GCFGetAllContactt(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datahp := GetAllContact(koneksyen, colname)
// 			yas, _ := json.Marshal(datahp)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }

//crawling

// get all crawling
// func GCFGetAllCrawling(MONGOCONNSTRINGENV, dbname, collectionname string) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	datacrawling := GetAllCrawling(mconn, collectionname)
// 	if datacrawling != nil {
// 		return GCFReturnStruct(CreateResponse(true, "success Get All Contact", datacrawling))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Get All Contact", datacrawling))
// 	}
// }

// func GCFGetAllCrawlingg(publickey, Mongostring, dbname, colname string, r *http.Request) string {
// 	resp := new(Credential)
// 	tokenlogin := r.Header.Get("Login")
// 	if tokenlogin == "" {
// 		resp.Status = false
// 		resp.Message = "Header Login Not Exist"
// 	} else {
// 		existing := IsExist(tokenlogin, os.Getenv(publickey))
// 		if !existing {
// 			resp.Status = false
// 			resp.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			koneksyen := SetConnection(Mongostring, dbname)
// 			datahp := GetAllCrawling(koneksyen, colname)
// 			yas, _ := json.Marshal(datahp)
// 			resp.Status = true
// 			resp.Message = "Data Berhasil diambil"
// 			resp.Token = string(yas)
// 		}
// 	}
// 	return ReturnStringStruct(resp)
// }
