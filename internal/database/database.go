package database

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strconv"
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
			return errors.New("Json is not valid.")
		}
	}

	if debug || forceRecreate {
		file, err := os.Create("./database.json")
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
		return initialStruct{}, errors.New("Failed to read bytes from file")
	}

	databaseData := initialStruct{}

	err = json.Unmarshal(readBytes, &databaseData)

	if err != nil {
		return initialStruct{}, errors.New("Failed unmarsheling bytes to Chirp slice")
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
		return nil, errors.New("Failed to read bytes from file")
	}

	dbStruct := initialStruct{}

	err = json.Unmarshal(readBytes, &dbStruct)

	if err != nil {
		return nil, errors.New("Failed unmarsheling bytes to Chirp slice")
	}

	if len(dbStruct.Chirps) == 0 {
		dbStruct.Chirps = make(Chirps, 0)
	}

	return dbStruct.Chirps, nil
}

func SaveChirpToDisk(createChirpData CreateChirpData) error {
	f, err := os.OpenFile("./database.json", os.O_CREATE, 0660)
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
	Email string `json:"email"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type Users = map[string]User

func CreateUser(userCreateData UserCreateData) error {
	f, err := os.OpenFile("./database.json", os.O_CREATE, 0660)
	defer f.Close()

	if err != nil {
		return err
	}

	users, err := GetUsers()

	if err != nil {
		return err
	}

	usersAny := make(map[string]interface{}, len(users))

	for key, value := range users {
		usersAny[key] = value
	}

	nextId := getNextId(usersAny)

	newUser := User{
		Email: userCreateData.Email,
		ID:    nextId,
	}

	users[nextId] = newUser

	newDbDataBytes, err := replaceUsers(users)

	if err != nil {
		return err
	}

	_, err = f.Write(newDbDataBytes)

	return err
}

func replaceUsers(newUsers Users) ([]byte, error) {
	databaseData, err := GetDatabaseData()

	if err != nil {
		return nil, err
	}

	databaseData.Users = newUsers

	databaseDataBytes, err := json.Marshal(databaseData)

	if err != nil {
		return nil, err
	}

	return databaseDataBytes, nil
}

func GetUser(id string) (User, error) {
	users, err := GetUsers()

	if err != nil {
		return User{}, err
	}

	user := users[id]

	return user, nil
}

func GetUsers() (Users, error) {
	readBytes, err := os.ReadFile("./database.json")

	if err != nil {
		return nil, err
	}

	dbStruct := initialStruct{}

	err = json.Unmarshal(readBytes, &dbStruct)

	if err != nil {
		return nil, errors.New("Failed unmarsheling bytes to Chirp slice")
	}

	if len(dbStruct.Users) == 0 {
		dbStruct.Users = make(Users, 0)
	}

	return dbStruct.Users, nil
}
