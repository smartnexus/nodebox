-- +goose Up
-- Create changelog entries table for tracking item modification history.
CREATE TABLE "changelogs" (
    "id"          uuid         NOT NULL,
    "created_at"  timestamptz  NOT NULL,
    "updated_at"  timestamptz  NOT NULL,
    "summary"     varchar(500) NOT NULL,
    "entity_id"   uuid         NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "changelogs_entities_changelog_entries"
        FOREIGN KEY ("entity_id") REFERENCES "entities" ("id") ON DELETE CASCADE
);

CREATE INDEX "changelog_entity_id" ON "changelogs" ("entity_id");

-- +goose Down
DROP INDEX IF EXISTS "changelog_entity_id";
DROP TABLE IF EXISTS "changelogs";
