package member

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/api/internal/models/responses"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
)

var _ apiModel.Registrable = (*SelfPermissionController)(nil)

type SelfPermissionController struct {
	companyService *service.CompanyService
}

func NewSelfPermissionController(companyService *service.CompanyService) *SelfPermissionController {
	return &SelfPermissionController{
		companyService: companyService,
	}
}

func (s SelfPermissionController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodGet
}

func (s SelfPermissionController) GetPath() string {
	return "/company/:id/permissions/self"
}

func (s SelfPermissionController) RequestRequirements() *apiModel.RequestRequirements {
	return apiModel.NewRequestRequirements(apiModel.URIData, requests.CompanyIdUri{})
}

func (s SelfPermissionController) Handle(ctx fiber.Ctx) error {
	request, ok := ctx.Locals(s.RequestRequirements().GetValidationType()).(*requests.CompanyIdUri)
	if !ok {
		return fiber.ErrInternalServerError
	}

	member, err := s.companyService.GetMember(ctx, request.CompanyId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return fiber.ErrUnauthorized
		}
		slog.Error("failed to get company membership", slog.String("companyId", request.CompanyId), slog.String("error", err.Error()))
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(responses.CompanyMemberPermissionResponse{
		User:        member.UserID,
		Permissions: member.Permission.StringifiedPermissions(),
	})
}
