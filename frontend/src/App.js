import { Route, Routes } from 'react-router-dom';
import Chat from './components/Chat'
import RegistrationPage from './components/authentication/RegistrationPage'
import LoginPage from './components/authentication/LoginPage';
import MainChatPage from './components/MainChatPage';

function App() {
  return (
    <div className="app">
      <Routes>
        <Route exact path="/chat" element={<MainChatPage /> } />
        <Route exact path="/register" element={<RegistrationPage /> } />
        <Route exact path="/login" element={<LoginPage /> } />
      </Routes>
    </div>
  );
}

export default App;
