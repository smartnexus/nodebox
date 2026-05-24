-- +goose Up
-- Create the changelog_tags table (group-scoped categories for changelog entries).
CREATE TABLE "changelog_tags" (
    "id"                   uuid NOT NULL,
    "created_at"           timestamptz NOT NULL,
    "updated_at"           timestamptz NOT NULL,
    "name"                 varchar(100) NOT NULL,
    "color"                varchar(20) DEFAULT '',
    "group_changelog_tags" uuid NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "changelog_tags_groups_changelog_tags"
        FOREIGN KEY ("group_changelog_tags") REFERENCES "groups" ("id") ON DELETE CASCADE
);

-- Add optional tag_id FK to changelogs table.
ALTER TABLE "changelogs" ADD COLUMN "tag_id" uuid;

ALTER TABLE "changelogs"
    ADD CONSTRAINT "changelogs_changelog_tags_changelogs"
    FOREIGN KEY ("tag_id") REFERENCES "changelog_tags" ("id") ON DELETE SET NULL;

-- +goose Down
ALTER TABLE "changelogs" DROP CONSTRAINT IF EXISTS "changelogs_changelog_tags_changelogs";
ALTER TABLE "changelogs" DROP COLUMN IF EXISTS "tag_id";
DROP TABLE IF EXISTS "changelog_tags";
