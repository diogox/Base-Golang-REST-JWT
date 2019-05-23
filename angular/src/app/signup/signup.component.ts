import { Component, OnInit } from '@angular/core';
import {AbstractControl, FormControl, FormGroup, NgForm, ValidationErrors, Validators} from '@angular/forms';
import {SignupValidators} from './signup.validators';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent implements OnInit {
  form = new FormGroup({
    username: new FormControl('', [
      Validators.required,
      Validators.minLength(3),
      Validators.pattern('^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$')
    ], SignupValidators.shouldBeUniqueUsername),
    email: new FormControl('', [
      Validators.required,
      Validators.email
    ], SignupValidators.shouldBeUniqueEmail),
    password: new FormControl('', [
      Validators.required,
      Validators.minLength(8)
    ]),
    passwordConfirmation: new FormControl()
  }, SignupValidators.mustMatch);

  constructor() { }

  ngOnInit() {
  }

  onSubmit() {
    // Ignore if fields are invalid
    if (this.form.invalid) {
      return;
    }
    // TODO: Make auth API call
    this.form.setErrors({
      invalidSignup: true
    });
  }

  get username(): AbstractControl {
    return this.form.get('username');
  }

  get email(): AbstractControl {
    return this.form.get('email');
  }

  get password(): AbstractControl {
    return this.form.get('password');
  }

  get passwordConfirmation(): AbstractControl {
    return this.form.get('passwordConfirmation');
  }
}
