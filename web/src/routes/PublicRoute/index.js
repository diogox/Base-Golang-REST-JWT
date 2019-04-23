import React from 'react';  
import { Redirect, Route } from 'react-router-dom';
import NavBar from '../../components/NavBar';

// Utils
import { loggedIn } from '../../utils/AuthService';

const PublicRoute = ({ component: Component, ...rest }) => (
  <Route {...rest} render={props => ( 
    loggedIn() ?

      // Redirect to login
      <Redirect to={{
          pathname: 'dashboard',
          state: { from: props.location }
      }} />
      
    // Return the component
    : <React.Fragment>
        <NavBar isLoggedIn={false} />
        <Component {...props} /> 
      </React.Fragment>
  )} />
);

export default PublicRoute;  
