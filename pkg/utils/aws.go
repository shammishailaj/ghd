package utils

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetAWSSession - returns an active AWS session for further use with the components of the AWS SDK
//  and error in case it is not able to create one.
func (u *Utils) GetAWSSession() (*session.Session, error) {
	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsProfile := os.Getenv("AWS_PROFILE")
	awsToken := ""
	t := strings.ToUpper(os.Getenv("AWS_SESSION_DEBUG"))
	debug := false
	if t == "TRUE" {
		debug = true
	}

	if debug {
		log.Printf("Initiating AWS Seesion with AWS_PROFILE = %s, AWS_REGION = %s, AWS_ACCESS_KEY_ID = %s, AWS_SECRET_ACCESS_KEY = %s", awsProfile, awsRegion, awsAccessKeyID, awsSecretAccessKey)
	} else {
		log.Printf("Initiating AWS Seesion with AWS_PROFILE = %s, AWS_REGION = %s", awsProfile, awsRegion)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsToken),
	})

	return sess, err
}

// ExitErrorf - A function to output error and kill the current process returning the control back to the OS.
func (u *Utils) ExitErrorf(msg string, args ...interface{}) {
	log.Printf(msg+"\n", args...)
	os.Exit(1)
}

