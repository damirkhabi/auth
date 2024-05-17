package converter

import (
	"github.com/arifullov/auth/internal/model"
	modelRepo "github.com/arifullov/auth/internal/repository/access/model"
)

func ToRolesFromRepo(routeAccesses []modelRepo.RouteAccesses) []model.Role {
	userRoles := make([]model.Role, 0, len(routeAccesses))
	for _, routeAccess := range routeAccesses {
		userRoles = append(userRoles, model.Role(routeAccess.Role))
	}
	return userRoles
}
