package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (users, error) {
	result := make([]User, 0, 100_000)
	scanner := bufio.NewScanner(r)

	i := 0
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		result = append(result, user)
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat, len(u))

	for _, user := range u {
		if idx := strings.Index(user.Email, "@"); idx != -1 {
			domainPart := strings.ToLower(user.Email[idx+1:])
			if strings.Contains(domainPart, domain) {
				result[domainPart]++
			}
		}
	}
	return result, nil
}
