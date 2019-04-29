import decode from 'jwt-decode';

const AUTH_TOKEN_KEY = 'auth_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

export const login = (username: string, password: string) => {
    // Get a token from api server using the fetch api
    return request(`http://localhost:8090/api/auth/login`, {
        method: 'POST',
        body: JSON.stringify({
            username,
            password,
        })
    }).then(res => {
        setTokens(res.auth_token, res.refresh_token); // Setting the token in localStorage
        return Promise.resolve(res);
    })
};

export const signup = (email: string, username: string, password: string) => {
    // Get a token from api server using the fetch api
    return request(`http://localhost:8090/api/auth/register`, {
        method: 'POST',
        body: JSON.stringify({
            email,
            username,
            password,
        })
    }).then(res => {
        return Promise.resolve(res);
    })
};

// Invalidate the refresh_token on the server
export const logout = () => {

    const headers: { [key:string]:string } = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };

    // Setting Authorization header
    if (loggedIn()) {
        headers['Authorization'] = 'Bearer ' + getAuthToken()
    }

    return fetch(`http://localhost:8090/api/auth/logout`, {
        headers,
        method: 'POST',
    }).then(() => {
        // Clear tokens from the localStorage
        localStorage.removeItem(AUTH_TOKEN_KEY);
        localStorage.removeItem(REFRESH_TOKEN_KEY);
    })
};

export const sendVerificationEmail = (email: string) => {
    // Get a token from api server using the fetch api
    return request(`http://localhost:8090/api/auth/verify`, {
        method: 'POST',
        body: JSON.stringify({
            email,
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
            email,
        })
    }).then(_checkStatus)
      .then(res => {
        return Promise.resolve(res);
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
            password,
        })
    }).then(_checkStatus)
      .then(res => {
        return Promise.resolve(res);
    })
};

export const loggedIn = (): boolean => {
    // Checks if there is a saved token and it's still valid
    const authToken = getAuthToken();
    const refreshToken = getRefreshToken();

    if (refreshToken === null || refreshToken === undefined) {
        return false;
    }

    if (authToken === null || authToken === undefined) {
        return false;
    }

    // If the authToken is invalid
    if (isTokenExpired(authToken)) {

        // If the refreshToken is valid, refresh the tokens
        if (!isTokenExpired(refreshToken)) {
            refreshAuth(refreshToken);
            return true;
        }
        return false
    }

    return true;
};

const refreshAuth = (refreshToken: string) => {
    // Refresh the tokens
    const headers: { [key:string]:string } = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };

    fetch(`http://localhost:8090/api/auth/refresh`, {
        headers,
        method: 'POST',
        body: JSON.stringify({
            'refresh_token': refreshToken,
        })
    }).then(_checkStatus)
        .then((res: Response) => res.json())
        .then((res: TokenServerResponse) => {
            setTokens(res.auth_token, res.refresh_token);
        });
};

interface TokenServerResponse {
    exp: number,
    auth_token: string,
    refresh_token: string,
}

// The model for the token received
interface TokenType {
    exp: number,
}

const isTokenExpired = (token: string) => {
    try {
        // Decode
        const decoded: TokenType = decode(token);

        // Get Dates
        let expDate = new Date(decoded.exp * 1000);
        let currentDate = new Date(Date.now());

        // Check if expiration date hasn't been reached yet
        return expDate < currentDate;
    }
    catch (err) {
        return false
    }
};

// Saves user token to localStorage
const setTokens = (authToken: string, refreshToken: string) => {
    localStorage.setItem(AUTH_TOKEN_KEY, authToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
};

// Retrieves the user token from localStorage
const getAuthToken = () => localStorage.getItem(AUTH_TOKEN_KEY);
const getRefreshToken = () => localStorage.getItem(REFRESH_TOKEN_KEY);

// Using jwt-decode npm package to decode the token
export const getProfile = () => {
    if (!loggedIn()) {
        return;
    }

    let token = getAuthToken();

    if (token !== null) {
        return decode(token as string);
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
        headers['Authorization'] = 'Bearer ' + getAuthToken()
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
        let error = new Error(response.statusText);
        throw error
    }
};