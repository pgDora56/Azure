package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/pgDora56/Azure/schemas"
)

func Insert(cfg schemas.DynamoConfig, content []schemas.IntroSchedule) error {
	table := initialize(cfg)
	for _, c := range content {
		err := table.Put(c).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func Get(cfg schemas.DynamoConfig) ([]schemas.IntroSchedule, error) {
	table := initialize(cfg)

	var iss []schemas.IntroSchedule
	err := table.Scan().All(&iss)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func initialize(dynamoConfig schemas.DynamoConfig) dynamo.Table {
	cred := credentials.NewStaticCredentials(dynamoConfig.AccessToken, dynamoConfig.Secret, "") // 最後の引数は[セッショントークン]

	sess := session.Must(session.NewSession())

	db := dynamo.New(sess, &aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})
	return db.Table("Azure")
}
