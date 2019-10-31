package services

import (
	"github.com/ega-forever/aws-demo-uploader/internal/interfaces"
	"github.com/tealeg/xlsx"
	"log"
	"strconv"
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

func (ss *ProcessService) Listen() error {

	eventCh, errCh := ss.queue.Subscribe()

	var serviceError error

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {

		for {
			select {
			case m := <-eventCh:
				{
					bytes, err := ss.bucket.GetFile(m.Filename)

					if err != nil {
						serviceError = err
						wg.Done()
					}

					err = ss.writeFileBufferToDb(bytes)

					if err != nil {
						serviceError = err
						wg.Done()
					}

					err = ss.queue.AckMessage(m.Id)

					if err != nil {
						serviceError = err
						wg.Done()
					}
				}
			case e := <-errCh:
				{
					log.Fatal(e)
				}

			}
		}

	}()

	wg.Wait()

	return serviceError
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

			name := row.Cells[0].Value
			sirname := row.Cells[1].Value
			score, _ := strconv.Atoi(row.Cells[2].Value)

			err = ss.database.SaveRecord(name, sirname, int64(score))

			if err != nil {
				return err
			}
		}
	}

	return nil
}
