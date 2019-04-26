import { LOGIN } from '../constants'
import { login as authLogin } from '../../utils/AuthService'

export interface Action {
    type: string,
    payload: object,
}

// TODO: Do I need this?!!
export function login(username: string, password: string) {
    authLogin(username, password).then(res => res.json())
        .then((res: any) => {
            return {
                type: LOGIN,
                payload: {
                    username,
                    email: res.email,
                }
            }
        });
}