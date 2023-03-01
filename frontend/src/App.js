import * as React from 'react';
import { Route, Routes, BrowserRouter } from "react-router-dom";

import Home from './pages/Home';
import Transaction from './pages/Transaction';
import Upload from './pages/Upload';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route index path="/" element={<Home />} />
        <Route path="/upload" element={<Upload />} />
        <Route path="/transaction" element={<Transaction />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
