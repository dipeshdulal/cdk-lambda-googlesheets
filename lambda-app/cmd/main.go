package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const SHEET_READONLY = "https://www.googleapis.com/auth/spreadsheets.readonly"

type MyData struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Data    string `json:"data"`
}

func HandleRequest(ctx context.Context, evt events.SNSEvent) error {
	refreshToken := os.Getenv("REFRESH_TOKEN")
	credentials := os.Getenv("CREDENTIALS")

	log.Println(refreshToken)
	log.Println(credentials)

	var oauthToken = oauth2.Token{
		RefreshToken: refreshToken,
	}

	config, err := google.ConfigFromJSON([]byte(credentials), SHEET_READONLY)
	if err != nil {
		log.Fatal(err)
	}

	client := config.Client(ctx, &oauthToken)

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
	}

	sheetID := "10KkEbphPRgoKfRdZNKkKT28QDt6VJhIARAWj6bi6Qxc"

	for i := range evt.Records {
		record := evt.Records[i]
		log.Println(record.SNS.Message)
		log.Println(record.SNS.MessageAttributes)

		data := MyData{}
		if err := json.Unmarshal([]byte(record.SNS.Message), &data); err != nil {
			log.Fatal(err)
		}

		writeRange := "Sheet1!F:F"
		valRange := sheets.ValueRange{
			Range: writeRange,
			Values: [][]interface{}{
				{time.Now().String()},
			},
		}
		op, err := service.Spreadsheets.Values.Append(sheetID, writeRange, &valRange).ValueInputOption("RAW").Do()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(op)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
