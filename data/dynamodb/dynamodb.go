package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"rocinante-books/data"
	"rocinante-books/utils/ptr"
)

type dynamodbStrategy struct {
	db *dynamodb.DynamoDB
}

func New() (data.Strategy, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("", "", ""),
		Region:      ptr.StrRef("ap-southeast-1"),
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
		TableName: ptr.StrRef("rocinante-books"),
	})
	if err != nil {
		return err
	}
	return nil
}
