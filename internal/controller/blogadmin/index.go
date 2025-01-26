package blogadmin

import (
	"fmt"
	"go-blog-admin/internal/config"
	controller "go-blog-admin/internal/controller"
	"go-blog-admin/internal/service"
	"time"

	"go-blog-admin/internal/i18n"
	"go-blog-admin/internal/mvc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BlogAdminIndexController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET  bool
	IsPOST bool

	webCtxt echo.Context // webCtxt

	DTO struct {
		Input struct {
		}
		Meta struct {
			IsFragment bool `json:"-"`
		}
		Output struct {
			mvc.ModelBaseDTO
			LangCode  string
			AppConfig struct {
				AppTitle string `json:"app_title,omitempty"`
				TmTitle  string `json:"tm_title,omitempty"`
			}
			Title     string
			LangWords map[string]string
		}
	}
}

func (x *BlogAdminIndexController) Handler() error {

	err := x.createDTO()
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

func NewBlogAdminIndexController(appService service.AppService, c echo.Context) *BlogAdminIndexController {

	return &BlogAdminIndexController{
		appService: appService,
		appConfig:  appService.Config(),
		userLang:   controller.UserLang(c, appService),
		IsGET:      controller.IsGET(c),
		IsPOST:     controller.IsPOST(c),
		webCtxt:    c,
	}
}

func (x *BlogAdminIndexController) validateFields() {

}

func (x *BlogAdminIndexController) createDTO() error {

	dto := &x.DTO
	c := x.webCtxt

	if err := c.Bind(dto); err != nil {
		return err
	}

	x.validateFields()

	return nil
}

func (x *BlogAdminIndexController) handleDTO() error {

	dto := &x.DTO
	// input := &dto.Input
	output := &dto.Output
	// meta := &dto.Meta
	// c := x.webCtxt

	userLang := x.userLang
	output.LangCode = userLang.LangCode()
	output.Title = userLang.Lang("Blog admin") // TODO /*Lang*/
	output.LangWords = userLang.LangWords()

	cfg := &output.AppConfig

	cfg.AppTitle = x.appConfig.Title
	cfg.TmTitle = fmt.Sprintf("%s © %d", x.appConfig.Title, time.Now().Year())

	return nil
}

func (x *BlogAdminIndexController) responseDTOAsMvc() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	appConfig := x.appConfig
	lang := x.userLang
	c := x.webCtxt

	data, err := mvc.NewModelWrap(c, output, meta.IsFragment, "Blog admin" /*Lang*/, appConfig, lang)
	if err != nil {
		return err
	}

	err = c.Render(http.StatusOK, "index.html", data)

	if err != nil {
		return err
	}

	return nil
}

func (x *BlogAdminIndexController) responseDTO() (err error) {

	return x.responseDTOAsMvc()

}
