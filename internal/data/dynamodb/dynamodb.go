package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	config2 "github.com/rgaquino/rocinante-books/bulk/raw2json/config"
	data2 "github.com/rgaquino/rocinante-books/internal/data"
)

type strategy struct {
	db *dynamodb.DynamoDB
}

func New(c *config2.AWS) (data2.Strategy, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Region:      aws.String(c.Region),
	})
	if err != nil {
		return nil, err
	}
	return &strategy{
		db: dynamodb.New(sess),
	}, nil
}

func (s *strategy) Create(entity data2.Entity) error {
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

func (s *strategy) CreateAll(entities []data2.Entity) error {
	panic("not implemented")
}
