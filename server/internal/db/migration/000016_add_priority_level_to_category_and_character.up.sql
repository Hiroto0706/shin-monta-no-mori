ALTER TABLE "characters"
ADD COLUMN "priority_level" SMALLINT NOT NULL DEFAULT 2;
ALTER TABLE "parent_categories"
ADD COLUMN "priority_level" SMALLINT NOT NULL DEFAULT 2;
ALTER TABLE "child_categories"
ADD COLUMN "priority_level" SMALLINT NOT NULL DEFAULT 2;