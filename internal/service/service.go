package service

import (
	"go-blog-admin/internal/config"
	"go-blog-admin/internal/i18n"
	"go-blog-admin/internal/repository"
	xlog "go-blog-admin/internal/util/utillog"
	"net/http"
	"os"
	"strings"
	"time"
)

// AppService all services
type AppService interface {
	Account() AccountService

	Config() *config.AppConfig
	// Logger() logger.AppLogger

	UserLang(code string) i18n.UserLang
	HasLang(code string) bool

	Vault() VaultService

	BlogAdmin() BlogAdminService
	Repository() repository.AppRepository
}

type defaultAppService struct {
	accountService AccountService
	// container      container.AppContainer
	vaultService VaultService

	blogAdminService BlogAdminService
	configSource     *config.AppConfigSource
	repository       repository.AppRepository
	lang             i18n.AppLang
}

func (x *defaultAppService) mustConfig() {

	d, _ := os.Getwd()

	xlog.Info("current work dir: %v", d)

	x.configSource = config.MustNewAppConfigSource()

	appConfig := x.configSource.Config() // first call, init

	mustConfigRuntime(appConfig)

}

func (x *defaultAppService) mustBuild() {

	var err error

	appConfig := x.configSource.Config() // first call, init

	x.lang = i18n.NewAppLang(appConfig)

	x.repository = repository.MustNewRepository(appConfig) // , appLogger)

	//

	mustCreateRepository(x)

	if x.vaultService, err = newVaultService(x); err != nil {
		xlog.Panic("vault service: %v", err)
	}

	x.blogAdminService = newBlogAdminService(x)

	x.accountService = newAccountService(x)

}

func mustConfigRuntime(appConfig *config.AppConfig) {
	t, ok := http.DefaultTransport.(*http.Transport)

	if ok {
		x := appConfig.HTTPTransport

		if x.MaxIdleConns > 0 {
			xlog.Info("set Http.Transport.MaxIdleConns=%v", x.MaxIdleConns)
			t.MaxIdleConns = x.MaxIdleConns
		}
		if x.IdleConnTimeout > 0 {
			xlog.Info("set Http.Transport.IdleConnTimeout=%v", x.IdleConnTimeout)
			t.IdleConnTimeout = time.Duration(x.IdleConnTimeout) * time.Second
		}
		if x.MaxConnsPerHost > 0 {
			xlog.Info("set Http.Transport.MaxConnsPerHost=%v", x.MaxConnsPerHost)
			t.MaxConnsPerHost = x.MaxConnsPerHost
		}

		if x.MaxIdleConnsPerHost > 0 {
			xlog.Info("set Http.Transport.MaxIdleConnsPerHost=%v", x.MaxIdleConnsPerHost)
			t.MaxIdleConnsPerHost = x.MaxIdleConnsPerHost
		}

	} else {
		xlog.Error("cannot init http.Transport")
	}
}

func MustNewAppServiceProd() AppService {

	appService := &defaultAppService{}

	appService.mustConfig()
	appService.mustBuild()

	return appService
}
func MustNewAppServiceTesting() AppService {

	env := map[string]string{
		"env": "testing",
	}

	for k, v := range env {
		_ = os.Setenv(strings.ToUpper("app_"+k), v)
	}

	return MustNewAppServiceProd()
}
func (x *defaultAppService) Account() AccountService   { return x.accountService }
func (x *defaultAppService) Config() *config.AppConfig { return x.configSource.Config() }

// func (x *appService) Logger() logger.AppLogger       { return x.container.Logger() }

func (x *defaultAppService) UserLang(code string) i18n.UserLang { return x.lang.UserLang(code) }
func (x *defaultAppService) HasLang(code string) bool           { return x.lang.HasLang(code) }

func (x *defaultAppService) Vault() VaultService { return x.vaultService }

func (x *defaultAppService) BlogAdmin() BlogAdminService { return x.blogAdminService }

func (x *defaultAppService) Repository() repository.AppRepository { return x.repository }
