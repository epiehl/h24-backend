package notificator

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/epiehl93/h24-notifier/internal/adapter"
	"github.com/epiehl93/h24-notifier/internal/utils"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var sesClient *ses.SES
var cognitoClient *cognito.CognitoIdentityProvider

type Notificator interface {
	Run() error
}

type NotificatorImpl struct {
	*adapter.Registry
}

func (n NotificatorImpl) Run() error {
	// time when the cycle started
	now := time.Now()

	// Determine point in time where the last notifications were sent
	var checkTime time.Time
	cycle, err := n.Cycle.GetLastSuccessfulCycle(models.NotificationCycle)
	if err != nil {
		// Assume we never checked notifications when there is not cycle available
		// Last checkpoint is 24 hours back from now
		if errors.Is(err, gorm.ErrRecordNotFound) {
			checkTime = now.Add(-24 * time.Hour)
		} else {
			return err
		}
	} else {
		checkTime = cycle.At
	}

	wishlists, err := n.Wishlist.FindWishlistsWithAvailableInOutlet()
	if err != nil {
		return err
	}

	// Compile data into maps of maps in schema { 'userSub': { 'itemSKU': *models.Item } }
	data := make(map[string]map[string]*models.Item)
	for _, list := range wishlists {
		for _, item := range list.Items {
			if item.AvailableInOutletSince.After(checkTime) {
				if data[list.UserSub] == nil {
					data[list.UserSub] = make(map[string]*models.Item)
				}
				data[list.UserSub][fmt.Sprint(item.SKU)] = item
			}
		}
	}

	for k1, v1 := range data {
		email, err := GetUserEmail(k1)
		if err != nil {
			return err
		}
		var items []*models.Item
		for _, v2 := range v1 {
			items = append(items, v2)
		}

		content := "Hey,\n\nthe following NEW items are available in the outlet:\n\n" + FormatItems(items)
		err = SendEmail(viper.GetString("aws.ses.from"), email, "New items available in outlet", content)
		if err != nil {
			utils.Log.Error(err)
		}

	}
	if err := n.Cycle.Create(&models.Cycle{At: now, Type: models.NotificationCycle, Successful: true}); err != nil {
		return err
	}
	return nil
}

func FormatItems(items []*models.Item) string {
	content := ""
	for _, item := range items {
		content += fmt.Sprintf("SKU: %d\nName: %s\nin outlet since %s", item.SKU, item.Name, item.AvailableInOutletSince)
	}

	return content
}

func SendEmail(from string, to string, subject string, body string) error {
	if sesClient == nil {
		sesClient = utils.NewSESClient()
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{ToAddresses: []*string{
			aws.String(to),
		}},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	// send mail
	result, err := sesClient.SendEmail(input)
	if err != nil {
		return err
	}

	utils.Log.Infof("email to %s succesfully sent: messageID: %s\n", to, *result.MessageId)

	return nil
}

func GetUserEmail(sub string) (string, error) {
	if cognitoClient == nil {
		cognitoClient = utils.NewCognitoClient()
	}
	input := &cognito.ListUsersInput{
		UserPoolId: aws.String(viper.GetString("cognito.poolid")),
		Filter:     aws.String(fmt.Sprintf("sub = \"%s\"", sub)),
	}

	output, err := cognitoClient.ListUsers(input)
	if err != nil {
		return "", err
	}

	if len(output.Users) > 1 {
		return "", errors.New("too many users returned")
	}

	// fetch first user
	user := output.Users[0]
	var email string
	for _, attr := range user.Attributes {
		if *attr.Name == "email" {
			email = *attr.Value
		}
	}

	if email == "" {
		return "", errors.New("could not find email address")
	}

	return email, nil
}

func NewNotificator(r *adapter.Registry) Notificator {
	return &NotificatorImpl{r}
}
