import React from 'react';  
import { Redirect, Route } from 'react-router-dom';

// Utils
import { loggedIn } from '../../utils/AuthService';

const PublicRoute = ({ component: Component, ...rest }) => (
  <Route {...rest} render={props => ( 
    loggedIn() ? 
    <Redirect to={{
        pathname: 'dashboard',
        state: { from: props.location }
    }} />
    : <Component {...props} />
  )} />
);

export default PublicRoute;  
