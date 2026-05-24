package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/schema/mixins"
)

// ChangelogTag holds the schema definition for a group-scoped changelog tag/category.
type ChangelogTag struct {
	ent.Schema
}

func (ChangelogTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		GroupMixin{ref: "changelog_tags"},
	}
}

func (ChangelogTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.String("color").
			MaxLen(20).
			Optional().
			Default(""),
	}
}

func (ChangelogTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("changelogs", Changelog.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.SetNull,
			}),
	}
}
