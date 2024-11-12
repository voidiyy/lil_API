package http_server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gigaAPI/internal/db_psql"
	types "gigaAPI/internal/type"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var URI = "https://localhost:1234/db_psql/worker"

func Request(method, urlPath string, buff []byte) (*http.Response, error) {
	// Створюємо запит
	req, err := http.NewRequest(method, URI+urlPath, bytes.NewReader(buff))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	return client.Do(req)
}
func TestWorkerHandler_All(t *testing.T) {
	server := NewServerHTTPS(":1234", "../../../localhost.crt", "../../../localhost.key", "../../../.env")

	go func() {
		err := server.RunServer()
		t.Error(err)
	}()

	time.Sleep(2 * time.Second)

	t.Run("CreateWorker", func(t *testing.T) {
		w := GenWorkerNoID()
		buff, er := json.Marshal(w)
		if er != nil {
			t.Fatalf("Failed to marshal worker: %v", er)
		}

		resp, e := Request(http.MethodPost, "/create", buff)
		if e != nil {
			t.Fatalf("Request failed: %v", e)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("GetWorkerByID", func(t *testing.T) {
		w := db_psql.CreateTestWorker()

		resp, e := Request(http.MethodGet, "/"+w.ID.String(), []byte("d"))
		if e != nil {
			t.Fatalf("Request failed: %v", e)
		}
		resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func TestHash(t *testing.T) {
	pass := "vovavova"

	for i := 0; i < 5; i++ {
		fmt.Println(db_psql.HashPass(pass))
		fmt.Println("#", i)
	}
}

func GenWorkerNoID() *types.Worker {
	return &types.Worker{
		Name:     genPassORName(6),
		Password: genPassORName(8),
		Email:    genPassORName(5) + "@mail.com",
		Role:     genRole(),
	}
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
