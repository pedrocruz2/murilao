import React, { useState } from 'react';
import axios from 'axios';

const Auth = ({ setIsAuthenticated, setCurrentUser }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    profilePic: null
  });
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleFileChange = (e) => {
    setFormData({
      ...formData,
      profilePic: e.target.files[0]
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    
    try {
      if (isLogin) {
        const response = await axios.post('/api/auth/login', {
          email: formData.email,
          password: formData.password
        });
        
        localStorage.setItem('token', response.data.token);
        setIsAuthenticated(true);
        setCurrentUser(response.data.user);
      } else {
        const data = new FormData();
        data.append('name', formData.name);
        data.append('email', formData.email);
        data.append('password', formData.password);
        if (formData.profilePic) {
          data.append('profilePic', formData.profilePic);
        }
        
        await axios.post('/api/auth/register', data);
        setIsLogin(true);
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Ocorreu um erro');
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-form">
        <h2>{isLogin ? 'Login' : 'Cadastro'}</h2>
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          {!isLogin && (
            <div className="form-group">
              <label>Nome</label>
              <input 
                type="text" 
                name="name" 
                value={formData.name} 
                onChange={handleChange} 
                required 
              />
            </div>
          )}
          
          <div className="form-group">
            <label>Email</label>
            <input 
              type="email" 
              name="email" 
              value={formData.email} 
              onChange={handleChange} 
              required 
            />
          </div>
          
          <div className="form-group">
            <label>Senha</label>
            <input 
              type="password" 
              name="password" 
              value={formData.password} 
              onChange={handleChange} 
              required 
            />
          </div>
          
          {!isLogin && (
            <div className="form-group">
              <label>Foto de Perfil</label>
              <input 
                type="file" 
                name="profilePic" 
                onChange={handleFileChange} 
              />
            </div>
          )}
          
          <button type="submit">{isLogin ? 'Entrar' : 'Cadastrar'}</button>
        </form>
        
        <div className="auth-toggle">
          {isLogin ? (
            <p>Não tem uma conta? <button onClick={() => setIsLogin(false)}>Cadastre-se</button></p>
          ) : (
            <p>Já tem uma conta? <button onClick={() => setIsLogin(true)}>Login</button></p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Auth;