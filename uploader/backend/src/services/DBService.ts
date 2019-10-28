import { Client } from 'pg';
import { EventEmitter } from 'events';

export default class DBService extends EventEmitter {

  private client: Client;

  constructor(uri: string) {
    super();
    this.client = new Client({
      connectionString: uri
    });

    this.client.on('error', err => {
      this.emit('error', err);
    })
  }

  public async connect() {
    await this.client.connect();
  }

  public async getExamResults(): Promise<{ rows: Array<{ name: string, sirname: string, score: number }> }> {
    return await new Promise((res, rej) => {
      this.client.query('SELECT * from exam_results LIMIT 10', (err, results) => {
        err ? rej(err) : res(results);
      })
    })


  }


}