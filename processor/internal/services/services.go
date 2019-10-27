package services

import (
	"github.com/ega-forever/aws-demo-uploader/internal/interfaces"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"sync"
)

type ProcessService struct {
	bucket   interfaces.Bucket
	queue    interfaces.Queue
	database interfaces.Database
}

func NewProcessService(bucket interfaces.Bucket, queue interfaces.Queue, database interfaces.Database) *ProcessService {

	return &ProcessService{
		bucket:   bucket,
		queue:    queue,
		database: database,
	}

}

func (ss *ProcessService) Listen() {

	eventCh, _ := ss.queue.Subscribe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		for {
			select {
			case m := <-eventCh:
				{
					log.Print(m)
					bytes, err := ss.bucket.GetFile(m.Filename)

					if err != nil {
						wg.Done()
						log.Error(err)
					}

					err = ss.writeFileBufferToDb(bytes)

					if err != nil {
						wg.Done()
						log.Error(err)
					}

					// todo process excel
					//todo save to db
					/*err := sqs.AckMessage(m)

					if err != nil {
						log.Print(err)
					}*/

					//wg.Done()
				}

			}
		}

	}()

	wg.Wait()

}

func (ss *ProcessService) writeFileBufferToDb(data *[]byte) error {

	xlFile, err := xlsx.OpenBinary(*data)
	if err != nil {
		return err
	}
	for _, sheet := range xlFile.Sheets {
		for index, row := range sheet.Rows {
			if index == 0 {
				continue
			}

			log.Print(row)

			/*for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("%s\n", text)
			}*/
		}
	}

	return nil
}
