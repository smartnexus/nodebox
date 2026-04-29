package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/changelog"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/entity"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/group"
)

// ChangelogRepository provides data access methods for changelog entries
// associated with items (entities) in the database.
type ChangelogRepository struct {
	db *ent.Client
}

// ChangelogEntryCreate is the input type for creating a new changelog entry.
type ChangelogEntryCreate struct {
	Summary string     `json:"summary" validate:"required"`
	TagID   *uuid.UUID `json:"tagId,omitempty"`
}

// ChangelogTagSummary is the summary representation of a changelog tag.
type ChangelogTagSummary struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
}

// ChangelogEntry is the output type for a single changelog entry.
type ChangelogEntry struct {
	ID        uuid.UUID            `json:"id"`
	EntityID  uuid.UUID            `json:"entityId"`
	CreatedAt time.Time            `json:"createdAt"`
	Summary   string               `json:"summary"`
	Tag       *ChangelogTagSummary `json:"tag,omitempty"`
}

var (
	mapChangelogEntryErr  = mapTErrFunc(mapChangelogEntry)
	mapEachChangelogEntry = mapTEachFunc(mapChangelogEntry)
)

func mapChangelogEntry(entry *ent.Changelog) ChangelogEntry {
	out := ChangelogEntry{
		ID:        entry.ID,
		EntityID:  entry.EntityID,
		CreatedAt: entry.CreatedAt,
		Summary:   entry.Summary,
	}

	if entry.Edges.Tag != nil {
		out.Tag = &ChangelogTagSummary{
			ID:    entry.Edges.Tag.ID,
			Name:  entry.Edges.Tag.Name,
			Color: entry.Edges.Tag.Color,
		}
	}

	return out
}

// Create inserts a new changelog entry for the given entity.
func (r *ChangelogRepository) Create(ctx context.Context, entityID uuid.UUID, input ChangelogEntryCreate) (ChangelogEntry, error) {
	q := r.db.Changelog.Create().
		SetEntityID(entityID).
		SetSummary(input.Summary)

	if input.TagID != nil && *input.TagID != uuid.Nil {
		q = q.SetTagID(*input.TagID)
	}

	entry, err := q.Save(ctx)
	if err != nil {
		return ChangelogEntry{}, err
	}

	// Reload with tag edge so the response includes the tag.
	entry, err = r.db.Changelog.Query().
		Where(changelog.ID(entry.ID)).
		WithTag().
		Only(ctx)

	return mapChangelogEntryErr(entry, err)
}

// GetByEntityID retrieves all changelog entries for a given entity, scoped to the group for
// multi-tenancy. Entries are returned in descending order by creation time (newest first).
func (r *ChangelogRepository) GetByEntityID(ctx context.Context, groupID, entityID uuid.UUID) ([]ChangelogEntry, error) {
	entries, err := r.db.Changelog.Query().
		Where(
			changelog.EntityID(entityID),
			changelog.HasEntityWith(
				entity.HasGroupWith(group.IDEQ(groupID)),
			),
		).
		WithTag().
		Order(ent.Desc(changelog.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return mapEachChangelogEntry(entries), nil
}

