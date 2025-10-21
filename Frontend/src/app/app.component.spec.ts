import {ComponentFixture, TestBed} from '@angular/core/testing';
import { AppComponent } from './app.component';
import { provideRouter} from '@angular/router';

describe('AppComponent', () => {
  let fixture: ComponentFixture<AppComponent>;
  let component: AppComponent
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppComponent],
      providers: [provideRouter([])]
    }).compileComponents();

    fixture = TestBed.createComponent(AppComponent);
    component = fixture.componentInstance
  });

  it('should create the app', () => {
    expect(component).toBeTruthy();
  });

  it(`should have the 'Frontend' title`, () => {
    expect(component.title).toEqual('Frontend');
  });
});
