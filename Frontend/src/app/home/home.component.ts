import {ChangeDetectionStrategy, Component, inject, Input, Output, EventEmitter, OnInit} from '@angular/core';
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

import validator from 'validator';

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
  user_id: number,
  user_name: string;
  first_name: string;
  last_name: string;
  email: string;
  user_status: UserStatus;
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
export class HomeComponent implements OnInit {
  readonly dialog = inject(MatDialog);

  displayedColumns: string[] = ['id', 'user_name', 'first_name', 'last_name', 'email', 'user_status', 'department', 'edit', 'delete'];

  users: User[] = [];

  openCreateDialog() {
    const dialogRef = this.dialog.open(DialogSaveUserComponent, {
      autoFocus: false,
      data: {user: undefined}
    });

    dialogRef.afterClosed().subscribe((result: User | undefined) => {
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
      data: {user: this.users[userIndex]}
    });

    dialogRef.afterClosed().subscribe((result: User | undefined) => {
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
      data: {user: this.users[userIndex]}
    });

    dialogRef.afterClosed().subscribe(async (result: boolean) => {
      if (result) {
        // Node(jovanni): This is bad for performance because you have to do a syscall to reallocate
        // but im not too worried about this case.
        try {
          const requestBody = {
            "user_id": this.users[userIndex].user_id
          }

          const response = await fetch("http://localhost:8080/User/Delete", {
            method: "DELETE",
            credentials: "include",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
          });

          this.users = this.users.filter((_, i) => {
            return i != userIndex;
          })
        } catch (e) {
          // This can be a toast
          console.log("Failed to delete" + e)
        }
      }
    });
  }

  async ngOnInit() {
    try {
      const response = await fetch("http://localhost:8080/User/Get/All", {
        method: "GET",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        }
      });

      this.users = await response.json();
    } catch (e) {
      console.log("Can't reach the server")
    }
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
      user_id: -1,
      user_name: '',
      first_name: '',
      last_name: '',
      email: '',
      user_status: UserStatus.ACTIVE,
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

  isEmailValid(): boolean {
    return validator.isEmail(this.user.email);
  }

  isFormFilledOut(): boolean {
    return (
      this.user.user_name &&
      this.user.first_name &&
      this.user.last_name &&
      this.user.email &&
      this.user.user_status != <UserStatus>(<unknown>-1) &&
      this.user.department != <Department>(<unknown>-1)
    ) == true;
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
