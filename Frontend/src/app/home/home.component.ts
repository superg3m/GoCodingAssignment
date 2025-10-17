import { Component } from '@angular/core';
import { UserTableComponent } from '../user-table/user-table.component';
import { MatIcon } from '@angular/material/icon';

@Component({
  selector: 'app-home',
  imports: [UserTableComponent, MatIcon],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {

}
