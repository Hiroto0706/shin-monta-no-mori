ALTER TABLE "images"
ADD CONSTRAINT title_not_empty CHECK (title <> '');