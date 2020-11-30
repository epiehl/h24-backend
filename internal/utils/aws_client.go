package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spf13/viper"
)

func NewCognitoClient() *cognito.CognitoIdentityProvider {
	svc := cognito.New(NewSession())
	return svc
}

func NewSESClient() *ses.SES {
	return ses.New(NewSession())
}

func NewSession() *session.Session {
	return session.Must(session.NewSession(aws.NewConfig().WithRegion(viper.GetString("aws.region"))))
}
