import decode from 'jwt-decode';

export default class AuthHelperMethods {
    
    // Initializing important variables

    login = (username, password) => {
        
        // Get a token from api server using the fetch api
        return this.fetch(`http://localhost:8090/api/auth/login`, {
            method: 'POST',
            body: JSON.stringify({
                'username': username,
                'password': password
            })
        }).then(res => {
            this.setToken(res.auth_token) // Setting the token in localStorage
            return Promise.resolve(res);
        })
    }


    loggedIn = () => {
        // Checks if there is a saved token and it's still valid
        const token = this.getToken() // Getting token from localstorage
        return !!token && !this.isTokenExpired(token) // handwaiving here
    }

    isTokenExpired = (token) => {
        console.log(token)
        try {
            const decoded = decode(token);
            if (decoded.expiration_interval_in_minutes < Date.now() / 1000) { // Checking if token is expired.
                return true;
            }
            else
                return false;
        }
        catch (err) {
            console.log("expired check failed! Line 42: AuthService.js");
            return false;
        }
    }

    setToken = (idToken) => {
        // Saves user token to localStorage
        localStorage.setItem('auth_token', idToken)
    }

    getToken = () => {
        // Retrieves the user token from localStorage
        return localStorage.getItem('auth_token')
    }

    logout = () => {
        // Clear user token and profile data from localStorage
        localStorage.removeItem('auth_token');
    }

    getConfirm = () => {
        // Using jwt-decode npm package to decode the token
        let answer = decode(this.getToken());
        console.log("Recieved answer!");
        return answer;
    }

    fetch = (url, options) => {
        // performs api calls sending the required authentication headers
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
        // Setting Authorization header
        // Authorization: Bearer xxxxxxx.xxxxxxxx.xxxxxx
        if (this.loggedIn()) {
            headers['Authorization'] = 'Bearer ' + this.getToken()
        }

        return fetch(url, {
            headers,
            ...options
        })
            .then(this._checkStatus)
            .then(response => response.json())
    }

    _checkStatus = (response) => {
        // raises an error in case response status is not a success
        let status = response.status

        if (status === 200) { // Success status
            return response
        } else { // Bad Request
            var error = new Error("Incorrect username and password!")
            error.response = response
            throw error
        }
    }
}