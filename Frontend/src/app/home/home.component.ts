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

import {
  MatSnackBar,
  MatSnackBarHorizontalPosition,
  MatSnackBarVerticalPosition,
} from '@angular/material/snack-bar';

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
  private snackBar = inject(MatSnackBar);
  readonly dialog = inject(MatDialog);

  displayedColumns: string[] = ['id', 'user_name', 'first_name', 'last_name', 'email', 'user_status', 'department', 'edit', 'delete'];

  users: User[] = [];

  showSuccess(message: string) {
    this.snackBar.open(message, 'Close', {
      duration: 5000,
      horizontalPosition: "end",
      verticalPosition: "top",
      panelClass: ['snackbar-success']
    });
  }

  showError(message: string) {
    this.snackBar.open(message, 'Close', {
      duration: 5000,
      horizontalPosition: "end",
      verticalPosition: "top",
      panelClass: ['snackbar-error']
    });
  }

  openCreateDialog() {
    const dialogRef = this.dialog.open(DialogSaveUserComponent, {
      autoFocus: false,
      data: {user: undefined}
    });

    dialogRef.afterClosed().subscribe(async (result: User | undefined) => {
      if (result) {
        try {
          const requestBody = {
            "user_name": result.user_name,
            "first_name": result.first_name,
            "last_name": result.last_name,
            "email": result.email,
            "user_status": result.user_status,
            "department": result.department,
          }

          const response = await fetch("http://localhost:8080/User/Create", {
            method: "POST",
            credentials: "include",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
          });

          if (!response.ok) {
            this.showError("Failed to create user: " + await response.text())
            return
          }

          const responseBody: User = await response.json()

          // Node(jovanni): This is bad for performance because you have to do a syscall to reallocate
          // but im not too worried about this case. I have to trigger change detection...
          this.users = [...this.users, responseBody];
          this.showSuccess("Successfully created the user!");
        } catch (e) {
          // This can be a toast
          console.log(e)
          this.showError("Failed to create user: " + e)
        }
      }
    });
  }

  openEditDialog(userIndex: number) {
    const dialogRef = this.dialog.open(DialogSaveUserComponent, {
      autoFocus: false,
      data: {user: this.users[userIndex]}
    });

    dialogRef.afterClosed().subscribe(async (result: User | undefined) => {
      if (result) {
        try {
          const requestBody = {
            "user_id": this.users[userIndex].user_id,
            "user_name": result.user_name,
            "first_name": result.first_name,
            "last_name": result.last_name,
            "email": result.email,
            "user_status": result.user_status,
            "department": result.department,
          }

          const response = await fetch("http://localhost:8080/User/Update", {
            method: "PATCH",
            credentials: "include",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(requestBody)
          });

          if (!response.ok) {
            this.showError("Failed to update user: " + await response.text());
            return;
          }

          // Node(jovanni): This is bad for performance because you have to do a syscall to reallocate
          // but im not too worried about this case. I have to trigger change detection...
          this.users[userIndex] = result;
          this.users = [...this.users];
          this.showSuccess("Successfully updated the user!");
        } catch (e) {
          this.showError("Failed to update" + e);
        }
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

          // Node(jovanni): This is bad for performance because you have to do a syscall to reallocate
          // but im not too worried about this case.
          this.users = this.users.filter((_, i) => {
            return i != userIndex;
          })
          this.showSuccess("Successfully deleted the user!");
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
