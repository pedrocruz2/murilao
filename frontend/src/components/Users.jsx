import React, { useState, useEffect } from 'react';
import axios from 'axios';

const Users = () => {
  const [users, setUsers] = useState([]);
  const [editingUser, setEditingUser] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    profilePic: null
  });
  const [error, setError] = useState(null);

  const fetchUsers = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get('/api/users', {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      setUsers(response.data);
    } catch (err) {
      setError('Erro ao carregar usuários');
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

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

  const handleEdit = (user) => {
    setEditingUser(user);
    setFormData({
      name: user.name,
      email: user.email,
      password: '',
      profilePic: null
    });
  };

  const handleDelete = async (id) => {
    if (window.confirm('Deseja realmente excluir este usuário?')) {
      try {
        const token = localStorage.getItem('token');
        await axios.delete(`/api/users/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        fetchUsers();
      } catch (err) {
        setError('Erro ao excluir usuário');
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    
    try {
      const token = localStorage.getItem('token');
      const data = new FormData();
      data.append('name', formData.name);
      data.append('email', formData.email);
      
      if (formData.password) {
        data.append('password', formData.password);
      }
      
      if (formData.profilePic) {
        data.append('profilePic', formData.profilePic);
      }
      
      if (editingUser) {
        await axios.put(`/api/users/${editingUser.id}`, data, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
      } else {
        await axios.post('/api/users', data, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
      }
      
      fetchUsers();
      setEditingUser(null);
      setFormData({ name: '', email: '', password: '', profilePic: null });
    } catch (err) {
      setError(err.response?.data?.message || 'Ocorreu um erro');
    }
  };

  return (
    <div className="users-container">
      <h2>Gerenciamento de Usuários</h2>
      {error && <div className="error-message">{error}</div>}
      
      <div className="users-form">
        <h3>{editingUser ? 'Editar Usuário' : 'Adicionar Usuário'}</h3>
        <form onSubmit={handleSubmit}>
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
            <label>Senha {editingUser && '(deixe em branco para manter a atual)'}</label>
            <input 
              type="password" 
              name="password" 
              value={formData.password} 
              onChange={handleChange} 
              required={!editingUser} 
            />
          </div>
          
          <div className="form-group">
            <label>Foto de Perfil</label>
            <input 
              type="file" 
              name="profilePic" 
              onChange={handleFileChange} 
            />
          </div>
          
          <div className="form-buttons">
            <button type="submit">{editingUser ? 'Atualizar' : 'Adicionar'}</button>
            {editingUser && (
              <button type="button" onClick={() => {
                setEditingUser(null);
                setFormData({ name: '', email: '', password: '', profilePic: null });
              }}>Cancelar</button>
            )}
          </div>
        </form>
      </div>
      
      <div className="users-list">
        <h3>Usuários Cadastrados</h3>
        <table>
          <thead>
            <tr>
              <th>Foto</th>
              <th>Nome</th>
              <th>Email</th>
              <th>Ações</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>
                  {user.profilePic ? (
                    <img src={`/uploads/${user.profilePic}`} alt={user.name} width="50" height="50" />
                  ) : (
                    <div className="no-image">Sem foto</div>
                  )}
                </td>
                <td>{user.name}</td>
                <td>{user.email}</td>
                <td>
                  <button onClick={() => handleEdit(user)}>Editar</button>
                  <button onClick={() => handleDelete(user.id)}>Excluir</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Users;