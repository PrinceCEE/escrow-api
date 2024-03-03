package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type RegStage int

const (
	RegStage1 RegStage = iota + 1
	RegStage2
	RegStage3
)

func StructToMap(data any) (map[string]any, error) {
	mapData := make(map[string]any)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &mapData)
	if err != nil {
		return nil, err
	}

	return mapData, nil
}

func IsInSlice(slice []any, item any) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

func GetUpdateQueryFromStruct(s any, tableName string) (string, error) {
	mapData, err := StructToMap(s)
	if err != nil {
		return "", err
	}

	t := reflect.TypeOf(s)
	userType := reflect.TypeOf(&models.User{})

	query := fmt.Sprintf("UPDATE %s\nSET ", tableName)
	for k, v := range mapData {
		if k == "id" {
			continue
		}

		if k == "version" {
			query += "version = version + 1, "
		}

		if t == userType {
			if (k == "email" && v == "") || (k == "phone_number" && v == "") {
				continue
			}
		}

		query += fmt.Sprintf("%s = %s, ", k, v)
	}
	query += fmt.Sprintf("\nWHERE id = %s AND version = %s\nRETURNING version", mapData["id"], mapData["version"])

	return query, nil
}

func GenerateHash(str string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
}

func ComparePassword(pwd string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pwd))
}

func GenerateRandomNumber() string {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)

	r := rand.New(source)
	v := r.Intn(10000)

	return fmt.Sprintf("%04d", v)
}

func Background(fn func()) {
	go func() {
		fn()
	}()
}
