import React, { useState, useEffect } from 'react';
import LanguageList from './components/LanguageList';
import './App.css';

function App() {
  const [languages, setLanguages] = useState([]);

  useEffect(() => {
    fetch('./data.json')
      .then(response => response.json())
      .then(data => setLanguages(data));
  }, []);

  return (
    <div className="App">
      <h1>Top Open-Source Programming Languages</h1>
      <LanguageList languages={languages} />
    </div>
  );
}

export default App;
