import { Location } from '@angular/common';
import {ComponentFixture, TestBed} from '@angular/core/testing';
import {Department, HomeComponent, UserStatus} from './home.component';
import {By} from '@angular/platform-browser';
import {provideRouter, Router} from '@angular/router';
import {provideLocationMocks} from '@angular/common/testing';

describe('HomeComponent', () => {
  let fixture: ComponentFixture<HomeComponent>;
  let component: HomeComponent;
  let router: Router;
  let location: Location;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ HomeComponent /* and possibly other modules */ ],
      providers: [
        provideRouter([
          { path: '', component: HomeComponent },
          { path: 'create', component: HomeComponent },
          { path: 'edit/:id', component: HomeComponent },
          { path: 'delete/:id', component: HomeComponent },
        ]),
        provideLocationMocks()
      ]
    }).compileComponents();

    router = TestBed.inject(Router);
    location = TestBed.inject(Location);
    fixture = TestBed.createComponent(HomeComponent);
    fixture.componentInstance.users = [
      {
        user_id: 1,
        user_name: 'bob',
        first_name: 'Bob',
        last_name: 'Smith',
        email: 'bob@test.com',
        user_status: UserStatus.ACTIVE,
        department: Department.NA
      }
    ];
    fixture.detectChanges();
    component = fixture.componentInstance;
    router.initialNavigation();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should open create dialog with query params', async () => {
    fixture.detectChanges();
    const createButton = fixture.debugElement.query(By.css('.create-btn'));
    createButton.nativeElement.click();
    await fixture.whenStable();

    expect(location.path()).toBe('/?create=true');
  });

  it('should open edit dialog with query params', async () => {
    fixture.detectChanges();
    const editButton = fixture.debugElement.query(By.css('.edit-btn'));
    editButton.nativeElement.click();
    await fixture.whenStable();

    expect(location.path()).toBe('/?edit=1');
  });

  it('should open delete dialog with query params', async () => {
    fixture.detectChanges();
    const deleteButton = fixture.debugElement.query(By.css('.delete-btn'));
    deleteButton.nativeElement.click();
    await fixture.whenStable();

    expect(location.path()).toBe('/?delete=1');
  });
});
