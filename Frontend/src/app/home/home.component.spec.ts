import {ComponentFixture, TestBed} from '@angular/core/testing';

import {Department, HomeComponent, UserStatus} from './home.component';
import {By} from '@angular/platform-browser';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HomeComponent]
    }).compileComponents();

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
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should open create dialog', () => {
    spyOn(component, "openCreateDialog")
    const editButton = fixture.debugElement.query(By.css('.create-btn'));
    editButton.triggerEventHandler('click', null);

    expect(component.openCreateDialog).toHaveBeenCalled();
  })

  it('should open edit dialog', () => {
    spyOn(component, "openEditDialog")
    const editButton = fixture.debugElement.query(By.css('.edit-btn'));
    editButton.triggerEventHandler('click', null);

    expect(component.openEditDialog).toHaveBeenCalled();
  })

  it('should open delete dialog', () => {
    spyOn(component, "openDeleteDialog")
    const deleteButton = fixture.debugElement.query(By.css('.delete-btn'));
    deleteButton.triggerEventHandler('click', null);

    expect(component.openDeleteDialog).toHaveBeenCalled();
  })
});
