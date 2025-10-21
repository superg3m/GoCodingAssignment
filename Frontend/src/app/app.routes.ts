import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { StubComponent } from './stub/stub.component';

// Look at the router stuff to do { path: '/User/Edit/:id', component: StubComponent }
export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'stub', component: StubComponent },
  { path: 'create', component: HomeComponent },
  { path: 'edit/:id', component: HomeComponent },
  { path: 'delete/:id', component: HomeComponent },
];
