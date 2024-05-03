package rest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

var netClient *http.Client

func CreateHttpClient() *http.Client {
	if netClient == nil {
		fmt.Println("client is empty")
		client := &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        30,
				MaxIdleConnsPerHost: 30,
				MaxConnsPerHost:     15,
				DisableKeepAlives:   false,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Duration(30 * time.Second),
		}
		netClient = client
	}
	fmt.Println("client is not empty")
	return netClient
}
