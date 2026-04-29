package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

// HandleChangelogTagsGetAll godoc
//
//	@Summary	Get All Changelog Tags
//	@Tags		Changelog Tags
//	@Produce	json
//	@Success	200	{array}		repo.ChangelogTagOut
//	@Router		/v1/changelog-tags [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleChangelogTagsGetAll() errchain.HandlerFunc {
	fn := func(r *http.Request) ([]repo.ChangelogTagOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.ChangelogTags.GetAll(r.Context(), auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandleChangelogTagCreate godoc
//
//	@Summary	Create Changelog Tag
//	@Tags		Changelog Tags
//	@Produce	json
//	@Param		payload	body		repo.ChangelogTagCreate	true	"Tag Data"
//	@Success	201		{object}	repo.ChangelogTagOut
//	@Router		/v1/changelog-tags [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleChangelogTagCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, body repo.ChangelogTagCreate) (repo.ChangelogTagOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.ChangelogTags.Create(r.Context(), auth.GID, body)
	}

	return adapters.Action(fn, http.StatusCreated)
}

// HandleChangelogTagDelete godoc
//
//	@Summary	Delete Changelog Tag
//	@Tags		Changelog Tags
//	@Produce	json
//	@Param		id	path	string	true	"Tag ID"
//	@Success	204
//	@Router		/v1/changelog-tags/{id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleChangelogTagDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		err := ctrl.repo.ChangelogTags.Delete(r.Context(), auth.GID, ID)
		return nil, err
	}

	return adapters.CommandID("id", fn, http.StatusNoContent)
}
