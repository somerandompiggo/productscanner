package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const get_url string = "http://localhost:8080/get"
const set_url string = "http://localhost:8080/set"

// request bodies
type SetRequest struct {
	AssetCode string `json:"assetCode"`
	AssetName string `json:"assetName"`
}

// response bodies
type GetResponse struct {
	AlreadyExists bool   `json:"alreadyExists"`
	AssetName     string `json:"assetName"`
	Success       bool   `json:"success"`
	ErrorMessage  string `json:"errorMessage"`
}
type SetResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

func clear() {
	fmt.Print("\033[H\033[2J")
}
func sleep() {
	time.Sleep(time.Millisecond * 1500)
}
func sleep_long() {
	time.Sleep(time.Second * 5)
}
func get(code string) (GetResponse, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", get_url, strings.NewReader(code))
	if err != nil {
		fmt.Println("Request creation failed with: ", err.Error())
		return GetResponse{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request execution failed with: ", err.Error())
		return GetResponse{}, err
	}
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read request body: ", err.Error())
		return GetResponse{}, err
	}
	var resp_struct GetResponse
	err = json.Unmarshal(resp_body, &resp_struct)
	if err != nil {
		fmt.Println("Failed to unmarshal response JSON: ", err.Error())
		return GetResponse{}, err
	}
	return resp_struct, nil
}
func set(code string, name string) error {
	req_struct := SetRequest{
		AssetCode: code,
		AssetName: name,
	}
	req_struct_json, err := json.Marshal(req_struct)
	if err != nil {
		fmt.Println("Could not serialize request JSON")
		return err
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", set_url, strings.NewReader(string(req_struct_json)))
	if err != nil {
		fmt.Println("Request creation failed with: ", err.Error())
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request execution failed with: ", err.Error())
		return err
	}
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read response body: ", err.Error())
		return err
	}
	var resp_struct SetResponse
	err = json.Unmarshal(resp_body, &resp_struct)
	if err != nil {
		fmt.Println("Failed to unmarshal response JSON: ", err.Error())
		return err
	}
	if !resp_struct.Success {
		fmt.Println("Server side error occurred: ", resp_struct.ErrorMessage)
		return fmt.Errorf("server side error occurred: %s", resp_struct.ErrorMessage)
	}
	return nil
}

func main() {
	var code string
	for {
		clear()
		fmt.Printf("Awaiting input...")
		fmt.Scanln(&code)
		clear()
		if code == "q" {
			fmt.Println("Exiting...")
			return
		}
		fmt.Println("Contacting server...")
		resp, err := get(code)
		if err != nil {
			fmt.Println("An error occurred whilst retrieving data. Restarting...")
			sleep_long()
			continue
		}
		if !resp.AlreadyExists {
			fmt.Printf("Asset does not yet exist in database. Please enter a name (q to cancel): ")
			var name string
			fmt.Scanln(&name)
			if name == "q" {
				fmt.Println("Cancelling...")
				sleep()
				continue
			}
			clear()
			fmt.Println("Contacting server...")
			err = set(code, name)
			if err != nil {
				fmt.Println("An error occurred whilst retrieving data. Restarting...")
				sleep_long()
				continue
			}
			fmt.Println("Set!")
		} else {
			fmt.Printf("'%s' was scanned.", resp.AssetName)
		}
		sleep()
	}
}
