package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"rocinante-books/config"
	"rocinante-books/data"
)

type dynamodbStrategy struct {
	db *dynamodb.DynamoDB
}

func New(c *config.AWS) (data.Strategy, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Region:      aws.String(c.Region),
	})
	if err != nil {
		return nil, err
	}
	return &dynamodbStrategy{
		db: dynamodb.New(sess),
	}, nil
}

func (s *dynamodbStrategy) Create(entity data.Entity) error {
	av, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		return err
	}
	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(entity.TableName()),
	})
	if err != nil {
		return err
	}
	return nil
}
