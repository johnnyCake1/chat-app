import { Route, Routes } from 'react-router-dom';
import RegistrationPage from './components/authentication/RegistrationPage'
import LoginPage from './components/authentication/LoginPage';
import MainChatPage from './components/MainChatPage';
import ProtectedRoute from './components/ProtectedRoute';
import { useEffect, useState } from 'react';
import { API_URL } from './constants';
import useLocalStorageState from './util/userLocalStorage';

function App() {
  const [user, setUser] = useState(null);
  const [token,] = useLocalStorageState('token');
  

  useEffect(() => {
    if (token) {
      const tokenParts = token.split('.');
      if (tokenParts.length === 3) {
        // Decode the payload (second part)
        const payload = JSON.parse(atob(tokenParts[1]));
        // Extract the user ID from the Issuer payload
        const userId = payload.iss;
        fetch(`${API_URL}/users/${userId}`, { headers: { Authorization: `Bearer ${token}` } })
        .then(response => {
          if (response.ok) {
            response.json().then(userData => {
              setUser(userData);
            })
          }
        })
      } else {
        console.error('Invalid JWT token format');
      }
    }
  }, [token]);

  console.log("App loaded");

  return (
    <div className="app">
      <Routes>
        <Route exact path="/" element={<ProtectedRoute><MainChatPage currentUser={user} /></ProtectedRoute> } />
        <Route exact path="/register" element={<RegistrationPage /> } />
        <Route exact path="/login" element={<LoginPage /> } />
      </Routes>
    </div>
  );
}

export default App;
