ALTER TABLE  products
ADD   COLUMN user_id BIGINT UNSIGNED NOT NULL,
ADD   INDEX  idx_products_user_id (user_id),
ADD   CONSTRAINT fk_products_user
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;