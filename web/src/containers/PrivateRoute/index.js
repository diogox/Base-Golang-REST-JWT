import React from 'react';  
import { Redirect, Route } from 'react-router-dom';

// Utils
import { loggedIn } from '../../utils/AuthService';

const PrivateRoute = ({ component: Component, ...rest }) => (
  <Route {...rest} render={props => ( 
    loggedIn() ? 
    <Component {...props} /> 
    : <Redirect to={{
        pathname: 'login',
        state: { from: props.location }
      }} /> 
  )} />
);

export default PrivateRoute;  
