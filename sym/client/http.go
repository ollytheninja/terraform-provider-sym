package client

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/symopsio/terraform-provider-sym/sym/utils"
)

func NewSymHttpClient(apiUrl, token string) SymHttpClient {
	return &symHttpClient{
		apiUrl: apiUrl,
		jwt:    token,
	}
}

type SymHttpClient interface {
	Do(method, path string, payload interface{}) (string, error)
	Create(path string, payload interface{}, result interface{}) (string, error)
	Read(path string, result interface{}) error
	Update(path string, payload interface{}, result interface{}) (string, error)
	Delete(path string) error
}

type symHttpClient struct {
	apiUrl string
	jwt    string
}

func (c *symHttpClient) getUrl(path string) string {
	base := strings.TrimRight(c.apiUrl, "/")
	return base + "/" + strings.TrimLeft(path, "/")
}

func (c *symHttpClient) Do(method string, path string, payload interface{}) (string, error) {
	url := c.getUrl(path)
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	log.Printf("submitting request: %s %s %s", method, path, string(b))
	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	requestID := uuid.New().String()
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sym-Request-ID", requestID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// no status code if there was an error at this point
		return "", utils.ErrAPIConnect(path, requestID)
	}

	body, err := io.ReadAll(resp.Body)

	// Return specific errors based on the status code from the Sym API.
	if resp.StatusCode == 400 {
		errorBody := utils.ErrorResponse{}
		err = json.Unmarshal(body, &errorBody)
		if err != nil {
			return "", err
		}
		return "", utils.ErrAPIBadRequest(errorBody.Errors)
	} else if resp.StatusCode == 401 {
		return "", utils.ErrConfigFileNoJWT
	} else if resp.StatusCode == 403 {
		return "", utils.ErrUserIsNotAdmin
	} else if resp.StatusCode == 404 {
		return "", utils.ErrAPINotFound(path, requestID)
	}

	if err != nil {
		// If we don't have a specific error for the status code, and we got an error message in the response,
		// display the full error message to the user.
		return "", err
	} else if resp.StatusCode > 400 {
		// We weren't able to get an error message from the API, and we don't have a specific error message
		// for this status code, but we know it failed. Display a generic error message to the user.
		return "", utils.ErrAPIUnexpected(path, requestID, resp.StatusCode)
	}

	return string(body), nil
}

func (c *symHttpClient) Create(path string, payload interface{}, result interface{}) (string, error) {
	body, err := c.Do("POST", path, payload)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(body), result); err != nil {
		return "", err
	}

	log.Printf("got response: %v", result)
	return body, nil
}

func (c *symHttpClient) Read(path string, result interface{}) error {
	body, err := c.Do("GET", path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(body), result)
}

func (c *symHttpClient) Update(path string, payload interface{}, result interface{}) (string, error) {
	body, err := c.Do("PATCH", path, payload)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(body), result); err != nil {
		return "", err
	}

	log.Printf("got response: %v", result)
	return body, nil
}

func (c *symHttpClient) Delete(path string) error {
	if _, err := c.Do("DELETE", path, nil); err != nil {
		return err
	}

	return nil
}
