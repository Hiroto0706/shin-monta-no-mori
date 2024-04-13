-- 外部キー制約の削除
ALTER TABLE
  "image_characters_relations" DROP CONSTRAINT IF EXISTS image_characters_relations_character_id_fkey;

ALTER TABLE
  "image_characters_relations" DROP CONSTRAINT IF EXISTS image_characters_relations_image_id_fkey;

ALTER TABLE
  "image_parent_categories_relations" DROP CONSTRAINT IF EXISTS image_parent_categories_relations_parent_category_id_fkey;

ALTER TABLE
  "image_parent_categories_relations" DROP CONSTRAINT IF EXISTS image_parent_categories_relations_image_id_fkey;

ALTER TABLE
  "child_categories" DROP CONSTRAINT IF EXISTS child_categories_parent_id_fkey;

ALTER TABLE
  "sessions" DROP CONSTRAINT IF EXISTS sessions_name_fkey;

-- インデックスの削除
DROP INDEX IF EXISTS "parent_categories_idx_src";

DROP INDEX IF EXISTS "characters_idx_src";

DROP INDEX IF EXISTS ":images_idx_simple_src";

DROP INDEX IF EXISTS ":images_idx_original_src";

DROP INDEX IF EXISTS "image_characters_relations_image_id_character_id_idx";

DROP INDEX IF EXISTS "image_parent_categories_relations_image_id_parent_category_id_idx";

-- テーブルの削除
DROP TABLE IF EXISTS "image_characters_relations";

DROP TABLE IF EXISTS "image_parent_categories_relations";

DROP TABLE IF EXISTS "child_categories";

DROP TABLE IF EXISTS "parent_categories";

DROP TABLE IF EXISTS "characters";

DROP TABLE IF EXISTS "images";

DROP TABLE IF EXISTS "sessions";

DROP TABLE IF EXISTS "operators";