package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		panic(err)
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ipString := string(ip)

	fmt.Println(ipString)

}
