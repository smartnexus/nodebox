package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

// HandleChangelogGet godoc
//
//	@Summary	Get Changelog
//	@Tags		Item Changelog
//	@Produce	json
//	@Param		id	path		string	true	"Item ID"
//	@Success	200	{array}		repo.ChangelogEntry
//	@Router		/v1/entities/{id}/changelog [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleChangelogGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) ([]repo.ChangelogEntry, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Changelog.GetByEntityID(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}

// HandleChangelogCreate godoc
//
//	@Summary	Create Changelog Entry
//	@Tags		Item Changelog
//	@Produce	json
//	@Param		id		path		string						true	"Item ID"
//	@Param		payload	body		repo.ChangelogEntryCreate	true	"Entry Data"
//	@Success	201		{object}	repo.ChangelogEntry
//	@Router		/v1/entities/{id}/changelog [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleChangelogCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, itemID uuid.UUID, body repo.ChangelogEntryCreate) (repo.ChangelogEntry, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Changelog.Create(auth, itemID, body)
	}

	return adapters.ActionID("id", fn, http.StatusCreated)
}
