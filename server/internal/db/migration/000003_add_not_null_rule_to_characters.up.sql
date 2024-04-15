ALTER TABLE "characters"
ADD CONSTRAINT name_not_empty CHECK (name <> '');