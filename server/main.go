package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

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

var assets map[string]string // key is AssetCode, value is AssetName

func get(c *fiber.Ctx) error {
	body := c.Body()
	resp := GetResponse{}
	if assets[string(body)] == "" {
		resp.AlreadyExists = false
	} else {
		resp.AlreadyExists = true
	}
	resp.AssetName = assets[string(body)]
	resp.Success = true
	resp_json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return c.SendString(string(resp_json))
}
func set(c *fiber.Ctx) error {
	body := c.Body()
	req := SetRequest{}
	resp := SetResponse{}
	resp.Success = true
	err := json.Unmarshal(body, &req)
	if err != nil {
		resp.Success = false
	}
	assets[req.AssetCode] = req.AssetName
	resp.Success = true
	resp_json, err := json.Marshal(resp)
	if err != nil {
		resp.Success = false
	}
	return c.SendString(string(resp_json))
}

func main() {
	assets = make(map[string]string, 1)
	app := fiber.New()

	app.Post("/get", get)
	app.Post("/set", set)

	app.Listen(":8080")
}
