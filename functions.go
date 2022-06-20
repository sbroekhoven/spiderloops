package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Shuffle(slice []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]string, len(slice))
	perm := r.Perm(len(slice))
	for i, randIndex := range perm {
		ret[i] = slice[randIndex]
	}
	return ret
}

func getWords(words []string, amount int) []string {
	var selected []string
	total := len(words)
	randomLine := rand.Intn(total - amount)
	for i := 1; i <= amount; i++ {
		selected = append(selected, words[randomLine])
		randomLine++
	}
	return selected
}

func findWord(slice []string, word string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[word]
	return ok
}

func (conf Config) ValidateUser(username, password string, c echo.Context) (bool, error) {

	auth_user := conf.AuthUsername
	auth_password := conf.AuthPassword
	real_ip := c.RealIP()
	user_agent := c.Request().UserAgent()
	request_url := c.Request().URL.String()

	auth_check := bcrypt.CompareHashAndPassword([]byte(auth_password), []byte(password))
	if auth_check != nil {
		// Passwords not equal
		logrus.WithFields(logrus.Fields{
			"authentication": "failed",
			"real_ip":        real_ip,
			"user_agent":     user_agent,
			"request_url":    request_url,
		}).Warn("Authentication was not successful")
	} else {
		// Login success!
		logrus.WithFields(logrus.Fields{
			"authentication": "success",
			"real_ip":        real_ip,
			"user_agent":     user_agent,
			"request_url":    request_url,
		}).Warn("Authentication was successful")
	}

	// org
	if auth_user == username && auth_check == nil {
		return true, nil
	}
	return false, nil
}

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Apache/2.4.7 (Ubuntu)")
		return next(c)
	}
}
