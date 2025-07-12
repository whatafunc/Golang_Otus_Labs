//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat_InvalidJSON(t *testing.T) {
	data := `{"Id":1,"Email":"user@example.com"}
			{"Id":2,"Email":"user2@example.com"}
			invalid json line
			{"Id":3,"Email":"user3@example.com"}`

	_, err := GetDomainStat(bytes.NewBufferString(data), "example.com")
	require.Error(t, err)
}

func TestGetDomainStat_MissingOrMalformedEmail(t *testing.T) {
	data := `{"Id":1,"Email":"valid@bam.com"}
{"Id":2,"Email":"invalid-email-without-at"}
{"Id":3,"Email":""}
{"Id":4}`

	result, err := GetDomainStat(bytes.NewBufferString(data), "bam.com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{"bam.com": 1}, result)
}

func TestGetDomainStat_CaseInsensitive(t *testing.T) {
	data := `{"Id":1,"Email":"User@Example.COM"}
{"Id":2,"Email":"another@EXample.com"}`

	result, err := GetDomainStat(bytes.NewBufferString(data), "eXample.com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{"example.com": 2}, result)
}

func TestGetDomainStat_EmptyInput(t *testing.T) {
	result, err := GetDomainStat(bytes.NewBufferString(""), "com")
	require.NoError(t, err)
	require.Empty(t, result)
}

func TestGetDomainStat_LargeInput(t *testing.T) {
	var b strings.Builder
	for i := 0; i < 10000; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		line := fmt.Sprintf(`{"Id":%d,"Email":"%s"}`+"\n", i, email)
		b.WriteString(line)
	}

	result, err := GetDomainStat(strings.NewReader(b.String()), "example.com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{"example.com": 10000}, result)
}

func TestGetDomainStat_DomainSuffixMatching(t *testing.T) {
	data := `{"ID":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"user@notexample.com","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"ID":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"user@example.com.au","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"ID":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"user@example.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}`

	result, err := GetDomainStat(bytes.NewBufferString(data), "example.com")
	require.NoError(t, err)
	require.Equal(t, DomainStat{"example.com": 1}, result)
}
