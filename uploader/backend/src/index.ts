import express from 'express';
import fileUpload from 'express-fileupload';
import cors from 'cors';
import bodyParser from 'body-parser';
import uploadRoutes from './routes/uploadRoutes';
import S3Service from './services/S3Service';
import config from './config';

const app = express();

app.use(fileUpload({
  createParentPath: true
}));

app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

const s3Service = new S3Service(config.bucket.region, config.bucket.name, config.bucket.apiVersion);

app.use('/upload', uploadRoutes(s3Service));

//start app
const port = process.env.PORT || 3000;

app.listen(port, () =>
  console.log(`App is listening on port ${ port }.`)
);