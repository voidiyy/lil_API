package db_psql

import (
	"fmt"
	types "gigaAPI/internal/type"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func HashPass(p string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hash pass error: ", err)
		return p
	}

	return string(hash)
}

func HashCompare(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func CreateTestWorker() *types.Worker {
	pass := genPassORName(8)
	hash := HashPass(pass)
	w := &types.Worker{
		ID:        uuid.New(),
		Name:      genPassORName(6),
		Password:  hash,
		Email:     genPassORName(5) + "@mail.com",
		Role:      genRole(),
		CreatedAt: time.Time{},
	}

	return w
}

func genPassORName(i int) string {
	chars := "abcdefghjklmnopqrstuvwxyz"
	str := make([]byte, i)
	for k, _ := range str {
		str[k] = chars[rand.Intn(len(chars))]
	}

	return string(str)
}

func genRole() string {
	var rls = []string{"Junior", "Middle", "Senior"}

	return rls[rand.Intn(len(rls))]
}
