package main

import (
	"io"
	"net"
	"net/http"
	"strings"
)

func checkIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func getRequest(url string) string {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	bodyString := string(bodyBytes)

	return bodyString
}

func isEmptyString(val string) bool {
	return strings.TrimSpace(val) == ""
}
