import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter as Router, Switch } from 'react-router-dom'
import { Provider } from 'react-redux'
import store from './redux/store'

// Css
import './index.css'
import './css/tailwind.css'

// Pages
import App from './components/App'
import DashboardPage from './pages/DashboardPage'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import VerifyEmailPage from './pages/VerifyEmailPage'
import ResendEmailVerificationPage from './pages/ResendEmailVerificationPage'
import SendResetPasswordEmailPage from './pages/SendResetPasswordEmailPage'
import ResetPasswordPage from './pages/ResetPasswordPage'
import NotFoundPage from './pages/NotFoundPage'

// Routes
import SharedRoute from './routes/SharedRoute'
import PrivateRoute from './routes/PrivateRoute'
import PublicRoute from './routes/PublicRoute'
import * as serviceWorker from './serviceWorker'

ReactDOM.render(
    <Provider store={store}>
        <Router>
            <Switch>
                <SharedRoute exact path='/' component={App} />
                <PublicRoute exact path='/login' component={LoginPage} />
                <PublicRoute exact path='/signup' component={SignupPage} />
                <PrivateRoute exact path='/dashboard' component={DashboardPage} />
                <SharedRoute exact path="/verify" component={ResendEmailVerificationPage} />
                <SharedRoute exact path="/verify/:token" component={VerifyEmailPage} />
                <SharedRoute exact path="/reset-password" component={SendResetPasswordEmailPage} />
                <SharedRoute exact path="/reset-password/:token" component={ResetPasswordPage} />
                <SharedRoute path="" component={NotFoundPage} />
            </Switch>
        </Router>
    </Provider>
    , document.getElementById('root')
);

// If you want your areact-router-dompp to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()
