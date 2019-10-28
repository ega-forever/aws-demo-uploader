import express from 'express';
import fileUpload from 'express-fileupload';
import cors from 'cors';
import bodyParser from 'body-parser';
import uploadRoutes from './routes/uploadRoutes';
import S3Service from './services/S3Service';
import config from './config';
import DBService from './services/DBService';
import examRoutes from './routes/examRoutes';
import * as bunyan from 'bunyan';
import * as path from 'path';

const logger = bunyan.createLogger({ name: 'backend' });

const init = async () => {

  const app = express();

  app.use(fileUpload({
    createParentPath: true
  }));

  app.use(cors());
  app.use(bodyParser.json());
  app.use(bodyParser.urlencoded({ extended: true }));

  app.use('/', express.static('public'));

  const s3Service = new S3Service(config.bucket.region, config.bucket.name, config.bucket.apiVersion);
  const dbService = new DBService(config.db.uri);

  await dbService.connect();
  await dbService.migrate();

  dbService.once('error', (err) => {
    logger.error(err);
    process.exit(1);
  });

  app.use('/upload', uploadRoutes(s3Service));
  app.use('/exams', examRoutes(dbService));

  app.listen(config.rest.port, () =>
    logger.info(`App is listening on port ${ config.rest.port }.`)
  );

};

module.exports = init().catch(e => {
  logger.error(e);
  process.exit(1);
});
