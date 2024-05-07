ALTER TABLE child_categories
ADD COLUMN image_id bigint,
  ADD CONSTRAINT fk_child_categories_images FOREIGN KEY (image_id) REFERENCES images(id);