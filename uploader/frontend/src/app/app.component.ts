import { Component } from '@angular/core';
import { Record } from './models/Record';
import { ObservableInput, throwError } from 'rxjs';
import { FormControl, FormGroup } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError } from 'rxjs/operators';
import Timer = NodeJS.Timer;
import { environment } from '../environments/environment';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  public title: string = 'uploader';
  public records$: Promise<Record[]>;
  public filter: FormControl = new FormControl('');
  public filterFormGroup: FormGroup = new FormGroup({ filter: this.filter });
  public uploadFormGroup: FormGroup = new FormGroup({});
  public file: File | null = null;
  public fileUploadInProcess: boolean = false;
  public fileUploadCompleted: boolean = false;
  private watchExamIntervalPid: Timer;


  constructor(private http: HttpClient) {
    this.records$ = this.search('');

    this.watchExams();

    this.filter.valueChanges.subscribe(async (text: string) => {
      if (text.length > 2 || !text.length) {
        this.records$ = this.search(text);
      }

      !text.length ? this.watchExams() : this.clearWatchExams();
    })
  }

  private watchExams() {
    this.watchExamIntervalPid = setInterval(() => {
      this.records$ = this.search('');
    }, 5000);
  }

  private clearWatchExams() {
    clearInterval(this.watchExamIntervalPid);
  }

  private async search(text: string): Promise<Record[]> {
    const endpoint = `${ environment.backend }/exams?text=${ text }`;

    return await this.http
      .get(endpoint, { headers: { Client: 'web' } })
      .pipe(
        catchError(this.handleBackendError)
      ).toPromise();

  }

  public async handleFileInput(files: FileList) {
    const file = files.item(0);

    if (!file.type.includes('excel') && !file.type.includes('openxml'))
      return;

    this.file = file;
  }


  private handleBackendError(e: HttpErrorResponse): ObservableInput<any> {
    return throwError(
      'Something bad happened; please try again later.');
  }

  postFile(fileToUpload: File): void {
    const endpoint = `${ environment.backend }/upload`;
    const formData: FormData = new FormData();
    formData.append('results', fileToUpload, fileToUpload.name);
    this.http
      .post(endpoint, formData, { headers: { Client: 'web' } })
      .pipe(
        catchError(this.handleBackendError)
      ).subscribe((result) => {
      this.fileUploadInProcess = false;
      this.file = null;
      this.fileUploadCompleted = true;

    })
  }


}
