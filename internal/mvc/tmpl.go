package mvc

import (
	"encoding/json"
	"fmt"
	"go-blog-admin/internal/config"
	"go-blog-admin/internal/config/consts"
	"go-blog-admin/internal/util/utilstring"
	xweb "go-blog-admin/internal/web"
	"html/template"
	"io"
	"io/fs"
	"time"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer interface {
	Render(w io.Writer, name string, data any, c echo.Context) error
}

type templateRenderer struct {
	templates *template.Template
}

type ModelAppConfig struct {
	AppTitle        string
	CopyrightTitle  string
	GlobalVersion   string
	AssetsPublicURL string
}

type ModelConst struct {
	DefaultTextLength int // 100

}
type ModelPrm struct {
	Title      string
	AppTitle   string
	Csrf       string
	IsFragment bool

	LangCode    string
	AppConfig   ModelAppConfig
	AppConst    ModelConst
	RawJSONData any
}

type ModelAPI struct {
	IsAuthenticated bool
	Lang            func(text string, args ...any) string
	URL             func(path string, args ...string) string // path?args[0]=args[1]&args[2]=args[3]#args[4]
}

type ModelWrap struct {
	Model any
	Prm   ModelPrm
	API   ModelAPI
}

func NewModelWrap(c echo.Context, model any, isFragment bool, title string, appConfig *config.AppConfig, lang UserLang) (*ModelWrap, error) {

	_csrf, _ := c.Get("_csrf").(string)

	res := &ModelWrap{
		Model: model,
		API: ModelAPI{

			Lang:            lang.Lang,
			URL:             utilstring.AppendURL,
			IsAuthenticated: xweb.IsSignedIn(c),
		},
		Prm: ModelPrm{

			Title: title,

			Csrf: _csrf,

			IsFragment: isFragment,

			LangCode: lang.LangCode(),

			AppConfig: ModelAppConfig{
				AppTitle:        appConfig.Title,
				CopyrightTitle:  fmt.Sprintf("© %v %v", time.Now().UTC().Year(), appConfig.Title),
				GlobalVersion:   appConfig.Assets.GlobalVersion,
				AssetsPublicURL: appConfig.Assets.AssetsPublicURL,
			},
			AppConst: ModelConst{
				DefaultTextLength: consts.DefaultTextLength, // 100
			},
		},
	}

	data, err := json.Marshal(map[string]any{
		"test":                  `"<>`,
		"prm_lang_code":         res.Prm.LangCode,
		"prm_assets_public_url": res.Prm.AppConfig.AssetsPublicURL,
		"prm_global_version":    res.Prm.AppConfig.GlobalVersion,
	})

	if err != nil {
		return nil, err
	}

	res.Prm.RawJSONData = data
	return res, nil
}

func NewTemplateRenderer(viewsFs fs.FS, patterns ...string) TemplateRenderer {

	res := templateRenderer{}
	res.templates = template.Must(template.ParseFS(viewsFs, patterns...))

	return &res
}

// Render renders a template document
func (x *templateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	//
	isFragment := false /*use global(wrap) layout*/

	if model, ok := data.(*ModelWrap); ok {
		isFragment = model.Prm.IsFragment
	}

	if !isFragment {
		err := x.templates.ExecuteTemplate(w, "layout_header", data)

		if err != nil {
			return err
		}
	}

	{

		err := x.templates.ExecuteTemplate(w, name, data)
		if err != nil {
			return err
		}
	}

	if !isFragment {
		err := x.templates.ExecuteTemplate(w, "layout_footer", data)
		if err != nil {
			return err
		}
	}

	return nil
}
