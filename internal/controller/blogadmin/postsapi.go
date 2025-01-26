package blogadmin

import (
	"go-blog-admin/internal/config"
	controller "go-blog-admin/internal/controller"
	"go-blog-admin/internal/util/utilaccess"
	"go-blog-admin/internal/util/utilorm"
	"go-blog-admin/internal/util/utilpaging"

	"go-blog-admin/internal/i18n"
	"go-blog-admin/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

var blogPostOmitColumns = []string{
	"content_markdown",
	// "content_html",
}

var blogPostFieldOptions = utilorm.FieldOptions{
	"code":             utilorm.ShortStringType,
	"title":            utilorm.ShortStringType,
	"content_markdown": utilorm.ShortStringType,
	"is_published":     utilorm.BoolType,
	"is_listed":        utilorm.BoolType,
}

var blogPostSortOptions = utilorm.SortOptions{
	"":      "-id", // default
	"id":    "id ASC",
	"-id":   "id DESC",
	"code":  "code ASC",
	"-code": "code DESC",
}

type PostsDTO struct {
	Input struct {
		utilpaging.PagingInputDTO
	}
	Meta struct {
		Status int
	}
	Output struct {
		Message string `json:"message,omitempty"`
		utilpaging.PagingOutputDTO[service.BlogPost]
		Permissions utilaccess.PermissionsDTO `json:"permissions,omitempty"`
	}
}

type PostsAPIController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET bool

	webCtxt echo.Context // webCtxt

	userAccount *service.UserAccount
	DTO         PostsDTO
}

func (x *PostsAPIController) Handler() error {
	// TODO sign out force

	err := x.validateDTO()
	if err != nil {
		return err
	}

	err = x.handleDTO()
	if err != nil {
		return err
	}

	err = x.responseDTO()
	if err != nil {
		return err
	}

	return nil
}

// NewAccountController is constructor.
func NewPostsAPIController(appService service.AppService, c echo.Context) *PostsAPIController {

	appConfig := appService.Config()

	return &PostsAPIController{
		appService:  appService,
		appConfig:   appConfig,
		userLang:    controller.UserLang(c, appService),
		IsGET:       controller.IsGET(c),
		userAccount: controller.GetAccount(c),
		webCtxt:     c,
	}
}

func (x *PostsAPIController) validateDTO() error {

	dto := &x.DTO
	input := &dto.Input

	c := x.webCtxt

	if err := c.Bind(input); err != nil {
		return err
	}

	input.LoadFilters(c.QueryParams()) //

	return nil
}

func (x *PostsAPIController) handleDTO() error {

	dto := &x.DTO
	input := &dto.Input
	meta := &dto.Meta
	output := &dto.Output
	// userLang := x.userLang
	// c := x.webCtxt
	// isInputValid := output.IsModelValid()

	if x.IsGET {

		bs := x.appService.BlogAdmin()

		{
			bs.Posts().Permissions(x.userAccount, &output.Permissions)
		}

		if err := bs.Posts().Query(&input.PagingInputDTO, &output.PagingOutputDTO, blogPostOmitColumns, blogPostFieldOptions, blogPostSortOptions); err != nil {
			return err
		}

	} else {
		meta.Status = http.StatusMethodNotAllowed
		output.Message = "method action undef"
	}

	return nil
}
func (x *PostsAPIController) responseDTOAsAPI() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt

	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}

	return c.JSON(meta.Status, output)

}

func (x *PostsAPIController) responseDTO() (err error) {

	return x.responseDTOAsAPI()

}
