import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Register from './pages/Register';
import NoPage from './pages/NoPage';
import Home from './pages/Home';
import Vault from './pages/Vault';
import Login from './pages/Login';

function App() {
  return (
    <>
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/vault" element={<Vault />} />
          <Route path="*" element={<NoPage />} />
        </Routes>
      </Router>
    </>
  )
}

export default App
