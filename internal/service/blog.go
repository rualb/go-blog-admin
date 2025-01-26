package service

import (
	"go-blog-admin/internal/config/consts"
	"go-blog-admin/internal/util/utilaccess"
	"go-blog-admin/internal/util/utilorm"

	"go-blog-admin/internal/util/utilpaging"
	"time"
)

// BlogPost blog post
type BlogPost struct {
	ID              int64  `json:"id" gorm:"size:255;primaryKey;autoIncrement"`
	Code            string `json:"code" gorm:"size:255;uniqueIndex"`
	Title           string `json:"title,omitempty" gorm:"size:255"`
	ContentMarkdown string `json:"content_markdown,omitempty" gorm:"size:32767"`
	// ContentHTML string    `json:"content_html,omitempty" gorm:"size:32767"` // 2^16-1
	Status    string    `json:"status" gorm:"size:255"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	IsPublished bool `json:"is_published,omitempty"` // Indicates whether the post is available for public use
	IsListed    bool `json:"is_listed,omitempty"`    // Indicates whether the post is visible in the posts list // IsIndexed

}

func (x *BlogPost) Fill() {
	// if x.ID == "" {
	// 	x.ID = uuid.New().String()
	// }
}

type BlogPostDAO struct {
	appService AppService
}

func (x *BlogPostDAO) Check(filter *utilpaging.PagingInputDTO) {
	filter.Limit = min(filter.Limit, 10) // validate
}
func (x *BlogPostDAO) Permissions(userAccount *UserAccount, dto *utilaccess.PermissionsDTO) {
	dto.Fill(userAccount.Roles, consts.BlogRolePrefix)
}

func (x *BlogPostDAO) Where(sql *utilorm.SQLBuilder,
	filter *utilpaging.PagingInputDTO,
	fieldOptions utilorm.FieldOptions,
) (err error) {

	if fieldOptions == nil {
		return nil
	}

	sql.WhereText.WriteString("1=1")

	if err = utilorm.WhereSearch(sql, filter.Search, fieldOptions); err != nil {
		return err
	}

	if err = utilorm.WhereFromFilter(sql, filter, fieldOptions); err != nil {
		return err
	}

	return err
}

func (x *BlogPostDAO) Sort(
	filter *utilpaging.PagingInputDTO,
	sortOptions utilorm.SortOptions,
) (sqlSort string, err error) {

	if filter.Sort == "" {
		filter.Sort = sortOptions[""] // init default
	}

	sqlSort, ok := sortOptions[filter.Sort]
	if !ok {
		filter.Sort = sortOptions[""]
		sqlSort = sortOptions[filter.Sort]
	}

	return sqlSort, err
}

func (x *BlogPostDAO) Query(
	filter *utilpaging.PagingInputDTO,
	output *utilpaging.PagingOutputDTO[BlogPost],
	omitColumns []string,
	fieldOptions utilorm.FieldOptions,
	sortOptions utilorm.SortOptions,
) (err error) {

	x.Check(filter)

	repo := x.appService.Repository()

	sql := &utilorm.SQLBuilder{}

	err = x.Where(sql, filter, fieldOptions)
	if err != nil {
		return err
	}

	unsafeSQLSort, err := x.Sort(filter, sortOptions)
	if err != nil {
		return err
	}
	var count int64

	err = repo.Model(&BlogPost{}).
		Where(sql.WhereText.String(), sql.WhereArgs...).
		Count(&count).Error

	if err != nil {
		return err
	}

	info := filter.Info(int(count))
	output.Fill(filter, info)
	output.Data = make([]*BlogPost, 0, info.Limit)

	if omitColumns == nil {
		omitColumns = []string{}
	}

	err = repo.
		Where(sql.WhereText.String(), sql.WhereArgs...).
		Order(unsafeSQLSort).
		Omit(omitColumns...). // ContentMarkdown ContentHTML
		Limit(info.Limit).
		Offset(info.Offset).
		Find(&output.Data).Error

	if err != nil {
		return err
	}

	return err
}

func (x *BlogPostDAO) FindByID(id int64) (*BlogPost, error) {
	if id == 0 {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	data := new(BlogPost)

	result := x.appService.Repository().Find(data, "id = ?", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return data, nil
}
func (x *BlogPostDAO) FindByCode(code string) (*BlogPost, error) {
	if code == "" {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	data := new(BlogPost)

	result := x.appService.Repository().Find(data, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return data, nil
}
func (x *BlogPostDAO) ID(id int64) (int64, error) {
	if id == 0 {
		return 0, nil // fmt.Errorf("id cannot be empty")
	}

	data := new(BlogPost)

	result := x.appService.Repository().Select("id").Limit(1).Find(data, "id = ? ", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return 0, result.Error
	}

	return data.ID, nil
}

func (x *BlogPostDAO) Code(code string) (int64, error) {
	if code == "" {
		return 0, nil // fmt.Errorf("id cannot be empty")
	}

	data := new(BlogPost)

	result := x.appService.Repository().Select("id").Find(data, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return 0, result.Error
	}

	return data.ID, nil
}

func (x *BlogPostDAO) Create(data *BlogPost) error {

	repo := x.appService.Repository()
	data.Fill()
	res := repo.Create(data)
	return res.Error

}
func (x *BlogPostDAO) Update(data *BlogPost) error {
	repo := x.appService.Repository()
	res := repo.Model(data).Select("*" /*over all columns*/).Updates(data)
	return res.Error
}
func (x *BlogPostDAO) Delete(id int64) error {

	if id == 0 {
		return nil
	}

	repo := x.appService.Repository()
	res := repo.Delete(&BlogPost{ID: id})
	return res.Error
}
