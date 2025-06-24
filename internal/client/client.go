package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zesty-co/terraform-provider-zesty/internal/models"
)

const DefaultHostURL string = "http://localhost:9000"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(host *string, token string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
		HostURL:    DefaultHostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	c.Token = token

	return &c, nil
}

func (c *Client) Validate() error {
	url := fmt.Sprintf("%s/validate", c.HostURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("x-api-key", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) CreateAccount(payload models.Payload) (*models.Account, error) {
	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/account", c.HostURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	account := models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) DeleteAccount(payload models.Payload) error {
	rb, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/account", c.HostURL)
	req, err := http.NewRequest("DELETE", url, bytes.NewReader(rb))
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}

func (c *Client) GetAccounts() (*[]models.Account, error) {
	url := fmt.Sprintf("%s/accounts", c.HostURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	account := []models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) GetAccount(accountID string) (*models.Account, error) {
	url := fmt.Sprintf("%s/account?accountID=%s", c.HostURL, accountID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	account := models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) UpdateAccount(payload models.Payload) (*models.Account, error) {
	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/account", c.HostURL)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	account := models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
