import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import './index.css';
import './css/tailwind.css';
import App from './containers/App';
import LoginPage from './containers/LoginPage';
import NotFoundPage from './containers/NotFoundPage';
import PrivateRoute from './containers/PrivateRoute';
import PublicRoute from './containers/PublicRoute';
import * as serviceWorker from './serviceWorker';

ReactDOM.render(
    <Router>
        <Switch>
            <PublicRoute exact path='/' component={App} />
            <PublicRoute exact path='/login' component={LoginPage} />
            <PrivateRoute exact path='/dashboard' component={App} />
            <Route path="" component={NotFoundPage} />
        </Switch>
    </Router>
    , document.getElementById('root'));

// If you want your areact-router-dompp to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
