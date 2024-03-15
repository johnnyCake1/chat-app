import { Route, Routes } from 'react-router-dom';
import RegistrationPage from './components/authentication/RegistrationPage'
import LoginPage from './components/authentication/LoginPage';
import MainChatPage from './components/MainChatPage';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  console.log("App loaded");
  return (
    <div className="app">
      <Routes>
        <Route exact path="/" element={<ProtectedRoute><MainChatPage /></ProtectedRoute> } />
        <Route exact path="/register" element={<RegistrationPage /> } />
        <Route exact path="/login" element={<LoginPage /> } />
      </Routes>
    </div>
  );
}

export default App;
