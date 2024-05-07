ALTER TABLE "images"
ADD CONSTRAINT "unique_original_src" UNIQUE ("original_src");
ALTER TABLE "images"
ADD CONSTRAINT "unique_simple_src" UNIQUE ("simple_src");
ALTER TABLE "images"
ADD CONSTRAINT "unique_original_filename" UNIQUE ("original_filename");
ALTER TABLE "images"
ADD CONSTRAINT "unique_simple_filename" UNIQUE ("simple_filename");