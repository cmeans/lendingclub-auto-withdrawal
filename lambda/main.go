package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

type availableFundsResponse struct {
	InvestorID    int32   `json:"investorId"`
	AvailableCash float32 `json:"availableCash"`
}

type withdrawlRequest struct {
	Amount float32 `json:"amount"`
}

type withdrawlResponse struct {
	InvestorID                 int32   `json:"investorId"`
	Amount                     float32 `json:"amount"`
	EstimatedFundsTransferDate string  `json:"estimatedFundsTransferDate"`
}

var (
	investorId        string = os.Getenv("INVESTOR_ID")
	lendingClubAPIKey string = os.Getenv("LENDING_CLUB_API_KEY")
	minimumAmount     float32 = 10.0
)

func init() {
	value, ok := os.LookupEnv("MINIMUM_AMOUNT")
	if ok {
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			fmt.Printf("error converting env var MINIMUM_AMOUNT: %s\n", err)
		} else {
			minimumAmount = float32(f)
		}
	}
}

func Handler(_ context.Context) {
	client := &http.Client{}

	availableFunds, err := availableFunds(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", availableFunds)
	if availableFunds.AvailableCash > minimumAmount {
		request := withdrawlRequest{Amount: availableFunds.AvailableCash}
		response, err := withdraw(client, request)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", response)
	}
}

func availableFunds(client *http.Client) (availableFunds availableFundsResponse, err error) {
	url := fmt.Sprintf("https://api.lendingclub.com/api/investor/%s/accounts/%s/availablecash", "v1", investorId)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", lendingClubAPIKey)

	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &availableFunds)
	if err != nil {
		return
	}

	return
}

func withdraw(client *http.Client, withdraw withdrawlRequest) (response withdrawlResponse, err error) {
	url := fmt.Sprintf("https://api.lendingclub.com/api/investor/%s/accounts/%s/funds/withdraw", "v1", investorId)

	data, _ := json.Marshal(withdraw)
	buffer := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		return response, err
	}

	req.Header.Set("Authorization", lendingClubAPIKey)

	res, err := client.Do(req)
	if err != nil {
		return response, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
