package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type ip_res struct {
	city           string
	region         string
	country_code   string
	continent_code string
	organisation   string
}

func doRequest(ip string) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.ipdata.co/%s", ip), nil)
	if err != nil {
		return fmt.Sprintf("Error creating request for IP %s: %s", ip, err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error retrieving IP %s: %s", ip, err)
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response body %s: %s", ip, err)
	}

	var jsonResp map[string]interface{}

	err = json.Unmarshal(resp_body, &jsonResp)
	if err != nil {
		// Use the same error text as the web service (api.ipdata.co).
		return fmt.Sprintf("%s is a private IP address", ip)
	} else {
		return fmt.Sprintf("%s, %s %s %s %s, %s", ip, jsonResp["city"].(string), jsonResp["region"].(string), jsonResp["country_code"].(string), jsonResp["continent_code"].(string), jsonResp["organisation"].(string))
	}
}

func usage() {
	fmt.Print(`Usage: locate-ip ...IP addresses
Examples:
    locate-ip 192.168.1.92 127.0.0.1
    traceroute example.com | awk '{print $2}' | xargs locate-ip
`)
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	// Match an IP address.
	// Note that this doesn't filter out illegitimate IPs, but that's not the goal (only to match).
	r, err := regexp.Compile("^([0-9]{1,3}\\.){3}([0-9]{1,3}).*?")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for n, ip := range os.Args[1:] {
		// For now, discard any entries that contain non-numeric chars (i.e., "*").
		// This is intended to filter entries if the `locate-ip` is called like:
		//
		//      traceroute example.com | awk '{print $2}' | xargs ./locate_ip
		//
		if r.MatchString(ip) {
			fmt.Printf("Hop %d %s\n", n, doRequest(r.FindString(ip)))
		}
	}
}
