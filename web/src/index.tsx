import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import './index.css'
import './css/tailwind.css'
import App from './components/App'
import DashboardPage from './pages/DashboardPage'
import LoginPage from './pages/LoginPage'
import NotFoundPage from './pages/NotFoundPage'
import SharedRoute from './routes/SharedRoute'
import PrivateRoute from './routes/PrivateRoute'
import PublicRoute from './routes/PublicRoute'
import * as serviceWorker from './serviceWorker'

ReactDOM.render(
    <Router>
        <Switch>
            <SharedRoute exact path='/' component={App} />
            <PublicRoute exact path='/login' component={LoginPage} />
            <PrivateRoute exact path='/dashboard' component={DashboardPage} />
            <Route path="" component={NotFoundPage} />
        </Switch>
    </Router>
    , document.getElementById('root')
)

// If you want your areact-router-dompp to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()
