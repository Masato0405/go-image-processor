package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucket := "your-s3-bucket"
	key := "uploaded-image.jpg"

	// AWS Config作成 (v2)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to load configuration"}, err
	}

	// S3クライアント作成
	svc := s3.NewFromConfig(cfg)

	// 15分有効な事前署名付きURLを生成
	presignClient := s3.NewPresignClient(svc)
	presignParams := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	presignResult, err := presignClient.PresignPutObject(context.TODO(), presignParams, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to generate presigned URL"}, err
	}

	// URLをクライアントに返す
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: fmt.Sprintf(`{"url": "%s"}`, presignResult.URL)}, nil
}

func main() {
	lambda.Start(handler)
}
