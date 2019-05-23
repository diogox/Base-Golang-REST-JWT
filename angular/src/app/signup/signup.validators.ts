import {AbstractControl, ValidationErrors} from '@angular/forms';

export class SignupValidators {

  public static mustMatch(control: AbstractControl): ValidationErrors | null {
    const password = control.get('password').value;
    const passwordConfirmation = control.get('passwordConfirmation').value;

    if (password !== passwordConfirmation) {
      control.get('passwordConfirmation').setErrors({ notMatch: true });
    }

    return null;
  }

  public static shouldBeUniqueUsername(control: AbstractControl): Promise<ValidationErrors | null> {
    return new Promise((resolve, reject) => {
      // TODO: Make API call

      setTimeout(() => {
        if (control.value === 'diogox') {
          resolve({ shouldBeUniqueUsername: true });
        } else {
          resolve(null);
        }
      }, 2000);
    });
  }

  public static shouldBeUniqueEmail(control: AbstractControl): Promise<ValidationErrors | null> {
    return new Promise((resolve, reject) => {
      // TODO: Make API call

      setTimeout(() => {
        if (control.value === 'dxmp@iol.pt') {
          resolve({ shouldBeUniqueEmail: true });
        } else {
          resolve(null);
        }
      }, 2000);
    });
  }

}
