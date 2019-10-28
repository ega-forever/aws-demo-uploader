package main

import (
	"github.com/ega-forever/aws-demo-uploader/internal/bucket"
	database_sql "github.com/ega-forever/aws-demo-uploader/internal/database"
	"github.com/ega-forever/aws-demo-uploader/internal/queue"
	"github.com/ega-forever/aws-demo-uploader/internal/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("QUEUE_URI", "")
	viper.SetDefault("QUEUE_REGION", "")
	viper.SetDefault("QUEUE_API_VERSION", "")

	viper.SetDefault("BUCKET_REGION", "")
	viper.SetDefault("BUCKET_API_VERSION", "")
	viper.SetDefault("BUCKET_NAME", "")

	viper.SetDefault("LOG_LEVEL", 30)

	viper.SetDefault("DATABASE_URI", "postgres://user:user@localhost:5432/app")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	logLevel := viper.GetInt("LOG_LEVEL")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(logLevel))

}

func main() {

	queueUri := viper.GetString("QUEUE_URI")
	queueRegion := viper.GetString("QUEUE_REGION")

	bucketName := viper.GetString("BUCKET_NAME")
	bucketRegion := viper.GetString("BUCKET_REGION")

	databaseUri := viper.GetString("DATABASE_URI")

	sqs := queue_sqs.New(queueUri, queueRegion, 1, 10)
	s3 := bucket_s3.New(bucketName, bucketRegion)
	db, err := database_sql.New(databaseUri)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrate()

	if err != nil {
		log.Fatal(err)
	}

	ps := services.NewProcessService(s3, sqs, db)

	err = ps.Listen()

	if err != nil {
		log.Fatal(err)
	}

}
