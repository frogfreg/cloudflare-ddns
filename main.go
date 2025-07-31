package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	if err := loadConfig(); err != nil {
		panic(err)
	}

	prevIp, err := os.ReadFile("./last-ip.txt")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ipString := strings.TrimSpace(string(ip))

	if len(ipString) < 1 {
		panic("ip is empty, something went wrong")
	}

	if ipString != strings.TrimSpace(string(prevIp)) {
		fmt.Println("IPs are different")
		panic("stopping before calling cloudflare")
		updateIp(ipString)
	}

}

func loadConfig() error {
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetDefault("zone-id", "")
			viper.SetDefault("dns-record-id", "")
			viper.SetDefault("cloudflare-email", "")
			viper.SetDefault("cloudflare-api-key", "")
			viper.SetDefault("domain-name", "")
			viper.SetDefault("ttl", 1)
			viper.SetDefault("type", "A")
			viper.SetDefault("ipv4", "")

			if writeErr := viper.WriteConfigAs("./config.toml"); writeErr != nil {
				fmt.Println("entering here")
				return writeErr
			}
			return fmt.Errorf("no config found, creating config with default values")

		} else {
			return err
		}
	}

	//todo: validate valid values in config

	return nil

}

type ARecordInfo struct {
	Domain string `json:"name"`
	Ttl    int    `json:"ttl"`
	Type   string `json:"type"`
	Ipv4   string `json:"content"`
}

type CloudflareDnsResponse struct {
	Success bool
	Result  any
}

func updateIp(newIp string) error {

	recordInfo := ARecordInfo{
		Domain: viper.GetString("domain-name"),
		Ttl:    viper.GetInt("ttl"),
		Type:   viper.GetString("Type"),
		Ipv4:   newIp,
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records/%v", viper.GetString("zone-id"), viper.GetString("dns-record-id"))

	body, err := json.Marshal(recordInfo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Auth-Email", viper.GetString("cloudflare-email"))
	req.Header.Add("X-Auth-Key", viper.GetString("cloudflare-api-key"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code other than 200 - status code is: %v-%v, body of request is: %v", resp.StatusCode, resp.Status, string(respBody))
	}

	var cdnsr CloudflareDnsResponse

	if err := json.Unmarshal(respBody, &cdnsr); err != nil {
		return err
	}

	if !cdnsr.Success {
		return fmt.Errorf("no success on request. %v", cdnsr.Result)
	}

	return writeIpToFile(newIp)
}

func writeIpToFile(ip string) error {
	if err := os.WriteFile("./last-ip.txt", []byte(ip), 0o666); err != nil {
		return err
	}
	return nil
}
