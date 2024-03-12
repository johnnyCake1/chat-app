import React, { useEffect, useState } from 'react';
import { Navigate } from 'react-router-dom';
import { API_URL } from '../constants';

// PrivateRoute is a component wrapper and is responsible for redirecting to authentication page if the user is not authorised to see the children components.
const PrivateRoute = ({ children }) => {
  const [isTokenValid, setIsTokenValid] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    try {
        fetch(`${API_URL}/validateToken`, {
            method: 'POST',
            credentials: 'include'
        }).then(response => {
            if (response.ok) {
                setIsTokenValid(true);
            }
            setIsLoading(false);
        }).catch(err => {
            console.log("err", err);
        });
    } catch (error) {
        console.error('Error validating user token:', error);
    }
  });

  if (isLoading) {
    return (
      <div className="flex flex-col justify-center items-center pt-8">
        Loading...
      </div>
    );
  }

  if (!isTokenValid) {
    return <Navigate to="/login" />;
  }

  return children;
};

export default PrivateRoute;
