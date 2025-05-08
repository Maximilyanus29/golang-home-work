package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = strings.ToLower(domain)
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		email := extractEmail(scanner.Bytes())
		if email == "" {
			continue
		}

		domainPart := getDomainPart(email)
		if strings.HasSuffix(domainPart, "."+domain) {
			result[domainPart]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func extractEmail(line []byte) string {
	key := []byte(`"Email":"`)
	start := bytes.Index(line, key)
	if start == -1 {
		return ""
	}

	start += len(key)
	end := bytes.IndexByte(line[start:], '"')
	if end == -1 {
		return ""
	}

	return string(line[start : start+end])
}

func getDomainPart(email string) string {
	at := strings.LastIndexByte(email, '@')
	if at == -1 {
		return ""
	}
	return strings.ToLower(email[at+1:])
}
