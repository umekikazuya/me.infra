package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type responseBody struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func handler(_ context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	body := responseBody{
		Message: "placeholder api deployed; app repo pipeline should replace this code",
		Status:  "placeholder",
	}
	statusCode := http.StatusServiceUnavailable

	if request.RequestContext.HTTP.Method == http.MethodGet && request.RawPath == "/up" {
		body.Status = "ok"
		statusCode = http.StatusOK
	}

	response, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"content-type":  "application/json",
			"cache-control": "no-store",
		},
		Body: string(response),
		Cookies: []string{
			"placeholder=1; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=300",
		},
	}, nil
}

func main() {
	lambda.StartWithOptions(handler, lambda.WithEnableSIGTERM(func() {
		time.Sleep(10 * time.Millisecond)
	}))
}
