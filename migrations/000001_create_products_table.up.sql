CREATE TABLE IF NOT EXISTS products 
(
  product_link TEXT PRIMARY KEY, 
  download_link TEXT UNIQUE NOT NULL
);
