package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type IPaystack interface {
	InitiateTransaction(InitiateTransactionDto) (*InitiateTransactionResponse, error)
}

type paystack struct{}

func NewPaystackAPI() *paystack {
	return &paystack{}
}

func sendRequest(method string, url string, body io.Reader) ([]byte, error) {
	baseUrl := os.Getenv("PAYSTACK_BASE_URL")

	req, _ := http.NewRequest(method, baseUrl+url, body)
	req.Header.Add("Authorization", os.Getenv("PAYSTACK_SECRET_KEY"))
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

type InitiateTransactionDto struct {
	Email     string `json:"email"`
	Amount    string `json:"amount"`
	Reference string `json:"reference"`
}

type InitiateTransactionResponse struct {
	AuthorizationUrl string
	AccessCode       string
	Reference        string
}

func (p *paystack) InitiateTransaction(data InitiateTransactionDto) (*InitiateTransactionResponse, error) {
	body, _ := json.Marshal(data)
	resp, err := sendRequest(http.MethodPost, "/transaction/initialize", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var responseBody map[string]any
	err = json.Unmarshal(resp, &responseBody)
	if err != nil {
		return nil, err
	}

	responseData := responseBody["data"].(map[string]any)
	return &InitiateTransactionResponse{
		AuthorizationUrl: responseData["authorization_url"].(string),
		AccessCode:       responseData["access_code"].(string),
		Reference:        responseData["reference"].(string),
	}, nil
}
