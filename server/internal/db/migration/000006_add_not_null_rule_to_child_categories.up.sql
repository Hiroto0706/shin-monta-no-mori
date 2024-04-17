ALTER TABLE "child_categories"
ADD CONSTRAINT name_not_empty CHECK (name <> '');