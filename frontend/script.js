document.addEventListener('DOMContentLoaded', () => {
    const formContainer = document.getElementById('form-container');
    const listContainer = document.getElementById('list-container');
    const productForm = document.getElementById('product-form');
    const productList = document.getElementById('product-list');
    const newProductButton = document.getElementById('new-product-button');
  
    productForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const newProduct = {
        name: productForm['product-name'].value,
        description: productForm['product-description'].value,
        price: parseFloat(productForm['product-price'].value),
        available: productForm['product-available'].value
      };
  
      await fetch('/api/products', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newProduct)
      });
  
      productForm.reset();
      showProductList();
    });
  
    newProductButton.addEventListener('click', () => {
      formContainer.style.display = 'block';
      listContainer.style.display = 'none';
    });
  
    const showProductList = async () => {
      formContainer.style.display = 'none';
      listContainer.style.display = 'block';
  
      const response = await fetch('/api/products');
      const products = await response.json();
  
      productList.innerHTML = '';
      products.sort((a, b) => a.price - b.price).forEach(product => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${product.name}</td>
          <td>${product.price.toFixed(2)}</td>
        `;
        productList.appendChild(row);
      });
    };
  
    // Mostrar a lista de produtos ao carregar a página
    showProductList();
  });
  