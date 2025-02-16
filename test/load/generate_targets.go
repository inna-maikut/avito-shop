package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

const (
	numWorkers      = 20
	totalUsers      = 10_000 // task requires 100_000, but we don't want to create all each run of generate_targets
	workerUsers     = totalUsers / numWorkers
	logProgressEach = 1_000
	targetUsers     = 10_000
)

//nolint:gocognit
func main() {
	f50, err := os.OpenFile("./test/load/targets_50-50", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		panic(fmt.Errorf("os.OpenFile: %w", err))
	}
	defer func() { _ = f50.Close() }()
	f10, err := os.OpenFile("./test/load/targets_10-90", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		panic(fmt.Errorf("os.OpenFile: %w", err))
	}
	defer func() { _ = f10.Close() }()
	f2, err := os.OpenFile("./test/load/targets_2-98", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		panic(fmt.Errorf("os.OpenFile: %w", err))
	}
	defer func() { _ = f10.Close() }()

	_, err = createUser(-1) // zero user to send him all money
	if err != nil {
		panic(fmt.Errorf("createUser: %w", err))
	}

	var counter atomic.Int32

	userTokens := make([]string, totalUsers) // create 100k users
	eg := errgroup.Group{}
	for j := range numWorkers {
		eg.Go(func() (err error) {
			var token string
			for i := j * workerUsers; i < (j+1)*workerUsers; i++ {
				token, err = createUser(i)
				userTokens[i] = token
				if err != nil {
					return fmt.Errorf("createUser %d: %w", i, err)
				}

				counter.Add(1)
				count := int(counter.Load())
				if count%logProgressEach == 0 {
					fmt.Printf("progress creating users %d from %d\n", count, totalUsers)
				}
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		panic(fmt.Errorf("eg.Wait: %w", err))
	}

	for i := range targetUsers {
		// generate only 200k requests
		token := userTokens[i]

		s := "POST http://localhost:8080/api/sendCoin\n" +
			"Authorization: " + token + "\n" +
			"Content-Type: application/json\n" +
			"@./test/load/send_coin_to_username_0.json\n\n" +
			"GET http://localhost:8080/api/buy/umbrella\n" +
			"Authorization: " + token + "\n\n"

		for range 18 {
			s += "GET http://localhost:8080/api/info\n" +
				"Authorization: " + token + "\n\n"
		}

		_, err = f10.WriteString(s)
		if err != nil {
			panic(fmt.Errorf("f.WriteString: %w", err))
		}

		s = "POST http://localhost:8080/api/sendCoin\n" +
			"Authorization: " + token + "\n" +
			"Content-Type: application/json\n" +
			"@./test/load/send_coin_to_username_0.json\n\n" +
			"GET http://localhost:8080/api/buy/socks\n" +
			"Authorization: " + token + "\n\n"

		for range 2 {
			s += "GET http://localhost:8080/api/info\n" +
				"Authorization: " + token + "\n\n"
		}

		_, err = f50.WriteString(s)
		if err != nil {
			panic(fmt.Errorf("f.WriteString: %w", err))
		}

		if i > targetUsers/10 {
			continue
		}

		s = "POST http://localhost:8080/api/sendCoin\n" +
			"Authorization: " + token + "\n" +
			"Content-Type: application/json\n" +
			"@./test/load/send_coin_to_username_0.json\n\n" +
			"GET http://localhost:8080/api/buy/socks\n" +
			"Authorization: " + token + "\n\n"

		for range 98 {
			s += "GET http://localhost:8080/api/info\n" +
				"Authorization: " + token + "\n\n"
		}

		_, err = f2.WriteString(s)
		if err != nil {
			panic(fmt.Errorf("f.WriteString: %w", err))
		}
	}
}

func createUser(i int) (string, error) {
	username := fmt.Sprintf("my-username-%d", i+1)

	body := `{"username":"` + username + `","password":"password"}`
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/auth", strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code not 200")
	}

	var out struct {
		Token string `json:"token"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.Unmarshal(bodyBytes, &out)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return out.Token, nil
}
