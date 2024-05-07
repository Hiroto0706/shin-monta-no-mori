CREATE TABLE "image_child_categories_relations" (
  "id" bigserial PRIMARY KEY,
  "image_id" bigint NOT NULL,
  "child_category_id" bigint NOT NULL
);

CREATE UNIQUE INDEX ON "image_child_categories_relations" ("image_id", "child_category_id");

ALTER TABLE
  "image_child_categories_relations"
ADD
  FOREIGN KEY ("image_id") REFERENCES "images" ("id");

ALTER TABLE
  "image_child_categories_relations"
ADD
  FOREIGN KEY ("child_category_id") REFERENCES "child_categories" ("id");