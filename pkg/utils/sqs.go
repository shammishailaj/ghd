package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/shammishailaj/ghd/pkg/schemas"
	"strconv"
	"strings"
	"time"
)

// InsertEnqueueRetry - A function to insert data into the redrive queue
func (u *Utils) InsertEnqueueRetry(query string, qURL string, dbProd *RDBMS) {
	t := time.Now()

	insertArray := make(map[string]string)

	insertArray["request"] = string(query)
	insertArray["queue_url"] = qURL
	insertArray["status"] = "0"
	insertArray["created_at"] = t.Format("2006-01-02 15:04:05")

	insertRedrive := dbProd.InsertRow("enqueue_retry", insertArray)

	log.Println("insertRedrive ", insertRedrive)
}

// SendCurl - A function to send a HTTP POST request via golang
func (u *Utils) SendCurl(query string, url string) (result *http.Response) {
	payload := strings.NewReader(string(query))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	return res
}

// EnqueueSQS - Add data to queue
func(u *Utils) EnqueueSQS(arrayDataJSON string, qURL string, dbM2 *RDBMS, dbM2Err error) (int, int) {
	qSuccess := 0
	qFailure := qSuccess
	postData := make(map[string]string)

	postData["data"] = string(arrayDataJSON)
	postData["qURL"] = qURL
	postData["operationType"] = "enqueue"
	query, _ := json.Marshal(postData)
	microServiceURL := os.Getenv("MQ_SERVICE_URL") //"https://em2hwb8dyl.execute-api.ap-southeast-1.amazonaws.com/development/mq"
	queueName := os.Getenv("QUEUE_NAME") //"https://em2hwb8dyl.execute-api.ap-southeast-1.amazonaws.com/development/mq"

	result := u.SendCurl(string(query), microServiceURL)

	defer result.Body.Close()

	mqResponse, mqResponseErr := ioutil.ReadAll(result.Body)
	mqResponseErrText := ""

	if mqResponseErr != nil {
		mqResponseErrText = mqResponseErr.Error()
	}

	var mqResponseList schemas.MqResponseBody

	mqResponseListErr := json.Unmarshal(mqResponse, &mqResponseList)

	if mqResponseListErr != nil {
		log.Errorf("FAILED to parse JSON response from MQ Service")
		return qSuccess, qFailure
	}

	if mqResponseList.Status == "true" {
		qSuccess++
	} else {
		qFailure++
		emailContent := fmt.Sprintf("<p>Inserting into %s Queue FAILED</p><p>%s</p><p>%s</p>", queueName, string(mqResponse), mqResponseErrText)
		// InsertEnqueueRetry(string(query), qURL, dbProd)
		if dbM2Err == nil {
			u.InsertEmailPool("shammi.shailaj@healthians.com", "BigQuery " + queueName + " Queue insert failure", emailContent, dbM2)
		} else {
			log.Printf("Unable to insert \"Queue insert failure\" email into M2 database. %s", dbM2Err.Error())
		}
	}
	return qSuccess, qFailure
}

// EnqueueSQSBatch - Add data to queue in batches of 10 messages
func (u *Utils) EnqueueSQSBatch(arrayDataJSON string, qURL string, dbM2 *RDBMS, dbM2Err error) (int, int) {
	qSuccess := 0
	qFailure := qSuccess
	postData := make(map[string]string)

	postData["data"] = string(arrayDataJSON)
	postData["qURL"] = qURL
	postData["operationType"] = "enqueue"
	query, _ := json.Marshal(postData)
	microServiceURL := os.Getenv("MQ_SERVICE_URL") //"https://em2hwb8dyl.execute-api.ap-southeast-1.amazonaws.com/development/mq"
	queueName := os.Getenv("QUEUE_NAME") //"https://em2hwb8dyl.execute-api.ap-southeast-1.amazonaws.com/development/mq"

	result := u.SendCurl(string(query), microServiceURL)

	defer result.Body.Close()

	mqResponse, mqResponseErr := ioutil.ReadAll(result.Body)

	var mqResponseList schemas.MqResponseBody

	json.Unmarshal(mqResponse, &mqResponseList)

	if mqResponseList.Status == "true" {
		// log.Println("inserted in queue")
		qSuccess++
	} else {
		qFailure++
		emailContent := fmt.Sprintf("<p>Inserting into %s Queue FAILED</p><p>%s</p><p>%s</p>", queueName, string(mqResponse), mqResponseErr.Error())
		// InsertEnqueueRetry(string(query), qURL, dbProd)
		if dbM2Err == nil {
			u.InsertEmailPool("shammi.shailaj@healthians.com", "BigQuery " + queueName + " Queue insert failure", emailContent, dbM2)
		} else {
			log.Printf("Unable to insert \"Queue insert failure\" email into M2 database. %s", dbM2Err.Error())
		}
	}
	return qSuccess, qFailure

}

// EnqueueSQSImport - Add Slice data to queue
func (u *Utils) EnqueueSQSImport(qMsgs []schemas.QueueMessageDatum, qURL string, dbProd *sql.DB, dbM2 *RDBMS, dbM2Err error) {
	qMsgsLen := len(qMsgs)
	log.Printf("qMsgsLen = %d", qMsgsLen)
	qSuccess := 0
	qFailure := qSuccess

	for i := 0; i < qMsgsLen; i++ {
		arrayDataJSON, err := json.Marshal(qMsgs[i])
		if err != nil {
			log.Printf("Error marshalling struct %#v into JSON", qMsgs[i])
			log.Printf("Error details = %#v", err)
		} else {
			qSuccess, qFailure = u.EnqueueSQS(string(arrayDataJSON), qURL, dbM2, dbM2Err)
		}
	}
	log.Printf("Messages in queue - Successful %d Unsuccessful %d", qSuccess, qFailure)
}

func (u *Utils) ReceiveMessageSQS(sess *session.Session, maxMsgs, qURL string) (*sqs.SQS, *sqs.ReceiveMessageOutput, error) {
	envMaxMsg := maxMsgs
	maxMsg, errMaxMsg := strconv.ParseInt(envMaxMsg, 10, 64)

	if errMaxMsg != nil {
		log.Errorf("Error errMaxMsg for receive message from SQS. %s", errMaxMsg)
		log.Infof("Using a default value of 10")
		maxMsg = 10
	}

	svc := sqs.New(sess)

	log.Infof("Receiving Messages from Queue...")

	sqsRecv, sqsRecvErr := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(maxMsg),
		VisibilityTimeout:   aws.Int64(60), // 20 seconds
		WaitTimeSeconds:     aws.Int64(15),
	})

	return svc, sqsRecv, sqsRecvErr
}

// DeleteMessage - Delete a message from queue
func (u *Utils) DeleteMessage(qURL string, svc *sqs.SQS, receiptHandler *string) error {
	resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: receiptHandler,
	})

	if err != nil {
		log.Errorf("Error deleting message from queue. Details: %s", err.Error())
		return err
	}

	log.Infof("Message deleted successfully. Details: %s", resultDelete.String())
	return nil
}

// Output - Function to write data to HTTP Writer stream
func Output(w http.ResponseWriter, jsonData []byte, httpStatus int) {
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, string(jsonData))
	return
}

