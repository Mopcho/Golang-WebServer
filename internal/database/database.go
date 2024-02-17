package database

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"slices"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type CreateChirpData struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type Chirps = map[string]Chirp

type initialStruct struct {
	Chirps Chirps `json:"chirps"`
	Users  Users  `json:"users"`
}

func setupCleanDb(file *os.File) error {
	initDbStruct := initialStruct{}

	dbBytes, err := json.Marshal(initDbStruct)

	if err != nil {
		return err
	}

	file.Write(dbBytes)

	return nil
}

func SetupDataBase(debug bool) error {
	forceRecreate := false
	stat, err := os.Stat("./database.json")

	if err != nil || stat.Size() == 0 {
		// File does not exist or it has no data inside
		forceRecreate = true
	} else {
		// File exists and has data lets check if it has valid json
		byteDataFromFile, err := os.ReadFile("./database.json")

		if err != nil {
			return err
		}

		isValidJson := json.Valid(byteDataFromFile)

		if !isValidJson && !debug {
			return errors.New("json is not valid")
		}
	}

	if debug || forceRecreate {
		file, err := os.Create("./database.json")

		if err != nil {
			return err
		}

		defer file.Close()

		if err != nil {
			return err
		}

		err = setupCleanDb(file)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetDatabaseData() (initialStruct, error) {
	readBytes, err := os.ReadFile("./database.json")

	if err != nil {
		return initialStruct{}, errors.New("failed to read bytes from file")
	}

	databaseData := initialStruct{}

	err = json.Unmarshal(readBytes, &databaseData)

	if err != nil {
		return initialStruct{}, errors.New("failed unmarsheling bytes to Chirp slice")
	}

	return databaseData, nil
}

func GetOneChirp(id string) (Chirp, error) {
	chirps, err := GetChirpsFromDisk()

	if err != nil {
		return Chirp{}, err
	}

	return chirps[id], nil
}

func GetChirpsFromDisk() (Chirps, error) {
	readBytes, err := os.ReadFile("./database.json")

	if err != nil {
		return nil, errors.New("failed to read bytes from file")
	}

	dbStruct := initialStruct{}

	err = json.Unmarshal(readBytes, &dbStruct)

	if err != nil {
		return nil, errors.New("failed unmarsheling bytes to Chirp slice")
	}

	if len(dbStruct.Chirps) == 0 {
		dbStruct.Chirps = make(Chirps, 0)
	}

	return dbStruct.Chirps, nil
}

func SaveChirpToDisk(createChirpData CreateChirpData) error {
	f, err := os.OpenFile("./database.json", os.O_CREATE, 0660)

	if err != nil {
		return err
	}

	defer f.Close()

	readChirps, err := GetChirpsFromDisk()

	if err != nil {
		return err
	}

	chirpsAny := make(map[string]interface{}, len(readChirps))

	for key, value := range readChirps {
		chirpsAny[key] = value
	}

	nextId := getNextId(chirpsAny)

	chirp := Chirp{
		ID:   nextId,
		Body: createChirpData.Body,
	}

	readChirps[nextId] = chirp

	newChirpsBytes, err := replaceChirpsInDbStruct(readChirps)

	if err != nil {
		return err
	}

	_, err = f.Write(newChirpsBytes)
	return err
}

func replaceChirpsInDbStruct(newChirps Chirps) ([]byte, error) {
	dbData, err := GetDatabaseData()

	if err != nil {
		return nil, err
	}

	dbData.Chirps = newChirps

	dbDataBytes, err := json.Marshal(dbData)

	if err != nil {
		return nil, err
	}

	return dbDataBytes, nil
}

func getNextId(mapData map[string]interface{}) string {
	ids := make([]int, 0)

	for key := range mapData {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		ids = append(ids, keyInt)
	}

	if len(ids) == 0 {
		return "0"
	}

	nextId := slices.Max(ids)

	nextIdString := strconv.Itoa(nextId + 1)

	return nextIdString
}

type UserCreateData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEditData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users = map[string]User

func CreateUser(userCreateData UserCreateData) error {
	users, err := GetUsers()

	if err != nil {
		return err
	}

	usersAny := make(map[string]interface{}, len(users))

	for key, value := range users {
		usersAny[key] = value
	}

	nextId := getNextId(usersAny)

	hashedPassword, err := hashPassword(userCreateData.Password)

	if err != nil {
		return err
	}

	newUser := User{
		Email:    userCreateData.Email,
		ID:       nextId,
		Password: hashedPassword,
	}

	users[nextId] = newUser

	err = replaceUsers(users)

	return err
}

func replaceUsers(newUsers Users) error {
	f, err := os.OpenFile("./database.json", os.O_CREATE, 0660)

	if err != nil {
		return err
	}

	defer f.Close()

	databaseData, err := GetDatabaseData()

	if err != nil {
		return err
	}

	databaseData.Users = newUsers

	databaseDataBytes, err := json.Marshal(databaseData)

	if err != nil {
		return err
	}

	_, err = f.Write(databaseDataBytes)

	return err
}

func GetUserById(id string) (User, error) {
	users, err := GetUsers()

	if err != nil {
		return User{}, err
	}

	user := users[id]

	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	users, err := GetUsers()

	if err != nil {
		return User{}, err
	}

	for key, value := range users {
		if value.Email == email {
			return users[key], nil
		}
	}

	return User{}, errors.New("User not found")
}

func GetUsers() (Users, error) {
	readBytes, err := os.ReadFile("./database.json")

	if err != nil {
		return nil, err
	}

	dbStruct := initialStruct{}

	err = json.Unmarshal(readBytes, &dbStruct)

	if err != nil {
		return nil, errors.New("failed unmarsheling bytes to Chirp slice")
	}

	if len(dbStruct.Users) == 0 {
		dbStruct.Users = make(Users, 0)
	}

	return dbStruct.Users, nil
}

func EditUser(userEditData UserEditData) error {
	dbUser, err := GetUserById(userEditData.ID)

	if err != nil {
		return err
	}

	fields := reflect.VisibleFields(reflect.TypeOf(dbUser))

	for _, field := range fields {
		fieldName := field.Name
		orgFieldValue := reflect.ValueOf(dbUser).FieldByName(fieldName)
		newFieldValue := reflect.ValueOf(userEditData).FieldByName(fieldName)

		if newFieldValue.IsZero() {
			continue
		}

		if fieldName == "ID" {
			continue
		}

		if orgFieldValue == newFieldValue {
			continue
		}

		if fieldName == "Password" {
			passwordsAreNotTheSame := ComparePassword(orgFieldValue.Interface().(string), newFieldValue.Interface().(string))

			if passwordsAreNotTheSame == nil {
				continue
			}

			newHash, err := hashPassword(newFieldValue.Interface().(string))

			if err != nil {
				return err
			}

			reflect.ValueOf(&dbUser).Elem().FieldByName(fieldName).SetString(newHash)
			continue
		}

		if reflect.ValueOf(&dbUser).Elem().FieldByName(fieldName).Type() != nil {
			reflect.ValueOf(&dbUser).Elem().FieldByName(fieldName).Set(newFieldValue)
			continue
		}
	}

	users, err := GetUsers()

	if err != nil {
		return err
	}

	users[dbUser.ID] = dbUser

	err = replaceUsers(users)

	return err
}

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	return string(hashBytes), err
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
