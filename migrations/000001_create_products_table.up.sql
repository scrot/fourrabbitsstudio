CREATE TABLE IF NOT EXISTS defaultdb.products 
(
  product_id INT PRIMARY KEY,
  product_link UUID UNIQUE NOT NULL,
  download_link TEXT UNIQUE NOT NULL
);
