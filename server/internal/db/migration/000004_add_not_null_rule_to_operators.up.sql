ALTER TABLE "operators"
ADD CONSTRAINT name_not_empty CHECK (name <> '');
ALTER TABLE "operators"
ADD CONSTRAINT hashed_password_not_empty CHECK (hashed_password <> '');
ALTER TABLE "operators"
ADD CONSTRAINT email_not_empty CHECK (email <> '');