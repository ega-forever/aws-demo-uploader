import express from 'express';
import S3Service from '../services/S3Service';
import DBService from '../services/DBService';
import asyncMiddleware from '../middlewares/asyncMiddleware';


export default function (dbService: DBService): express.Router {

  const router = express.Router();

  router.get('/', asyncMiddleware(async (req: express.Request, res: express.Response) => {
    const data = await dbService.getExamResults();
    return res.send(data.rows);
  }));


  return router;
}