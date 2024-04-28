ALTER TABLE "images"
ADD COLUMN "original_filename" varchar NOT NULL DEFAULT '';
ALTER TABLE "images"
ADD COLUMN "simple_filename" varchar;