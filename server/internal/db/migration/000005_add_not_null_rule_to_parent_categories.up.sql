ALTER TABLE "parent_categories"
ADD CONSTRAINT name_not_empty CHECK (name <> '');