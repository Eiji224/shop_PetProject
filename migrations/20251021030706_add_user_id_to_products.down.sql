ALTER TABLE products
 DROP FOREIGN KEY fk_products_user,
 DROP INDEX   idx_products_user_id,
 DROP COLUMN  user_id;