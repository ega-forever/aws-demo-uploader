import express from 'express';
import S3Service from '../services/S3Service';
import asyncMiddleware from '../middlewares/asyncMiddleware';


export default function (s3Service: S3Service): express.Router {

  const router = express.Router();

  router.post('/', asyncMiddleware(async (req: express.Request, res: express.Response) => {

    const file: { data: Buffer, name: string } = (req as any).files.results;

    if (!file) {
      return res.send({ ok: 0 });
    }

    const extensionMatch = file.name.match(/\.[0-9a-z]+$/);

    if (!extensionMatch.length) {
      return res.send({ ok: 0 });
    }

    await s3Service.upload(file.data, extensionMatch[0]);
    res.send({ ok: 1 });
  }));


  return router;
}