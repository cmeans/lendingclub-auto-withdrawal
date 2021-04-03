package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type availableFundsResponse struct {
	InvestorID    int32   `json:"investorId"`
	AvailableCash float32 `json:"availableCash"`
}

type withdrawalRequest struct {
	Amount float32 `json:"amount"`
}

type withdrawalResponse struct {
	InvestorID                 int32   `json:"investorId"`
	Amount                     float32 `json:"amount"`
	EstimatedFundsTransferDate string  `json:"estimatedFundsTransferDate"`
}

var (
	investorId        string  = os.Getenv("INVESTOR_ID")
	lendingClubAPIKey string  = os.Getenv("LENDING_CLUB_API_KEY")
	minimumAmount     float32 = 10.0
)

func init() {
	value, ok := os.LookupEnv("MINIMUM_AMOUNT")
	if ok {
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			fmt.Printf("error converting environment variable MINIMUM_AMOUNT: %s\n", err)
		} else {
			minimumAmount = float32(f)
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		investorId = os.Args[1]
		if len(os.Args) > 2 {
			lendingClubAPIKey = os.Args[2]
		}

		Handler(nil)
	} else {
		lambda.Start(Handler)
	}
}

func Handler(_ context.Context) {
	err := validateSettings()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	client := &http.Client{}

	availableFunds, err := availableFunds(client)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unable to unmarshal response") {
			fmt.Println("environment variable LENDING_CLUB_API_KEY/param 2 appears to be invalid/expired/revoked")
			os.Exit(2)
		} else {
			fmt.Println(err.Error())
			os.Exit(4)
		}
	}

	fmt.Printf("Available Cash: $%.2f\n", availableFunds.AvailableCash)

	if investorId != fmt.Sprintf("%d", availableFunds.InvestorID) {
		fmt.Println("INVESTOR_ID does not match returned InvestorID, invalid or mismatched with API key")
		os.Exit(3)
	}

	if availableFunds.AvailableCash > minimumAmount {
		request := withdrawalRequest{Amount: availableFunds.AvailableCash}
		response, err := withdraw(client, request)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(5)
		}
		fmt.Printf("Transferring: $%.2f\nEstimated Transfer Date: %s\n", response.Amount, response.EstimatedFundsTransferDate)
	}
}

func availableFunds(client *http.Client) (availableFunds availableFundsResponse, err error) {
	url := fmt.Sprintf("https://api.lendingclub.com/api/investor/%s/accounts/%s/availablecash", "v1", investorId)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return availableFunds, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", lendingClubAPIKey)

	res, err := client.Do(req)
	if err != nil {
		return availableFunds, fmt.Errorf("unable execute request: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return availableFunds, fmt.Errorf("unable to read response body: %w", err)
	}

	err = json.Unmarshal(body, &availableFunds)
	if err != nil {
		return availableFunds, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return
}

func withdraw(client *http.Client, withdraw withdrawalRequest) (response withdrawalResponse, err error) {
	url := fmt.Sprintf("https://api.lendingclub.com/api/investor/%s/accounts/%s/funds/withdraw", "v1", investorId)

	data, _ := json.Marshal(withdraw)
	buffer := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		return response, fmt.Errorf("unable to create withdrawal request: %w", err)
	}

	req.Header.Set("Authorization", lendingClubAPIKey)

	res, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("unable execute withdrawal request: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("unable to read withdrawal response body: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("unable to unmarshal withdrawal response: %w", err)
	}

	return
}

func validateSettings() error {
	if investorId == "" {
		return errors.New("environment variable INVESTOR_ID/param 1 cannot be blank")
	}

	if lendingClubAPIKey == "" {
		return errors.New("environment variable LENDING_CLUB_API_KEY/param 2 cannot be blank")
	}

	if _, err := strconv.Atoi(investorId); err != nil {
		return errors.New("environment variable INVESTOR_ID/param 1 does not appears to be valid")
	}

	if _, err := base64.StdEncoding.DecodeString(lendingClubAPIKey); err != nil {
		return errors.New("environment variable LENDING_CLUB_API_KEY/param 2 does not appear to be a valid key (is not a valid Base64 encoded string)")
	}

	return nil
}
