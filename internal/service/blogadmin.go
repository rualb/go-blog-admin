package service

type BlogAdminService interface {
	Posts() *BlogPostDAO
}

type defaultBlogAdminService struct {
	appService AppService
	blogPost   BlogPostDAO
}

func newBlogAdminService(appService AppService) BlogAdminService {

	res := &defaultBlogAdminService{

		appService: appService,
		blogPost: BlogPostDAO{
			appService: appService,
		},
	}

	return res
}

func (x *defaultBlogAdminService) Posts() *BlogPostDAO {
	return &x.blogPost
}
