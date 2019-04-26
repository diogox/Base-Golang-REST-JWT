import decode from 'jwt-decode';

const AUTH_TOKEN_KEY = 'auth_token'

export const login = (username: string, password: string) => {
    // Get a token from api server using the fetch api
    return request(`http://localhost:8090/api/auth/login`, {
        method: 'POST',
        body: JSON.stringify({
            username,
            password
        })
    }).then(res => {
        setToken(res.auth_token) // Setting the token in localStorage
        return Promise.resolve(res)
    })
};

export const signup = (email: string, username: string, password: string) => {
    // Get a token from api server using the fetch api
    return request(`http://localhost:8090/api/auth/register`, {
        method: 'POST',
        body: JSON.stringify({
            email,
            username,
            password
        })
    }).then(res => {
        return Promise.resolve(res)
    })
};

export const sendVerificationEmail = (email: string) => {
    // Get a token from api server using the fetch api
    return request(`http://localhost:8090/api/auth/verify`, {
        method: 'POST',
        body: JSON.stringify({
            email
        })
    }).then(res => {
        return Promise.resolve(res)
    })
};

export const verifyEmail = (token: string) => {
    // Get a token from api server using the fetch api
    return fetch(`http://localhost:8090/api/auth/verify/${token}`, {
        method: 'GET',
    }).then(res => {
        return Promise.resolve(res)
    })
};

export const sendResetPasswordEmail = (email: string) => {
    // Get a token from api server using the fetch api
    const headers: { [key:string]:string } = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };

    return fetch('http://localhost:8090/api/auth/reset-password', {
        headers,
        method: 'POST',
        body: JSON.stringify({
            email
        })
    }).then(_checkStatus)
      .then(res => {
        return Promise.resolve(res)
    })
};

export const resetPassword = (token: string, password: string) => {
    // Get a token from api server using the fetch api
    const headers: { [key:string]:string } = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };

    return fetch(`http://localhost:8090/api/auth/reset-password/${token}`, {
        headers,
        method: 'POST',
        body: JSON.stringify({
            password
        })
    }).then(_checkStatus)
      .then(res => {
        return Promise.resolve(res)
    })
};

export const loggedIn = () => {
    // Checks if there is a saved token and it's still valid
    const token = getToken(); // Getting token from localstorage

    // If there's no value, return
    if (token === undefined) {
        return false
    }

    // Check validity
    return (!!token) && !isTokenExpired(token)
};

// The model for the token received
interface TokenType {
    exp: number,
    auth_token: string,
}

const isTokenExpired = (token: string) => {
    try {
        const decoded: TokenType = decode(token)
        if (decoded.exp < Date.now() / 1000) { // Checking if token is expired.
            return true
        }
        else
            return false
    }
    catch (err) {
        return false
    }
};

// Saves user token to localStorage
const setToken = (authToken: string) => localStorage.setItem(AUTH_TOKEN_KEY, authToken)

// Retrieves the user token from localStorage
const getToken = () => localStorage.getItem(AUTH_TOKEN_KEY)

// Clear user token and profile data from localStorage
export const logout = () => localStorage.removeItem(AUTH_TOKEN_KEY)

// Using jwt-decode npm package to decode the token
export const getProfile = () => {
    let token = getToken()

    if (token !== null) {
        return decode(token as string)
    }
};

const request = (url: string, options: any) => {
    // Performs api calls sending the required authentication headers
    const headers: { [key:string]:string } = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };

    // Setting Authorization header
    // Authorization: Bearer xxxxxxx.xxxxxxxx.xxxxxx
    if (loggedIn()) {
        headers['Authorization'] = 'Bearer ' + getToken()
    }

    return fetch(url, {
        headers,
        ...options
    })
        .then(_checkStatus)
        .then(response => response.json())
};

const _checkStatus = (response: Response) => {
    // Raises an error in case response status is not a success
    if (response.status >= 200 && response.status < 300) { // Success status lies between 200 to 300
        return response
    } else {
        let error = new Error(response.statusText)
        throw error
    }
}