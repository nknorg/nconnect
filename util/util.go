package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/nknorg/tuna"
)

func MergeStrings(src, target []string) []string {
	resSet := make(map[string]struct{}, len(src)+len(target))
	for _, s := range src {
		resSet[s] = struct{}{}
	}
	for _, s := range target {
		resSet[s] = struct{}{}
	}

	res := make([]string, 0, len(resSet))
	for s := range resSet {
		res = append(res, s)
	}

	return res
}

func RemoveStrings(src, target []string) []string {
	resSet := make(map[string]struct{}, len(src))
	for _, s := range src {
		resSet[s] = struct{}{}
	}
	for _, s := range target {
		delete(resSet, s)
	}

	res := make([]string, 0, len(resSet))
	for s := range resSet {
		res = append(res, s)
	}

	return res
}

func JSONConvert(src, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

func MatchRegex(patterns []string, s string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			log.Println("Regexp match error:", err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func ParseExecError(err error) string {
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return string(ee.Stderr)
		}
		return err.Error()
	}
	return ""
}

// IsValidUrl tests a string to determine if it is a well-structured url or not.
func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func GetRemotePrice(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	price := strings.TrimSpace(string(b))
	_, _, err = tuna.ParsePrice(price)
	if err != nil {
		return "", err
	}
	return price, nil
}
