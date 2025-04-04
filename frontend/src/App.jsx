import React, { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Auth from './components/Auth';
import Users from './components/Users';
import Products from './components/Products';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [currentUser, setCurrentUser] = useState(null);
  
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      fetch('/api/auth/verify', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      .then(res => {
        if (res.ok) {
          setIsAuthenticated(true);
          return res.json();
        } else {
          localStorage.removeItem('token');
          setIsAuthenticated(false);
          throw new Error('Token inválido');
        }
      })
      .then(data => {
        setCurrentUser(data.user);
      })
      .catch(err => {
        console.error(err);
      });
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsAuthenticated(false);
    setCurrentUser(null);
  };

  const ProtectedRoute = ({ children }) => {
    if (!isAuthenticated) {
      return <Navigate to="/login" />;
    }
    return children;
  };

  return (
    <BrowserRouter>
      <div className="app">
        {isAuthenticated && (
          <nav className="navbar">
            <span>Olá, {currentUser?.name}</span>
            <button onClick={handleLogout}>Logout</button>
          </nav>
        )}
        <Routes>
          <Route path="/login" element={
            isAuthenticated ? <Navigate to="/" /> : <Auth setIsAuthenticated={setIsAuthenticated} setCurrentUser={setCurrentUser} />
          } />
          <Route path="/" element={
            <ProtectedRoute>
              <div className="container">
                <div className="sidebar">
                  <ul>
                    <li><a href="/users">Usuários</a></li>
                    <li><a href="/products">Produtos</a></li>
                  </ul>
                </div>
                <div className="content">
                  <h1>Dashboard</h1>
                  <p>Bem-vindo ao sistema de gerenciamento</p>
                </div>
              </div>
            </ProtectedRoute>
          } />
          <Route path="/users" element={
            <ProtectedRoute>
              <Users />
            </ProtectedRoute>
          } />
          <Route path="/products" element={
            <ProtectedRoute>
              <Products />
            </ProtectedRoute>
          } />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;