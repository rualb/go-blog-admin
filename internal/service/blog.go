package service

import (
	"go-blog-admin/internal/config/consts"
	"go-blog-admin/internal/util/utilaccess"
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

func (x *BlogPostDAO) Where(filter *utilpaging.PagingInputDTO) (whereCondition string, whereArgs []any, err error) {
	whereCondition = "1=1"
	whereArgs = []any{}

	if v := filter.Search; v != "" { // filter.GetFilter("text");
		// content_markdown content_html
		whereCondition += " and (title ilike ? or content_markdown ilike ?)" // " and (title ilike ? or content_markdown ilike ?)"
		whereArgs = append(whereArgs, "%"+v+"%", "%"+v+"%")
	}

	if whereCondition == "1=1" {
		whereCondition = ""
	}

	return whereCondition, whereArgs, err
}

func (x *BlogPostDAO) Sort(filter *utilpaging.PagingInputDTO) (sqlSort string, err error) {

	sqlSort = "id desc"
	switch filter.Sort {
	case "-id":
		sqlSort = "id desc"
	case "id":
		sqlSort = "id asc"
	case "-code":
		sqlSort = "code desc"
	case "code":
		sqlSort = "code asc"
	default:
		filter.Sort = "-id"
	}

	return sqlSort, err
}

func (x *BlogPostDAO) Query(filter *utilpaging.PagingInputDTO, output *utilpaging.PagingOutputDTO[BlogPost], omitColumns *[]string) (err error) {

	x.Check(filter)

	repo := x.appService.Repository()

	sqlWhere, sqlWhereArgs, _ := x.Where(filter)
	sqlSort, _ := x.Sort(filter)

	var count int64

	err = repo.Model(&BlogPost{}).
		Where(sqlWhere, sqlWhereArgs...).
		Count(&count).Error

	if err != nil {
		return err
	}

	info := filter.Info(int(count))
	output.Fill(filter, info)
	output.Data = make([]*BlogPost, 0, info.Limit)

	if omitColumns == nil {
		omitColumns = &[]string{}
	}

	err = repo.
		Where(sqlWhere, sqlWhereArgs...).
		Order(sqlSort).
		Omit(*omitColumns...). // ContentMarkdown ContentHTML
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

	user := new(BlogPost)

	result := x.appService.Repository().Find(user, "id = ?", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
func (x *BlogPostDAO) FindByCode(code string) (*BlogPost, error) {
	if code == "" {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Find(user, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
func (x *BlogPostDAO) ID(id int64) (int64, error) {
	if id == 0 {
		return 0, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Select("id").Limit(1).Find(user, "id = ? ", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return 0, result.Error
	}

	return user.ID, nil
}

func (x *BlogPostDAO) Code(code string) (int64, error) {
	if code == "" {
		return 0, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Select("id").Find(user, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return 0, result.Error
	}

	return user.ID, nil
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
