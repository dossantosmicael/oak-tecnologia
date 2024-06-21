const express = require('express')
const bodyParser = require('body-parser')
const fs = require('fs')
const { stringify } = require('querystring')

const app = express()
const PORT = 3000

app.use(bodyParser.json())
app.use(express.static('frontend'))

app.get('/api/products', (req, res) => {
    const data = JSON.parse(fs.readFileSync('backend/data.json'))
    res.json(data.products)
})

app.post('/api/products', (req, res) => {
    const newProduct = req.body
    const data = JSON.parse(fs.readFileSync('backend/data.json'))
    data.products.push(newProduct)
    fs.writeFileSync('backend/data.json', JSON.stringify(data))
    res.status(201).json(newProduct)
})

app.listen(PORT, ()=> {
    console.log(`Server is running on http://localhost:${PORT}`);
})