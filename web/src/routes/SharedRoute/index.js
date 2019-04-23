import React from 'react';  
import { Route } from 'react-router-dom';
import NavBar from '../../components/NavBar';

// Utils
import { loggedIn } from '../../utils/AuthService';

const SharedRoute = ({ component: Component, ...rest }) => (
  <Route {...rest} render={props => (

    // Return the component
    <React.Fragment>
      <NavBar isLoggedIn={loggedIn()} />
      <Component {...props} /> 
    </React.Fragment>
  )} />
);

export default SharedRoute;  
