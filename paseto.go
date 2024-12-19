package pasetobackend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GCFInsertWorkout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collworkout string, r *http.Request) string {
	var response Credential
	response.Status = false

	// Set koneksi ke database
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname) // Mengembalikan *mongo.Database
	var admindata Admin
	gettoken := r.Header.Get("Login")

	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Proses token Login
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		admindata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			admin2 := FindAdmin(mconn, colladmin, admindata)
			if admin2.Role == "admin" {
				var workoutData Workout
				err := json.NewDecoder(r.Body).Decode(&workoutData)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Generate NumberID
					workoutData.NumberID = GenerateNumberID(mconn, collworkout)

					// Insert data ke MongoDB
					insertedID := insertWorkout(mconn, collworkout, workoutData)

					// Validasi hasil insert
					if insertedID == primitive.NilObjectID {
						response.Message = "Insert workout gagal"
					} else {
						response.Status = true
						response.Message = fmt.Sprintf("Berhasil Insert Workout. ID: %v, NumberID: %d", insertedID, workoutData.NumberID)
					}
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
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
			dataworkout := GetAllHp(conn, colname)
			if dataworkout == nil {
				req.Status = false
				req.Message = "Data Workout tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Workout berhasil diambil"
				req.Data = dataworkout
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

// WORKOUT

func GetWorkoutData(PublicKey, MongoConnStringEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)

	// Menggunakan fungsi SetConnection yang Anda miliki
	conn := SetConnection(MongoConnStringEnv, dbname)
	if conn == nil {
		req.Status = false
		req.Message = "Failed to connect to MongoDB"
		return ReturnStringStruct(req)
	}

	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Verifikasi token login
		_, err := DecodeGetHp(os.Getenv(PublicKey), tokenlogin) // Sesuaikan fungsi DecodeGetHp
		if err != nil {
			req.Status = false
			req.Message = "Invalid token: " + tokenlogin
		} else {
			// Mengambil data workout dari MongoDB
			dataWorkout := GetAllWorkout(conn, colname)
			if dataWorkout == nil {
				req.Status = false
				req.Message = "No workout data found"
			} else {
				req.Status = true
				req.Message = "Workout data retrieved successfully"
				req.Data = dataWorkout
			}
		}
	}
	return ReturnStringStruct(req)
}

func GCFUpdateWorkout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collworkout string, r *http.Request) string {
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
				var workoutData Workout
				err := json.NewDecoder(r.Body).Decode(&workoutData)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Update data workout di MongoDB
					UpdatedWorkout(mconn, collworkout, bson.M{"number_id": workoutData.NumberID}, workoutData)
					response.Status = true
					response.Message = "Berhasil Update Workout"
					return GCFReturnStruct(CreateResponse(true, "Success Update Workout", workoutData))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFDeleteWorkout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collworkout string, r *http.Request) string {
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
				var workoutData Workout
				err := json.NewDecoder(r.Body).Decode(&workoutData)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Delete data workout di MongoDB
					DeleteWorkout(mconn, collworkout, workoutData)
					response.Status = true
					response.Message = "Berhasil Delete Workout"
				}
			} else {
				response.Message = "Anda tidak dapat Delete data karena bukan admin"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFGetWorkoutByID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	// Membuat koneksi ke database
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Parsing ID dari request body
	var workoutData Workout
	err := json.NewDecoder(r.Body).Decode(&workoutData)
	if err != nil {
		return GCFReturnStruct(CreateResponse(false, "Error parsing application/json: "+err.Error(), nil))
	}

	// Mendapatkan data workout berdasarkan ID
	workout := GetWorkoutByID(mconn, collectionname, workoutData.ID)
	if workout != (Workout{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get Workout By ID", workout))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get Workout By ID", nil))
	}
}
