import React, { useState } from 'react';
import './App.css';
import {API_URL} from "./constants";

function App() {
  const [serverTime, setServerTime] = useState('');

  const fetchServerTime = async () => {
    const response = await fetch(`${API_URL}/api/v1/time`);
    const data = await response.json();
    setServerTime(data.serverTime);
  };

  return (
      <div className="App">
        <header className="App-header">
          <p>Server Time: {serverTime}</p>
          <button onClick={fetchServerTime}>Fetch Server Tim`e!</button>
        </header>
      </div>
  );
}

export default App;
