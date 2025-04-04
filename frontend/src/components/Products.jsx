import React, { useState, useEffect } from 'react';
import axios from 'axios';

const Products = () => {
  const [products, setProducts] = useState([]);
  const [editingProduct, setEditingProduct] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
    quantity: '',
    image: null
  });
  const [error, setError] = useState(null);

  const fetchProducts = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get('/api/products', {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      setProducts(response.data);
    } catch (err) {
      setError('Erro ao carregar produtos');
    }
  };

  useEffect(() => {
    fetchProducts();
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
      image: e.target.files[0]
    });
  };

  const handleEdit = (product) => {
    setEditingProduct(product);
    setFormData({
      name: product.name,
      description: product.description,
      price: product.price,
      quantity: product.quantity,
      image: null
    });
  };

  const handleDelete = async (id) => {
    if (window.confirm('Deseja realmente excluir este produto?')) {
      try {
        const token = localStorage.getItem('token');
        await axios.delete(`/api/products/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        fetchProducts();
      } catch (err) {
        setError('Erro ao excluir produto');
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
      data.append('description', formData.description);
      data.append('price', formData.price);
      data.append('quantity', formData.quantity);
      
      if (formData.image) {
        data.append('image', formData.image);
      }
      
      if (editingProduct) {
        await axios.put(`/api/products/${editingProduct.id}`, data, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
      } else {
        await axios.post('/api/products', data, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
      }
      
      fetchProducts();
      setEditingProduct(null);
      setFormData({ name: '', description: '', price: '', quantity: '', image: null });
    } catch (err) {
      setError(err.response?.data?.message || 'Ocorreu um erro');
    }
  };

  return (
    <div className="products-container">
      <h2>Gerenciamento de Produtos</h2>
      {error && <div className="error-message">{error}</div>}
      
      <div className="products-form">
        <h3>{editingProduct ? 'Editar Produto' : 'Adicionar Produto'}</h3>
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
            <label>Descrição</label>
            <textarea 
              name="description" 
              value={formData.description} 
              onChange={handleChange} 
              required 
            />
          </div>
          
          <div className="form-group">
            <label>Preço</label>
            <input 
              type="number" 
              step="0.01" 
              name="price" 
              value={formData.price} 
              onChange={handleChange} 
              required 
            />
          </div>
          
          <div className="form-group">
            <label>Quantidade</label>
            <input 
              type="number" 
              name="quantity" 
              value={formData.quantity} 
              onChange={handleChange} 
              required 
            />
          </div>
          
          <div className="form-group">
            <label>Imagem</label>
            <input 
              type="file" 
              name="image" 
              onChange={handleFileChange} 
              required={!editingProduct} 
            />
          </div>
          
          <div className="form-buttons">
            <button type="submit">{editingProduct ? 'Atualizar' : 'Adicionar'}</button>
            {editingProduct && (
              <button type="button" onClick={() => {
                setEditingProduct(null);
                setFormData({ name: '', description: '', price: '', quantity: '', image: null });
              }}>Cancelar</button>
            )}
          </div>
        </form>
      </div>
      
      <div className="products-list">
        <h3>Produtos Cadastrados</h3>
        <div className="products-grid">
          {products.map((product) => (
            <div key={product.id} className="product-card">
              <div className="product-image">
                {product.image ? (
                  <img src={`/uploads/${product.image}`} alt={product.name} />
                ) : (
                  <div className="no-image">Sem imagem</div>
                )}
              </div>
              <div className="product-info">
                <h4>{product.name}</h4>
                <p className="product-description">{product.description}</p>
                <p className="product-price">R$ {product.price.toFixed(2)}</p>
                <p className="product-quantity">Quantidade: {product.quantity}</p>
                <div className="product-actions">
                  <button onClick={() => handleEdit(product)}>Editar</button>
                  <button onClick={() => handleDelete(product.id)}>Excluir</button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Products;