-- name: ProductPageSize :many
SELECT 
    p.id,          
    p.sku,         
    p.name,        
    p.description, 
    c.name AS category_name,
    m.name AS manufacturer_name,
    p.weight,
    p.unit,   
    p.price,       
    p.baseprice,   
    p.issueyear   
FROM products AS p
INNER JOIN categories AS c ON p.category_id = c.id
INNER JOIN manufacturers AS m ON p.manufacturer_id = m.id
OFFSET $1
LIMIT $2;

-- name: ProductTotal :one
SELECT COUNT(*) AS total
FROM products;

-- name: ProductPageSizeCategory :many
SELECT 
    p.id,          
    p.sku,         
    p.name,        
    p.description, 
    c.name AS category_name,
    m.name AS manufacturer_name,
    p.weight,
    p.unit,   
    p.price,       
    p.baseprice,   
    p.issueyear   
FROM products AS p
INNER JOIN categories AS c ON p.category_id = c.id
INNER JOIN manufacturers AS m ON p.manufacturer_id = m.id
WHERE c.id = $3
OFFSET $1
LIMIT $2;

-- name: ProductTotalCategory :one
SELECT COUNT(category_id) AS total
FROM products
WHERE category_id = $1;