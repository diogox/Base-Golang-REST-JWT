import React from 'react';  
import { Redirect, Route } from 'react-router-dom';
import NavBar from '../../components/NavBar';

// Utils
import { loggedIn } from '../../utils/AuthService';

const PrivateRoute = ({ component: Component, ...rest }) => (
  <Route {...rest} render={props => ( 
    loggedIn() ?

      // Return the component
      <React.Fragment>
        <NavBar isLoggedIn={true} />
        <Component {...props} /> 
      </React.Fragment>

      // Redirect to login
    : <Redirect to={{
        pathname: 'login',
        state: { from: props.location }
      }} /> 
  )} />
);

export default PrivateRoute;  
