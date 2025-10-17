import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { UserTableComponent } from '../user-table/user-table.component';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIcon } from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {FormsModule} from '@angular/forms';
import { MatSelectModule } from '@angular/material/select';

@Component({
  selector: 'app-home',
  imports: [
    UserTableComponent, 
    MatButtonModule, 
    MatDialogModule,
    MatIcon,
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  readonly dialog = inject(MatDialog);

  openDialog() {
    const dialogRef = this.dialog.open(DialogCreateUser, {autoFocus: false});

    dialogRef.afterClosed().subscribe(result => {
      console.log(`Dialog result: ${result}`);
    });
  }
}

@Component({
  selector: 'dialog-create-user',
  templateUrl: 'dialog-create-user.html',
  imports: [
    MatDialogModule, 
    MatButtonModule, 
    MatInputModule, 
    MatFormFieldModule, 
    FormsModule,
    MatSelectModule,
    MatIcon
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DialogCreateUser {}