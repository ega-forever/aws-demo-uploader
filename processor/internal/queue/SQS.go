package queue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	AWSSQS "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ega-forever/aws-demo-uploader/internal/models"
	"time"
)

type SQS struct {
	sqs             *AWSSQS.SQS
	uri             string
	unsubscribeCh   chan bool
	eventCh         chan *models.QueueMessage
	errorCh         chan error
	maxMessages     int
	currentMessages map[string]*models.QueueMessage
	lockTimeSeconds int64
}

type SQSMessageBodyRecordsObject struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}
type SQSMessageBodyRecordsS3 struct {
	Object SQSMessageBodyRecordsObject `json:"object"`
}

type SQSMessageBodyRecords struct {
	EventSource string                  `json:"eventSource"`
	AwsRegion   string                  `json:"awsRegion"`
	S3          SQSMessageBodyRecordsS3 `json:"s3"`
}

type SQSMessage struct {
	Id      string
	Records []SQSMessageBodyRecords `json:"records"`
}

func New(uri string, region string, maxMessages int, lockTimeSeconds int64) *SQS {

	creds := credentials.NewEnvCredentials()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		},
	}))

	sqs := AWSSQS.New(sess)

	return &SQS{
		sqs:             sqs,
		uri:             uri,
		unsubscribeCh:   make(chan bool),
		eventCh:         make(chan *models.QueueMessage),
		errorCh:         make(chan error),
		maxMessages:     maxMessages,
		currentMessages: make(map[string]*models.QueueMessage),
		lockTimeSeconds: lockTimeSeconds,
	}
}

func (sqs *SQS) Subscribe() (chan *models.QueueMessage, chan error) {

	go func() {

		interval := time.NewTicker(time.Second)

		for {
			select {
			case <-interval.C:
				{

					if len(sqs.currentMessages) >= sqs.maxMessages {
						continue
					}

					message, err := sqs.getMessage()

					// todo
					if err != nil {
						sqs.errorCh <- err
						return
					}

					if message != nil {

						sqs.currentMessages[message.Id] = message
						sqs.eventCh <- message
					}

				}
			case <-sqs.unsubscribeCh:
				return
			}

		}

	}()

	return sqs.eventCh, sqs.errorCh

}

func (sqs *SQS) Unsubscribe() {
	sqs.unsubscribeCh <- true
	sqs.currentMessages = make(map[string]*models.QueueMessage)
}

func (sqs *SQS) AckMessage(messageId string) error {

	_, err := sqs.sqs.DeleteMessage(&AWSSQS.DeleteMessageInput{
		QueueUrl:      &sqs.uri,
		ReceiptHandle: &messageId,
	})

	if err != nil {
		return err
	}

	delete(sqs.currentMessages, messageId)

	return nil
}

func (sqs *SQS) getMessage() (*models.QueueMessage, error) {
	result, err := sqs.sqs.ReceiveMessage(&AWSSQS.ReceiveMessageInput{
		QueueUrl: &sqs.uri,
		AttributeNames: aws.StringSlice([]string{
			"SentTimestamp",
		}),
		MaxNumberOfMessages: aws.Int64(int64(sqs.maxMessages)),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		WaitTimeSeconds:   aws.Int64(10),
		VisibilityTimeout: aws.Int64(sqs.lockTimeSeconds),
	})
	if err != nil {
		return nil, err
	}

	if len(result.Messages) > 0 {

		message := SQSMessage{}
		err := json.Unmarshal([]byte(*result.Messages[0].Body), &message)

		if err != nil {
			return nil, err
		}

		fileMessage := &models.QueueMessage{
			Filename: message.Records[0].S3.Object.Key,
			Size:     message.Records[0].S3.Object.Size,
			Id:       *result.Messages[0].ReceiptHandle,
		}

		return fileMessage, nil
	}

	return nil, nil
}
