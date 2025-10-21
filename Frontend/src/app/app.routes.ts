import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { StubComponent } from './stub/stub.component';

export const routes: Routes = [
    { path: '', component: HomeComponent },
    { path: 'stub', component: StubComponent },
];
