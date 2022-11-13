package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const authHeader = "X-API-Key"

type client struct {
	baseURL string
	apiKey  string
}

type blockRequest struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Z        int    `json:"z"`
	Material string `json:"material"`
}

type blockResponse struct {
	ID       string `json:"id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Z        int    `json:"z"`
	Material string `json:"material"`
}

type schemaRequest struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Z        int    `json:"z"`
	Rotation int    `json:"rotation"`
	Schema   string `json:"schema"`
}

func newClient(url string, apiKey string) *client {
	return &client{url, apiKey}
}

func (c *client) createBlock(block blockRequest) (*blockResponse, error) {
	url := fmt.Sprintf("%s/v1/block", c.baseURL)

	// convert the object to json
	d, err := json.Marshal(block)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal block to json: %s", err)
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(d))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err)
	}
	r.Header.Add(authHeader, c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected status 200, got status: %d, message: %s", resp.StatusCode, body)
	}

	// process the block
	blockResp := &blockResponse{}
	err = json.NewDecoder(resp.Body).Decode(blockResp)
	if err != nil {
		return nil, fmt.Errorf("unable to decode block: %s", err)
	}

	return blockResp, nil
}

func (c *client) deleteBlock(block blockRequest) error {
	url := fmt.Sprintf("%s/v1/block/%d/%d/%d", c.baseURL, block.X, block.Y, block.Z)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	r.Header.Add(authHeader, c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("expected status 200, got status: %d, message: %s", resp.StatusCode, body)
	}

	return nil
}

func (c *client) getBlock(x, y, z int) (*blockResponse, error) {
	url := fmt.Sprintf("%s/v1/block/%d/%d/%d", c.baseURL, x, y, z)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err)
	}

	r.Header.Add(authHeader, c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		// get the response body
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected status 200, got status: %d, message: %s", resp.StatusCode, body)
	}

	// process the block
	block := &blockResponse{}
	err = json.NewDecoder(resp.Body).Decode(block)
	if err != nil {
		return nil, fmt.Errorf("unable to decode block: %s", err)
	}

	return block, nil
}

func (c *client) createSchema(schema schemaRequest) (string, error) {
	url := fmt.Sprintf("%s/v1/schema/%d/%d/%d/%d", c.baseURL, schema.X, schema.Y, schema.Z, schema.Rotation)

	// read the zip file
	f, err := os.Open(schema.Schema)
	if err != nil {
		return "", fmt.Errorf("unable to open schema file: %s, err: %s", schema.Schema, err)
	}

	r, err := http.NewRequest(http.MethodPost, url, f)
	if err != nil {
		return "", fmt.Errorf("unable to create request: %s", err)
	}
	r.Header.Add(authHeader, c.apiKey)
	r.Header.Add("Content-Type", "application/zip")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", fmt.Errorf("unable to execute request: %s", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status 200, got status: %d, message: %s", resp.StatusCode, body)
	}

	return string(body), nil
}

func (c *client) undoSchema(undoID string) error {
	url := fmt.Sprintf("%s/v1/schema/undo/%s", c.baseURL, undoID)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	r.Header.Add(authHeader, c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("expected status 200, got status: %d, message: %s", resp.StatusCode, body)
	}

	return nil
}

type schemaDetailsResponse struct {
	StartX int `json:"startX"`
	StartY int `json:"startY"`
	StartZ int `json:"startZ"`
	EndX   int `json:"endX"`
	EndY   int `json:"endY"`
	EndZ   int `json:"endZ"`
}

func (c *client) getSchemaDetails(undoID string) (*schemaDetailsResponse, error) {
	url := fmt.Sprintf("%s/v1/schema/details/%s", c.baseURL, undoID)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err)
	}

	r.Header.Add(authHeader, c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	schemaResp := &schemaDetailsResponse{}
	err = json.NewDecoder(resp.Body).Decode(schemaResp)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %s", err)
	}

	return schemaResp, nil
}
