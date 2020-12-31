package client

import (
	"fmt"
	"log"
)

type SymTarget struct {
	Id            string   `json:"id,omitempty"`
	Type          string   `json:"type"`
	Label         string   `json:"label"`
	IntegrationId string   `json:"integration_id"`
	Settings      Settings `json:"settings"`
}

func (s SymTarget) String() string {
	return fmt.Sprintf(
		"{id=%s, type=%s, label=%s, integration_id=%s, settings=%v",
		s.Id,
		s.Type,
		s.Label,
		s.IntegrationId,
		s.Settings,
	)
}

type TargetClient interface {
	Create(target SymTarget) (string, error)
	Read(id string) (*SymTarget, error)
	Update(target SymTarget) (string, error)
	Delete(id string) (string, error)
}

func NewTargetClient(httpClient SymHttpClient) TargetClient {
	return &targetClient{
		HttpClient: httpClient,
	}
}

type targetClient struct {
	HttpClient SymHttpClient
}

func (c *targetClient) Create(target SymTarget) (string, error) {
	log.Printf("Creating Sym Target: %v", target)
	result := SymTarget{}

	if _, err := c.HttpClient.Create("/targets/", &target, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates target was not created")
	}

	log.Printf("Created Sym Target: %s", result.Id)
	return result.Id, nil
}

func (c *targetClient) Read(id string) (*SymTarget, error) {
	log.Printf("Getting Sym Target: %s", id)
	result := SymTarget{}

	if err := c.HttpClient.Read(fmt.Sprintf("/targets/%s/", id), &result); err != nil {
		return nil, err
	}

	log.Printf("Got Sym Target: %s", id)
	return &result, nil
}

func (c *targetClient) Update(target SymTarget) (string, error) {
	log.Printf("Updating Sym Target: %v", target)
	result := SymTarget{}

	if _, err := c.HttpClient.Update(fmt.Sprintf("/targets/%s/", target.Id), &target, &result); err != nil {
		return "", err
	}

	if result.Id == "" {
		return "", fmt.Errorf("response indicates Sym Target was not updated")
	}

	log.Printf("Updated Sym Target: %s", result.Id)
	return result.Id, nil
}

func (c *targetClient) Delete(id string) (string, error) {
	log.Printf("Deleting Sym Target: %s", id)

	if err := c.HttpClient.Delete(fmt.Sprintf("/targets/%s/", id)); err != nil {
		return "", err
	}

	return id, nil
}
