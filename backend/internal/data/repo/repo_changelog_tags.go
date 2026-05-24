package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/changelogtag"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/group"
)

// ChangelogTagRepository provides data access methods for group-scoped changelog tags.
type ChangelogTagRepository struct {
	db *ent.Client
}

// ChangelogTagCreate is the input for creating a changelog tag.
type ChangelogTagCreate struct {
	Name  string `json:"name"  validate:"required"`
	Color string `json:"color"`
}

// ChangelogTagOut is the full representation of a changelog tag.
type ChangelogTagOut struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"createdAt"`
}

func mapChangelogTag(ct *ent.ChangelogTag) ChangelogTagOut {
	return ChangelogTagOut{
		ID:        ct.ID,
		Name:      ct.Name,
		Color:     ct.Color,
		CreatedAt: ct.CreatedAt,
	}
}

// GetAll returns all changelog tags for a group, ordered by name.
func (r *ChangelogTagRepository) GetAll(ctx context.Context, gid uuid.UUID) ([]ChangelogTagOut, error) {
	tags, err := r.db.ChangelogTag.Query().
		Where(changelogtag.HasGroupWith(group.ID(gid))).
		Order(changelogtag.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ChangelogTagOut, len(tags))
	for i, ct := range tags {
		out[i] = mapChangelogTag(ct)
	}
	return out, nil
}

// Create creates a new changelog tag for a group.
func (r *ChangelogTagRepository) Create(ctx context.Context, gid uuid.UUID, input ChangelogTagCreate) (ChangelogTagOut, error) {
	ct, err := r.db.ChangelogTag.Create().
		SetName(input.Name).
		SetColor(input.Color).
		SetGroupID(gid).
		Save(ctx)
	if err != nil {
		return ChangelogTagOut{}, err
	}
	return mapChangelogTag(ct), nil
}

// Delete removes a changelog tag by ID, verifying group ownership.
func (r *ChangelogTagRepository) Delete(ctx context.Context, gid, id uuid.UUID) error {
	_, err := r.db.ChangelogTag.Delete().
		Where(
			changelogtag.ID(id),
			changelogtag.HasGroupWith(group.ID(gid)),
		).Exec(ctx)
	return err
}
