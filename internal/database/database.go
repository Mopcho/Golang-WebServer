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
}

func SetupDataBase() error {
	file, err := os.Create("./database.json")
	defer file.Close()

	if err != nil {
		return err
	}

	initDbStruct := initialStruct{}

	dbBytes, err := json.Marshal(initDbStruct)

	if err != nil {
		return err
	}

	file.Write(dbBytes)

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

	nextId := getNextChirpId(readChirps)

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

func getNextChirpId(chirps Chirps) string {
	ids := make([]int, 0)

	for key := range chirps {
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
