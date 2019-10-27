import dotenv from 'dotenv';

dotenv.config();

export default {
  bucket: {
    region: process.env.BUCKET_REGION,
    apiVersion: process.env.BUCKET_API_VERSION || '2006-03-01',
    name: process.env.BUCKET_NAME
  }
}