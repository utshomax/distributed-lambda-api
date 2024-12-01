package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"api-lens/pkg/config"
	"api-lens/pkg/metrics"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse configuration
	cfg, err := config.Parse([]byte(request.Body))
	log.Println(cfg)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	// Collect metrics
	metricsCollection := metrics.CollectMetrics(ctx, cfg)
	log.Println(metricsCollection)

	// Convert to JSON response
	response, err := json.Marshal(metricsCollection)
	log.Println(string(response))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(handler)
}
