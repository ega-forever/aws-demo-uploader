package interfaces

import "github.com/ega-forever/aws-demo-uploader/internal/models"

type Bucket interface {
	GetFile(filename string) (*[]byte, error)
}

type Queue interface {
	Subscribe() (chan *models.QueueMessage, chan error)
	Unsubscribe()
	AckMessage(messageId string) error
}

type Database interface {
	SaveRecord(name string, sirname string, score int64) error
}
