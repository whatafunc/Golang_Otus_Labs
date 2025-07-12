package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go" //nolint:depguard // allow json-iterator temporarily
)

// "github.com/minio/simdjson-go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	domain = strings.ToLower(domain)

	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		email := strings.ToLower(user.Email)
		at := strings.IndexByte(email, '@')
		if at < 0 || at+1 >= len(email) {
			continue
		}
		userDomain := email[at+1:]

		if userDomain == domain || strings.HasSuffix(userDomain, "."+domain) {
			result[userDomain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil // map of domain/count of matches.
}
