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


  public async migrate() {
    await this.makeRequest(`
      create table if not exists exam_results (
        id SERIAL,
        name VARCHAR(255),
        sirname VARCHAR(255),
        score double precision
      )`
    );

    await this.makeRequest(`CREATE INDEX IF NOT EXISTS exam_results_name_idx on exam_results using btree (name)`);
    await this.makeRequest(`CREATE INDEX IF NOT EXISTS exam_results_sirname_idx on exam_results using btree (sirname)`);

  }

  public async getExamResults(text: string): Promise<{ rows: Array<{ name: string, sirname: string, score: number }> }> {

    if (text.length) {
      return this.makeRequest('SELECT * from exam_results where name LIKE $1 or sirname LIKE $1 LIMIT 10', [`%${ text }%`]);
    }


    return this.makeRequest('SELECT * from exam_results LIMIT 10')
  }


  private async makeRequest(query: string, params: any[] = []): Promise<any> {
    return await new Promise((res, rej) => {
      this.client.query(query, params, (err, results) => {
        err ? rej(err) : res(results);
      })
    })
  }


}