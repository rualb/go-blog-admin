package blogadmin

import (
	"go-blog-admin/internal/config/consts"
	"net/http"

	"github.com/labstack/echo/v4"
)

var mapPathToRoles = map[string][]string{
	// POST /api HTTP/1.1

	http.MethodGet + " " + consts.PathBlogAdminStatusAPI: {consts.BlogRoleAccess},
	http.MethodGet + " " + consts.PathBlogAdminConfigAPI: {consts.BlogRoleAccess},

	http.MethodGet + " " + consts.PathBlogAdminPosts:       {consts.BlogRoleAccess},
	http.MethodGet + " " + consts.PathBlogAdminPostsEntity: {consts.BlogRoleAccess},
	http.MethodGet + " " + consts.PathBlogAdminPostsAPI:    {consts.BlogRoleAccess},

	http.MethodGet + " " + consts.PathBlogAdminPostsEntityAPI:    {consts.BlogRoleEdit, consts.BlogRoleView},
	http.MethodPost + " " + consts.PathBlogAdminPostsAPI:         {consts.BlogRoleAdd},
	http.MethodPut + " " + consts.PathBlogAdminPostsEntityAPI:    {consts.BlogRoleEdit},
	http.MethodDelete + " " + consts.PathBlogAdminPostsEntityAPI: {consts.BlogRoleDelete},

	http.MethodGet + " " + consts.PathBlogAdminPostsEntityByCodeAPI: {consts.BlogRoleAccess},
}

func RolesForAPI(c echo.Context) []string {
	key := c.Request().Method + " " + c.Path() // Method + Grop path + Route path
	roles := mapPathToRoles[key]
	return roles
}

func RolesForAssets(_ echo.Context) []string {
	return []string{consts.BlogRoleAccess}
}
