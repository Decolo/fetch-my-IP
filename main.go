package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const EXPIRE = 60 * 60 * 24

func fetchIpInfos(conn redis.Conn, ipString string, parsedData map[string]interface{}) {
	ipConfigUrl := fmt.Sprintf("https://ipapi.co/%s/json/", ipString)

	ipConfigString := getRequest(ipConfigUrl)

	fmt.Println("current IP config:", ipConfigString)

	if strings.Contains(ipConfigString, `"error": true,`) {
		return
	}

	if parsedData == nil {
		parsedData = make(map[string]interface{})
	}

	parsedData[ipString] = ipConfigString

	jsonData, err := json.Marshal(parsedData)

	if err != nil {
		panic(err)
	}

	conn.Do("SET", "ipInfos", string(jsonData), "EX", EXPIRE)
}

func run(conn redis.Conn) {
	ipString := getRequest("https://ifconfig.me")

	if !checkIP(ipString) {
		panic("Invalid IP")
	}

	fmt.Println("current IP:" + ipString)

	ipInfoString, err := redis.String(conn.Do("GET", "ipInfos"))

	if err != nil {
		fetchIpInfos(conn, ipString, nil)

		panic(err)
	}

	if isEmptyString(ipInfoString) {
		fetchIpInfos(conn, ipString, nil)

		return
	}

	var result interface{}

	err = json.Unmarshal([]byte(ipInfoString), &result)

	if err != nil {
		fetchIpInfos(conn, ipString, nil)

		panic(err)
	}

	parsedData, ok := result.(map[string]interface{})

	if !ok {
		fetchIpInfos(conn, ipString, nil)
		return
	}

	value, exists := parsedData[ipString]

	if !exists {
		fetchIpInfos(conn, ipString, parsedData)
		return
	}

	fmt.Println("current IP config:", value)

}

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	// 在结束时关闭连接
	defer conn.Close()

	run(conn)
}
