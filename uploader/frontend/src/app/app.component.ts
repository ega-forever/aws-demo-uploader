import { Component } from '@angular/core';
import { Record } from './models/Record';
import { Observable, ObservableInput, throwError } from 'rxjs';
import { FormControl, FormGroup } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError, tap } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  private title: string = 'uploader';
  private records$: Promise<Record[]>;
  private filter: FormControl = new FormControl('');
  private filterFormGroup: FormGroup = new FormGroup({ filter: this.filter });
  private uploadFormGroup: FormGroup = new FormGroup({ });
  private file: File | null = null;
  private fileUploadInProcess: boolean = false;


  constructor(private http: HttpClient) {
    this.filter.valueChanges.subscribe(async (text: string) => {
      if (text.length > 4)
        this.records$ = this.search(text);
    })

  }

  // todo implement
  private async search(text: string): Promise<Record[]> {
    console.log('super')
    return [
      {
        name: 'Test',
        sirname: Date.now().toString(),
        age: Math.random() * 10
      }
    ]
  }

  private async handleFileInput(files: FileList) {
    const file = files.item(0);

    if (!file.type.includes('excel') && !file.type.includes('openxml'))
      return;

    this.file = file;
  }

  private handleUploadError(e: HttpErrorResponse): ObservableInput<any> {
    console.log(e);
    return throwError(
      'Something bad happened; please try again later.');
  }

  postFile(fileToUpload: File): void {
    const endpoint = 'http://localhost:3000/upload';
    const formData: FormData = new FormData();
    formData.append('results', fileToUpload, fileToUpload.name);
    this.http
      .post(endpoint, formData, { headers: { Client: 'web' } })
      .pipe(
        catchError(this.handleUploadError)
      ).subscribe((result) => {
      console.log(result);
      this.fileUploadInProcess = false;
      this.file = null;
    })
  }


}
