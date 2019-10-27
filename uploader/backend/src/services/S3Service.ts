import AWS from 'aws-sdk';
import {promisify} from 'util';

export default class S3Service {

  private readonly name: string;
  private readonly s3: AWS.S3;

  constructor(region: string, name: string, apiVersion: string){
    this.name = name;
    this.s3 = new AWS.S3({apiVersion, region})

  }


  public async upload(file: Buffer, extension: string): Promise<AWS.S3.Types.PutObjectRequest>{
    const uploadParams = {
      Bucket: this.name,
      Key: `${Date.now()}${extension}`,
      Body: file
    };

    return promisify(this.s3.upload.bind(this.s3)) (uploadParams)
  }

}