package blogadmin

import (
	"go-blog-admin/internal/config"
	"go-blog-admin/internal/config/consts"
	controller "go-blog-admin/internal/controller"
	"go-blog-admin/internal/mvc"

	"strings"

	"go-blog-admin/internal/i18n"
	"go-blog-admin/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BlogPostDTO struct {
	service.BlogPost
}
type PostsEntityDTO struct {
	Input struct {
		Code string      `param:"code"`
		ID   int64       `param:"id"`
		Data BlogPostDTO `json:"data,omitempty"`
	}
	Meta struct {
		Status int
	}
	Output struct {
		mvc.ModelBaseDTO
		Data BlogPostDTO `json:"data,omitempty"`
	}
}
type PostsEntityAPIController struct {
	appService service.AppService
	appConfig  *config.AppConfig
	userLang   i18n.UserLang

	IsGET    bool
	IsPOST   bool
	IsPUT    bool
	IsDELETE bool

	webCtxt echo.Context // webCtxt

	DTO PostsEntityDTO
}

func (x *PostsEntityAPIController) Handler() error {
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
func NewPostsEntityAPIController(appService service.AppService, c echo.Context) *PostsEntityAPIController {

	appConfig := appService.Config()

	return &PostsEntityAPIController{
		appService: appService,
		appConfig:  appConfig,
		userLang:   controller.UserLang(c, appService),
		IsGET:      controller.IsGET(c),
		IsPOST:     controller.IsPOST(c),
		IsPUT:      controller.IsPUT(c),
		IsDELETE:   controller.IsDELETE(c),
		webCtxt:    c,
	}
}

func (x *PostsEntityAPIController) validateDTOFields() (err error) {

	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	meta := &dto.Meta
	srv := x.appService.BlogAdmin()

	if x.IsPOST || x.IsPUT {

		// validate input: add update

		{
			input.Data.Code = strings.TrimSpace(input.Data.Code)
			input.Data.Title = strings.TrimSpace(input.Data.Title)
			// // input.Data.ContentHTML = "" // reset
		}

		{
			v := output.NewModelValidatorStr(x.userLang, "code", "Code", input.Data.Code, consts.DefaultTextLength)
			v.Required()
		}

		{
			v := output.NewModelValidatorStr(x.userLang, "title", "Title", input.Data.Title, consts.DefaultTextLength)
			v.Required()
		}

		{
			_ = output.NewModelValidatorStr(x.userLang, "content_markdown", "Markdown", input.Data.ContentMarkdown, consts.LongTextLength)
		}

	}

	if !output.IsModelValid() {
		meta.Status = http.StatusUnprocessableEntity // 422 validation
		return nil
	}

	if x.IsDELETE || x.IsPUT {
		// exists: delete, update

		id, err := srv.Posts().ID(input.ID)
		if err != nil {
			return err
		}

		if id == 0 {
			meta.Status = http.StatusNotFound // 404
			return nil
		}

	}

	if x.IsPOST || x.IsPUT {
		// code dupl: add, update

		id, err := srv.Posts().Code(input.Data.Code)
		if err != nil {
			return err
		}

		if !(id == 0 || input.ID == id) {
			meta.Status = http.StatusConflict           // e.g., duplicate data 409
			output.AddError("code", "Duplicate entry.") // Lang
			return nil
		}

	}

	// if x.IsPOST || x.IsPUT {

	// 	if input.Data.ContentHTML, err = utilhtml.Sanitize(input.Data.ContentHTML); err != nil {

	// 		meta.Status = http.StatusBadRequest
	// 		output.AddError("content_html", "Error on sanitize HTML.") // Lang
	// 		return err
	// 	}

	// 	// if input.Data.ContentHTML, err = utilmd.MD2HTML(input.Data.ContentMarkdown); err != nil {

	// 	// 	meta.Status = http.StatusBadRequest                                    // e.g., duplicate data 409
	// 	// 	output.AddError("content_markdown", "Error on converting markdown to HTML.") // Lang
	// 	// 	return err
	// 	// }
	// }

	return nil

}

func (x *PostsEntityAPIController) validateDTO() error {

	dto := &x.DTO
	input := &dto.Input

	c := x.webCtxt

	if err := c.Bind(input); err != nil {
		return err
	}

	return x.validateDTOFields()

}
func (x *PostsEntityAPIController) handleGET() (err error) {
	dto := &x.DTO
	input := &dto.Input
	meta := &dto.Meta
	output := &dto.Output
	srv := x.appService.BlogAdmin()

	var res *service.BlogPost

	if input.Code != "" { // /code/:code
		res, err = srv.Posts().FindByCode(input.Code)
	} else { // /:id
		res, err = srv.Posts().FindByID(input.ID)
	}

	if err != nil {
		return err
	}

	if res == nil {
		meta.Status = http.StatusNotFound
	} else {
		output.Data.BlogPost = *res // copy
		// output.Data.ContentMarkdown = ""
	}

	return nil
}

func (x *PostsEntityAPIController) handlePOST() (err error) {
	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	srv := x.appService.BlogAdmin()

	output.Data = input.Data
	output.Data.ID = 0 // reset ID

	return srv.Posts().Create(&output.Data.BlogPost)

}
func (x *PostsEntityAPIController) handlePUT() (err error) {
	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	srv := x.appService.BlogAdmin()

	output.Data = input.Data
	output.Data.ID = input.Data.ID // reset ID
	return srv.Posts().Update(&output.Data.BlogPost)

}
func (x *PostsEntityAPIController) handleDELETE() error {

	dto := &x.DTO
	input := &dto.Input
	output := &dto.Output
	srv := x.appService.BlogAdmin()

	output.Data.ID = input.ID // reset ID
	return srv.Posts().Delete(output.Data.ID)

}
func (x *PostsEntityAPIController) handleDTO() error {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output

	if meta.Status > 0 {
		return nil // stop processing
	}

	switch {

	case x.IsGET:
		return x.handleGET()
	case x.IsPOST:
		return x.handlePOST()
	case x.IsPUT:
		return x.handlePUT()
	case x.IsDELETE:
		return x.handleDELETE()
	default:
		{
			meta.Status = http.StatusMethodNotAllowed
			output.Message = "Method action undef"
		}
	}

	return nil
}
func (x *PostsEntityAPIController) responseDTOAsAPI() (err error) {

	dto := &x.DTO
	meta := &dto.Meta
	output := &dto.Output
	c := x.webCtxt
	controller.CsrfToHeader(c)

	if meta.Status == 0 {
		meta.Status = http.StatusOK
	}

	// if x.IsPOST || x.IsPUT {
	// 	// clean or modyfy omitColumns
	// 	// output.Data.ContentMarkdown = ""
	// 	// output.Data.ContentHTML = ""
	// }

	return c.JSON(meta.Status, output)

}

func (x *PostsEntityAPIController) responseDTO() (err error) {

	return x.responseDTOAsAPI()

}
