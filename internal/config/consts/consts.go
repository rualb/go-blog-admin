package consts

const AppName = "go-blog-admin"

// App consts
const (
	LongTextLength    = 32767 // int(int16(^uint16(0) >> 1)) // equivalent of short.MaxValue
	DefaultTextLength = 100

	TitleTextLengthTiny   = 12
	TitleTextLengthSmall  = 25
	TitleTextLengthInfo   = 35
	TitleTextLengthMedium = 50
	TitleTextLengthLarge  = 100
)

const (
	// PathAPI represents the group of PathAPI.
	PathAPI = "/api"
)
const (
	RoleAdmin = "admin"
)

const (
	BlogRolePrefix  = "blog_"
	BlogRoleAccess  = "blog_access"
	BlogRoleAdd     = "blog_add"
	BlogRoleEdit    = "blog_edit"
	BlogRoleView    = "blog_view"
	BlogRoleDelete  = "blog_delete"
	BlogRolePublish = "blog_publish"
)

//nolint:gosec
const (
	PathSysMetricsAPI = "/sys/api/metrics"

	PathBlogAdminPingDebugAPI = "/blog-admin/api/ping"

	PathBlog = "/blog"

	PathBlogAdmin            = "/blog-admin"
	PathBlogAdminAssets      = "/blog-admin/assets"
	PathBlogAdminPosts       = "/blog-admin/posts"
	PathBlogAdminPostsEntity = "/blog-admin/posts/:code" // GET
	PathBlogAdminStatusAPI   = "/blog-admin/api/status"  // private
	PathBlogAdminConfigAPI   = "/blog-admin/api/config"  // public

	PathBlogAdminPostsAPI             = "/blog-admin/api/posts"            // LIST POST
	PathBlogAdminPostsEntityAPI       = "/blog-admin/api/posts/:id"        // GET PUT DELETE
	PathBlogAdminPostsEntityByCodeAPI = "/blog-admin/api/posts/:code/code" // GET
)
