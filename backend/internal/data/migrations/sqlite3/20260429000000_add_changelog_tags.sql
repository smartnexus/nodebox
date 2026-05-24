-- +goose Up
-- Create the changelog_tags table (group-scoped categories for changelog entries).
CREATE TABLE `changelog_tags` (
    `id`                   uuid NOT NULL,
    `created_at`           datetime NOT NULL,
    `updated_at`           datetime NOT NULL,
    `name`                 varchar(100) NOT NULL,
    `color`                varchar(20) NULL DEFAULT '',
    `group_changelog_tags` uuid NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `changelog_tags_groups_changelog_tags`
        FOREIGN KEY (`group_changelog_tags`) REFERENCES `groups` (`id`) ON DELETE CASCADE
);

-- Add optional tag_id FK to changelogs table.
ALTER TABLE `changelogs` ADD COLUMN `tag_id` uuid NULL;

-- +goose Down
ALTER TABLE `changelogs` DROP COLUMN `tag_id`;
DROP TABLE `changelog_tags`;
