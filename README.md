# Data processor

## Description
the application for monitoring students exam scores. The app includes 2 services: 

#### uploader
A file uploader application with an ability to view the current scores and search for students (the test data can be found in ```test_data``` folder).

#### processor
An application, responsible for parsing the data from excel and push to postgres database.

## Requirements
1) postgres version 9.6
2) sqs
3) s3 bucket
4) registry (for docker images)
5) ecs or eks (for running docker images)

## environment

#### uploader

| variable | default value | description
| --- | --- | --- |
| AWS_ACCESS_KEY_ID | | aws service account's access key id
| AWS_SECRET_ACCESS_KEY | | aws service account's secret key
| PORT | 3000 | rest port
| BUCKET_REGION | | aws s3 bucket region
| BUCKET_API_VERSION | 2006-03-01 | aws s3 bucket api version (optional)
| BUCKET_NAME | | aws s3 bucket name
| DATABASE_URI | | postgres URI, for example ```postgres://user:123@localhost:5432/otus```

#### processor

| variable | default value | description
| --- | --- | --- |
| AWS_ACCESS_KEY_ID | | aws service account's access key id
| AWS_SECRET_ACCESS_KEY | | aws service account's secret key
| QUEUE_URI | | aws sqs URI
| QUEUE_REGION | | aws SQS region
| QUEUE_API_VERSION | | aws SQS api version
| BUCKET_REGION | | aws s3 bucket region
| BUCKET_API_VERSION | | aws s3 bucket version
| BUCKET_NAME | | aws s3 bucket name
| LOG_LEVEL |  | level of logging (check out [logrus](github.com/sirupsen/logrus) for more details)
| DATABASE_URI | | postgres URI, for example ```postgres://user:123@localhost:5432/otus```

#### Extra setup
Before deploying the app, make sure you already have:
1) service account with enough rights to push / pull messages from/to sqs, and to upload / read files from s3
2) the notification rule under your created s3 bucket: so when the s3 receives new file, it sends the notification via sqs. 
You can setup the rule in s3 (choose bucket, then properties->events, then add new rule on PUT)