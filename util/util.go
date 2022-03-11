package util

import (
	"encoding/json"
	"log"
	"net"
	"os/exec"
	"regexp"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

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
