package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func newClient(url string, apiKey string) *client {
	return &client{url, apiKey}
}

func (c *client) createBlock(block blockRequest) error {
	url := fmt.Sprintf("%s/block", c.baseURL)

	// convert the object to json
	d, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("unable to marshal block to json: %s", err)
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(d))
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}
	r.Header.Add("X-Minecraft-ID", c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got status %d", resp.StatusCode)
	}

	return nil
}

func (c *client) deleteBlock(block blockRequest) error {
	url := fmt.Sprintf("%s/block/%d/%d/%d", c.baseURL, block.X, block.Y, block.Z)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	r.Header.Add("X-Minecraft-ID", c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got status %d", resp.StatusCode)
	}

	return nil
}

func (c *client) getBlock(x, y, z int) error {
	url := fmt.Sprintf("%s/block/%d/%d/%d", c.baseURL, x, y, z)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}

	r.Header.Add("X-Minecraft-ID", c.apiKey)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("unable to execute request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got status %d", resp.StatusCode)
	}

	return nil
}
