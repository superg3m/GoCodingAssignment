import { ChangeDetectionStrategy, Component, inject, Input, Output, EventEmitter } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatIcon } from '@angular/material/icon';

import { MatInput, MatInputModule } from "@angular/material/input";
import { MatFormFieldModule } from "@angular/material/form-field";
import { CommonModule,  } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';

import {MAT_DIALOG_DATA} from '@angular/material/dialog';
import { Inject } from '@angular/core';

export enum UserStatus {
  ACTIVE = "active",
  INACTIVE = "inactive",
  TERMINATED = "terminated",
}

export enum Department {
  ENGINEER = "engineer",
  SALES = "sales",
  HR = "hr",
}

export interface User {
  id: number,
  username: string;
  firstname: string;
  lastname: string;
  email: string;
  userStatus: UserStatus;
  department?: Department;
}

@Component({
  selector: 'app-home',
  imports: [
    MatButtonModule, 
    MatDialogModule,
    MatIcon,
    MatTableModule
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  readonly dialog = inject(MatDialog);

  displayedColumns: string[] = ['username', 'firstname', 'lastname', 'email', 'userStatus', 'department', 'edit', 'delete'];

  users: User[] = [
    {
      id: 1,
      username: 'JakeG32',
      firstname: 'Jake',
      lastname: 'Gore',
      email: 'JakeG32@gmail.com',
      userStatus: UserStatus.ACTIVE,
      department: Department.ENGINEER
    },
    // {name: 'WillS-23', weight: 4.0026, symbol: 'He'},
    // {name: 'JamesB', weight: 6.941, symbol: 'Li'},
    // {name: 'BarryA90', weight: 9.0122, symbol: 'Be'},
  ];

  openCreateDialog() {
    const dialogRef = this.dialog.open(DialogSaveUserComponent, {
      autoFocus: false,
      data: { user: undefined }
    });

    dialogRef.afterClosed().subscribe((result: User|undefined) => {
      if (result) {
        console.log(result);
        this.users = [...this.users, result];
        // api request here
        console.log(this.users);
      }
    });
  }

  openEditDialog(userIndex: number) {
    const dialogRef = this.dialog.open(DialogSaveUserComponent, {
      autoFocus: false,
      data: { user: this.users[userIndex] }
    });

    dialogRef.afterClosed().subscribe((result: User|undefined) => {
      if (result) {
        console.log(result)
        this.users[userIndex] = result
        this.users = [...this.users]
      }
    });
  }

  openDeleteDialog(userIndex: number) {
    const dialogRef = this.dialog.open(DialogDeleteUserComponent, {
      autoFocus: false,
      data: { user: this.users[userIndex] }
    });

    dialogRef.afterClosed().subscribe((result: boolean) => {
      if (result) {
        // Node(jovanni): This is bad for performance because you have to do a syscall to reallocate
        // but im not too worried for this case.
        this.users = this.users.filter((_, i) => {
          return i != userIndex;
        })
      }
    });
  }
}

@Component({
  selector: 'dialog-save-user',
  templateUrl: 'dialog-save-user.component.html',
  imports: [
    MatDialogModule, 
    MatButtonModule, 
    MatInputModule, 
    MatFormFieldModule, 
    FormsModule,
    MatSelectModule,
    MatIcon,
    CommonModule
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DialogSaveUserComponent {
  user: User;

  constructor(
    private dialogRef: MatDialogRef<DialogSaveUserComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { user: User }
  ) {
    this.user = data.user ? { ...data.user } : {
      id: -1,
      username: '',
      firstname: '',
      lastname: '',
      email: '',
      userStatus: UserStatus.ACTIVE,
      department: Department.ENGINEER
    };
  }

  readonly userStatusValues = Object.values(UserStatus);
  readonly departmentValues = Object.values(Department);

  saveUser() {
    this.dialogRef.close(this.user)
  }

  cancel() {
    console.log(this.user)
    this.dialogRef.close()
  }

  isFormFilledOut() {
    return (
      this.user.username&&
      this.user.firstname &&
      this.user.lastname &&
      this.user.email && 
      this.user.userStatus != <UserStatus>(<unknown>-1) &&
      this.user.department != <Department>(<unknown>-1)
    );
  }
}

@Component({
  selector: 'dialog-delete-user',
  templateUrl: 'dialog-delete-user.component.html',
  imports: [
    MatDialogModule, 
    MatButtonModule, 
    MatInputModule, 
    MatFormFieldModule, 
    FormsModule,
    MatSelectModule,
    MatIcon,
    CommonModule
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DialogDeleteUserComponent {
  constructor(
    private dialogRef: MatDialogRef<DialogDeleteUserComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { user: User }
  ) {}

  readonly userStatusValues = Object.values(UserStatus);
  readonly departmentValues = Object.values(Department);

  deleteUser() {
    this.dialogRef.close(true)
  }

  cancel() {
    this.dialogRef.close(false)
  }
}