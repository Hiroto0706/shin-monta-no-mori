CREATE TABLE "operators" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "images" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "original_src" varchar NOT NULL,
  "simple_src" varchar,
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "characters" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "src" varchar NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "parent_categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "src" varchar NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "child_categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "parent_id" bigint NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "image_parent_categories_relations" (
  "id" bigserial PRIMARY KEY,
  "image_id" bigint NOT NULL,
  "parent_category_id" bigint NOT NULL
);

CREATE TABLE "image_characters_relations" (
  "id" bigserial PRIMARY KEY,
  "image_id" bigint NOT NULL,
  "character_id" bigint NOT NULL
);

CREATE INDEX ":images_idx_original_src" ON "images" ("original_src");

CREATE INDEX ":images_idx_simple_src" ON "images" ("simple_src");

CREATE INDEX "characters_idx_src" ON "characters" ("src");

CREATE INDEX "parent_categories_idx_src" ON "parent_categories" ("src");

CREATE UNIQUE INDEX ON "image_parent_categories_relations" ("image_id", "parent_category_id");

CREATE UNIQUE INDEX ON "image_characters_relations" ("image_id", "character_id");

COMMENT ON COLUMN "images"."original_src" IS '文字ありの画像';

COMMENT ON COLUMN "images"."simple_src" IS '文字無しの画像.オリジナルが文字無しの時もあるので、nullableにする.';

ALTER TABLE "sessions" ADD FOREIGN KEY ("name") REFERENCES "operators" ("name");

ALTER TABLE "child_categories" ADD FOREIGN KEY ("parent_id") REFERENCES "parent_categories" ("id");

ALTER TABLE "image_parent_categories_relations" ADD FOREIGN KEY ("image_id") REFERENCES "images" ("id");

ALTER TABLE "image_parent_categories_relations" ADD FOREIGN KEY ("parent_category_id") REFERENCES "parent_categories" ("id");

ALTER TABLE "image_characters_relations" ADD FOREIGN KEY ("image_id") REFERENCES "images" ("id");

ALTER TABLE "image_characters_relations" ADD FOREIGN KEY ("character_id") REFERENCES "characters" ("id");
