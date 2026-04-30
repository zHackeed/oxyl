package member

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/api/internal/models/responses"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*ListMembers)(nil)

type ListMembers struct {
	company *service.CompanyService
}

func NewListMemberController(companyService *service.CompanyService) *ListMembers {
	return &ListMembers{
		company: companyService,
	}
}

func (l ListMembers) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (l ListMembers) GetPath() string {
	return "/company/:id/member"
}

func (l ListMembers) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (l ListMembers) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(l.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}
	members, err := l.company.GetMembers(ctx, request.CompanyId)

	if err != nil {
		if errors.Is(err, models.ErrPermissionDenied) {
			return fiber.ErrForbidden
		}
		slog.Error("failed to get members of the company", err)
		return fiber.ErrInternalServerError
	}

	wrappedUsers := make([]*responses.CompanyUserValueWrapper, 0)

	for _, member := range members {
		wrappedMember := new(responses.CompanyUserValueWrapper)

		wrappedMember.User = member.User
		wrappedMember.Permissions = member.Permissions.StringifiedPermissions()
		wrappedMember.CreatedAt = member.CreatedAt

		wrappedUsers = append(wrappedUsers, wrappedMember)
	}

	slog.Info(request.CompanyId, wrappedUsers)

	return ctx.JSON(wrappedUsers)
}
