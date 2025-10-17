import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { ResourcesComponent } from './resources/resources.component';

export const routes: Routes = [
    { path: '', component: HomeComponent },
    { path: 'resources', component: ResourcesComponent },
];