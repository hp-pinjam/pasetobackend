package pasetobackend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
				resp.Message = "Selamat Datang"
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
		// Validasi data yang wajib diisi
		if userdata.Username == "" || userdata.Email == "" || userdata.Password == "" {
			resp.Message = "Username, Email, dan Password tidak boleh kosong"
			return ReturnStringStruct(resp)
		}

		if userdata.Height <= 0 || userdata.Weight <= 0 || userdata.Age <= 0 {
			resp.Message = "Height, Weight, dan Age harus lebih besar dari 0"
			return ReturnStringStruct(resp)
		}

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
				"height":   userdata.Height,
				"weight":   userdata.Weight,
				"age":      userdata.Age,
			})
			if err != nil {
				resp.Message = "Gagal Input data ke user: " + err.Error()
			} else {
				resp.Status = true
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
					// Set ID secara otomatis jika tidak diberikan
					if workoutData.ID.IsZero() {
						workoutData.ID = primitive.NewObjectID()
					}

					// Insert data ke MongoDB
					insertedID := insertWorkout(mconn, collworkout, workoutData)
					if insertedID != nil {
						response.Status = true
						response.Message = fmt.Sprintf("Berhasil Insert Workout. ID: %v", insertedID)
					} else {
						response.Message = "Gagal Insert Workout"
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

	// Membuat koneksi ke MongoDB
	conn := SetConnection(MongoConnStringEnv, dbname)
	if conn == nil {
		req.Status = false
		req.Message = "Failed to connect to MongoDB"
		return ReturnStringStruct(req)
	}

	// Cek token login
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
		return ReturnStringStruct(req)
	}

	// Verifikasi token login
	_, err := DecodeGetHp(os.Getenv(PublicKey), tokenlogin) // Sesuaikan fungsi DecodeGetHp
	if err != nil {
		req.Status = false
		req.Message = "Invalid token: " + tokenlogin
		return ReturnStringStruct(req)
	}

	// Mengambil semua data workout dari koleksi
	dataWorkout, err := GetAllWorkout(conn, colname)
	if err != nil {
		req.Status = false
		req.Message = "Failed to retrieve workout data: " + err.Error()
		return ReturnStringStruct(req)
	}

	if len(dataWorkout) == 0 {
		req.Status = false
		req.Message = "No workout data found"
		return ReturnStringStruct(req)
	}

	// Jika berhasil
	req.Status = true
	req.Message = "Workout data retrieved successfully"
	req.Data = dataWorkout
	return ReturnStringStruct(req)
}

func GCFUpdateWorkout(publickey, MONGOCONNSTRINGENV, dbname, colladmin, collworkout string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname) // Koneksi ke database
	var admindata Admin

	// Cek header token
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
		return GCFReturnStruct(response)
	}

	// Decode token Login
	checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
	admindata.Email = checktoken
	if checktoken == "" {
		response.Message = "Kamu kayaknya belum punya akun"
		return GCFReturnStruct(response)
	}

	// Cek apakah user adalah admin
	admin2 := FindAdmin(mconn, colladmin, admindata)
	if admin2.Role != "admin" {
		response.Message = "Anda tidak dapat Update data karena bukan admin"
		return GCFReturnStruct(response)
	}

	// Decode data Workout dari body request
	var workoutData Workout
	err := json.NewDecoder(r.Body).Decode(&workoutData)
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Pastikan NumberID tersedia untuk filter
	if workoutData.NumberID == 0 {
		response.Message = "NumberID is required to update workout"
		return GCFReturnStruct(response)
	}

	// Filter dan data update
	filter := bson.M{"number_id": workoutData.NumberID}
	update := bson.M{
		"$set": bson.M{
			"name":       workoutData.Name,
			"gif":        workoutData.Gif,
			"repetition": workoutData.Repetition,
			"calories":   workoutData.Calories,
			"status":     workoutData.Status,
		},
	}

	// Update data di koleksi workout
	collection := mconn.Collection(collworkout)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.Message = "Failed to update workout: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Cek apakah dokumen ditemukan dan diperbarui
	if result.MatchedCount == 0 {
		response.Message = "No workout found with the given NumberID"
		return GCFReturnStruct(response)
	}

	if result.ModifiedCount == 0 {
		response.Message = "Workout data is already up-to-date"
		return GCFReturnStruct(response)
	}

	// Ambil data yang diperbarui
	var updatedWorkout Workout
	err = collection.FindOne(ctx, filter).Decode(&updatedWorkout)
	if err != nil {
		response.Message = "Workout updated, but failed to fetch updated data: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Jika berhasil
	response.Status = true
	response.Message = "Workout updated successfully"
	response.Data = updatedWorkout
	return GCFReturnStruct(response)
}

func GCFDeleteWorkout(publickey, MONGOCONNSTRINGENV, dbname, colluser, collworkout string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname) // Koneksi ke database
	var userdata User

	// Cek header token
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
		return GCFReturnStruct(response)
	}

	// Decode token Login
	checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
	userdata.Username = checktoken
	if checktoken == "" {
		response.Message = "Kamu kayaknya belum punya akun"
		return GCFReturnStruct(response)
	}

	// Cek apakah user adalah admin
	// user2 := FindUser(mconn, colluser, userdata)
	// if user2.Role != "user" {
	// 	response.Message = "Anda tidak dapat Delete data karena bukan admin"
	// 	return GCFReturnStruct(response)
	// }

	// Decode data Workout dari body request
	var workoutData Workout
	err := json.NewDecoder(r.Body).Decode(&workoutData)
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Validasi NumberID untuk filter
	if workoutData.NumberID == 0 {
		response.Message = "NumberID is required to delete workout"
		return GCFReturnStruct(response)
	}

	// Filter untuk menghapus workout berdasarkan NumberID
	filter := bson.M{"number_id": workoutData.NumberID}
	collection := mconn.Collection(collworkout)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete workout dari MongoDB
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.Message = "Failed to delete workout: " + err.Error()
		return GCFReturnStruct(response)
	}

	// Cek apakah dokumen ditemukan dan dihapus
	if result.DeletedCount == 0 {
		response.Message = "No workout found with the given NumberID"
		return GCFReturnStruct(response)
	}

	// Jika berhasil
	response.Status = true
	response.Message = "Workout deleted successfully"
	return GCFReturnStruct(response)
}

func GCFGetWorkoutByID(PublicKey, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	req := new(Response)

	conn := SetConnection(MONGOCONNSTRINGENV, dbname)
	if conn == nil {
		req.Status = false
		req.Message = "Failed to connect to MongoDB"
		return ReturnStringStruct(req)
	}

	// Cek token login
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
		return ReturnStringStruct(req)
	}

	// Verifikasi token login
	_, err := DecodeGetHp(os.Getenv(PublicKey), tokenlogin) // Sesuaikan fungsi DecodeGetHp
	if err != nil {
		req.Status = false
		req.Message = "Invalid token: " + tokenlogin
		return ReturnStringStruct(req)
	}

	// Parsing NumberID dari body request
	var input struct {
		NumberID int `json:"number_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		req.Status = false
		req.Message = "Error parsing JSON body: " + err.Error()
		return ReturnStringStruct(req)
	}

	// Validasi NumberID
	if input.NumberID == 0 {
		req.Status = false
		req.Message = "NumberID is required"
		return ReturnStringStruct(req)
	}

	// Filter untuk mencari workout berdasarkan NumberID
	filter := bson.M{"number_id": input.NumberID}

	// Koleksi database
	collection := conn.Collection(collectionname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mendapatkan data workout berdasarkan NumberID
	var workout Workout
	err = collection.FindOne(ctx, filter).Decode(&workout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			req.Status = false
			req.Message = "Workout not found"
			return ReturnStringStruct(req)
		}
		req.Status = false
		req.Message = "Failed to retrieve workout: " + err.Error()
		return ReturnStringStruct(req)
	}

	// Jika berhasil
	req.Status = true
	req.Message = "Success: Get Workout By NumberID"
	req.Data = workout
	return ReturnStringStruct(req)
}

func GetUserData(PublicKey, MongoConnStringEnv, dbname, colname string, r *http.Request) string {
	req := new(Credential)

	// Membuat koneksi ke MongoDB
	conn := SetConnection(MongoConnStringEnv, dbname)
	if conn == nil {
		req.Status = false
		req.Message = "Failed to connect to MongoDB"
		return ReturnStringStruct(req)
	}

	// Ambil token login dari header
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
		return ReturnStringStruct(req)
	}

	// Verifikasi token login dan ambil Hp (username)
	userID, err := DecodeGetHp(PublicKey, tokenlogin)
	if err != nil {
		req.Status = false
		req.Message = "Invalid token: " + err.Error()
		return ReturnStringStruct(req)
	}
	fmt.Println("Decoded userID from token:", userID)

	// Query MongoDB berdasarkan Hp (username)
	collection := conn.Collection(colname)
	var userdata User
	err = collection.FindOne(context.Background(), bson.M{"username": userID}).Decode(&userdata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			req.Status = false
			req.Message = "User not found"
			return ReturnStringStruct(req)
		}
		req.Status = false
		req.Message = "Failed to retrieve user data: " + err.Error()
		return ReturnStringStruct(req)
	}

	// Kosongkan password sebelum mengembalikan data user
	userdata.Password = ""

	// Jika berhasil
	req.Status = true
	req.Message = "User data retrieved successfully"
	req.Data = userdata
	return ReturnStringStruct(req)
}

func Registrasi(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname) // conn bertipe *mongo.Database
	admindata := new(User)

	// Decode request body ke dalam struct User
	err := json.NewDecoder(r.Body).Decode(&admindata)
	if err != nil {
		resp.Message = "Error parsing application/json: " + err.Error()
	} else {
		// Validasi data yang wajib diisi
		if admindata.Username == "" || admindata.Email == "" || admindata.Password == "" {
			resp.Message = "Username, Email, dan Password tidak boleh kosong"
			return ReturnStringStruct(resp)
		}

		if admindata.Height <= 0 || admindata.Weight <= 0 || admindata.Age <= 0 {
			resp.Message = "Height, Weight, dan Age harus lebih besar dari 0"
			return ReturnStringStruct(resp)
		}

		// Hash password sebelum menyimpan ke database
		hash, err := HashPass(admindata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password: " + err.Error()
		} else {
			// Masukkan data user ke dalam koleksi MongoDB
			collection := conn.Collection("admin") // Ambil koleksi "user" dari database
			_, err := collection.InsertOne(context.Background(), bson.M{
				"username": admindata.Username,
				"email":    admindata.Email,
				"password": hash, // Simpan password yang sudah di-hash
				"height":   admindata.Height,
				"weight":   admindata.Weight,
				"age":      admindata.Age,
			})
			if err != nil {
				resp.Message = "Gagal Input data: " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Berhasil Input data"
			}
		}
	}

	// Konversi respons ke format string
	response := ReturnStringStruct(resp)
	return response
}
