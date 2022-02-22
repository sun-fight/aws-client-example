package initial

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/sun-fight/aws-client/mdynamodb"
)

var _cfg aws.Config

// 初始化aws-service
func InitAwsService() (err error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	_cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-2"),
		// local dynamodb
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: "http://localhost:8000"}, nil
				}),
		),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("97rwl", "oqzqpi", "")),
	)

	if err != nil {
		return err
	}

	// init dynamodb
	mdynamodb.Init(_cfg)
	return nil
}
