package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegStage int

const (
	RegStage1 RegStage = iota + 1
	RegStage2
	RegStage3
)

type ContextKey struct{}

type Pagination struct {
	Offset int
	Limit  int
}

func GetPagination(page, pageSize int) Pagination {
	return Pagination{
		Offset: (page - 1) * pageSize,
		Limit:  pageSize,
	}
}

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

type queryFromStruct struct {
	Query string
	Args  []any
}

func GetUpdateQueryFromStruct(s any, tableName string) (*queryFromStruct, error) {
	mapData, err := StructToMap(s)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("UPDATE %s\nSET version = version + 1, \n", tableName)

	keys := []string{}
	args := []any{}

	for k, v := range mapData {
		if k == "id" || k == "version" || v == nil {
			continue
		}

		keys = append(keys, k)
		args = append(args, v)
	}

	for i, v := range keys {
		if i+1 == len(keys) {
			query += fmt.Sprintf("%s = $%d \n", v, i+1)
		} else {
			query += fmt.Sprintf("%s = $%d, \n", v, i+1)
		}
	}

	version := mapData["version"].(float64)
	query += fmt.Sprintf("\nWHERE id = '%s' AND version = %d\nRETURNING version", mapData["id"], int64(version))

	return &queryFromStruct{Query: query, Args: args}, nil
}

func GeneratePasswordHash(str string) ([]byte, error) {
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

func GetPage(v int) int {
	if v == 0 {
		return 1
	}

	return v
}

func GetPageSize(v int) int {
	if v == 0 {
		return 20
	}

	return v
}
