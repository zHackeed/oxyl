package company

import (
	"github.com/gofiber/fiber/v3"
	bind "github.com/idan-fishman/fiber-bind"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/models/requests"
	"zhacked.me/oxyl/shared/pkg/service"
)

var _ apiModel.Registrable = (*AddMemberController)(nil)

type AddMemberController struct {
	companyService *service.CompanyService
}

func NewAddMemberController(companyService *service.CompanyService) *AddMemberController {
	return &AddMemberController{
		companyService: companyService,
	}
}

func (a *AddMemberController) GetMethod() apiModel.HttpMethod {
	return apiModel.MethodPost
}

func (a *AddMemberController) GetPath() string {
	return "/company/:id/member"
}

func (a *AddMemberController) GetRequestModel() interface{} {
	return requests.AddMemberRequest{}
}

func (a *AddMemberController) Handle(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	if companyID == "" {
		return fiber.ErrBadRequest
	}
	request, ok := ctx.Locals(bind.JSON).(*requests.AddMemberRequest)

	if !ok {
		return fiber.ErrInternalServerError
	}

	if err := a.companyService.AddUserToCompany(ctx.Context(), companyID, request.UserEmail, int(request.Permission)); err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member added"})
}
