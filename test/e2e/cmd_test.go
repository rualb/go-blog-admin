package e2e

import (
	"encoding/base64"
	"fmt"
	xcmd "go-blog-admin/internal/cmd"
	"go-blog-admin/internal/config"
	"go-blog-admin/internal/config/consts"
	"go-blog-admin/internal/repository"
	"go-blog-admin/internal/service"
	"go-blog-admin/internal/util/utilhttp"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func setup() {

	cfgSrc := config.MustNewAppConfigSource()
	cfg := cfgSrc.Config()
	db := repository.MustNewRepository(cfg) // not via app-service

	{
		// migrate
		for _, x := range []any{
			&service.UserAccount{},
			&service.VaultKey{},
			&service.BlogPost{},
		} {
			if err := db.AutoMigrate(x); err != nil {
				log.Fatalf("migration: %v", err)
			}
		}

		// seed
		{
			key := &service.VaultKey{}
			key.ID = "test-only-1"
			key.AuthKey = base64.StdEncoding.EncodeToString([]byte(strings.Repeat("A", 64)))
			// insert or update
			if res := db.Save(key); res.Error != nil {
				log.Fatalf("new vault key: %v", res.Error)
			}
		}
	}
}

func TestMain(m *testing.M) {
	fmt.Println("testMain")
	setup()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestBlogAddData(t *testing.T) {

	os.Setenv("APP_ENV", "testing")
	os.Setenv("APP_BLOG_ADMIN", "true")

	srv := service.MustNewAppServiceProd()
	blogSrv := srv.BlogAdmin()
	dao := blogSrv.Posts()

	for i := 1; i < 199; i++ {
		tmp := service.BlogPost{
			Code:            fmt.Sprintf("test-post-%v", i),
			Title:           fmt.Sprintf("test post %v", i),
			ContentMarkdown: fmt.Sprintf("test post %v **Content**", i),
			// ContentHTML: fmt.Sprintf("test post %v <b>Content</b>", i),
		}

		if id, _ := dao.Code(tmp.Code); id > 0 {
			continue
		}

		tmp.Fill()
		if err := dao.Create(&tmp); err != nil {
			t.Fatal(err)
		}
	}

}

func getAuthToken(srv service.AppService) (authToken string, err error) {

	// TODO add user admin

	test_user := "blog_user_test"

	{
		//acc, _ := srv.Account().FindByID(test_user)
		//if acc == nil {
		repo := srv.Repository()

		res := repo.Driver().Save(&service.UserAccount{
			ID:       test_user,
			Username: test_user,
			Roles: strings.Join([]string{
				consts.BlogRoleAccess,
				consts.BlogRoleAdd,
				consts.BlogRoleView,
				consts.BlogRoleEdit,
				consts.BlogRolePublish,
				consts.BlogRoleDelete,
			}, " "),
		})
		if res.Error != nil {
			return "", res.Error
		}
		//}

	}

	secret, err := srv.Vault().CurrentKey()
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = secret.ID
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = test_user
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 720).Unix()
	claims["iss"] = "auth"

	tokenString, err := token.SignedString(secret.AuthKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func TestBlogAdmin(t *testing.T) {
	// Setup Echo context

	//   Trafalgar

	var err error
	var authToken string
	os.Setenv("APP_ENV", "testing")
	os.Setenv("APP_BLOG_ADMIN", "true")

	cmd := xcmd.Command{}

	go cmd.Exec()

	time.Sleep(1 * time.Second)

	// Authorization:Bearer jwt

	{
		srv := cmd.AppService
		blogSrv := srv.BlogAdmin()
		dao := blogSrv.Posts()

		if authToken, err = getAuthToken(srv); err != nil {
			t.Fatal(err)
		}

		data := &service.BlogPost{
			Code:            "test-post-1",
			Title:           "Test post 1",
			ContentMarkdown: "Test post 1 **Content**",
			// ContentHTML: "Test post 1 <b>Content</b>",
		}

		if data.ID, err = dao.Code(data.Code); err != nil {
			t.Fatal(err)
		}

		data.Fill()

		if err := dao.Delete(data.ID); err != nil {
			t.Fatal(err)
		}

		if id, err := dao.Code(data.Code); err != nil {
			t.Fatal(err)
		} else if id > 0 {
			t.Fatal("rec exists")
		}

		if err := dao.Create(data); err != nil {
			t.Fatal(err)
		}

		if id, err := dao.Code(data.Code); err != nil {
			t.Fatal(err)
		} else if id == 0 {
			t.Fatal("rec not exists")
		}

		data.ContentMarkdown += " #UPDATED"
		// data.ContentHTML += " #UPDATED"

		if err := dao.Update(data); err != nil {
			t.Fatal(err)
		}

		header := map[string]string{
			`Authorization`: `Bearer ` + authToken,
		}

		urls := []struct {
			title  string
			url    string
			query  map[string]string
			header map[string]string
			search []string
		}{

			{title: "test 1", search: []string{"Test post 1"},
				url:   "http://127.0.0.1:18180" + strings.ReplaceAll(consts.PathBlogAdminPostsEntityAPI, ":id", strconv.FormatInt(data.ID, 10)),
				query: map[string]string{}, header: header},
		}

		for _, itm := range urls {

			t.Run(itm.title, func(t *testing.T) {

				t.Logf("url %v", itm.url)
				dataArr, err := utilhttp.GetBytes(itm.url, itm.query, itm.header)

				if err != nil {
					t.Errorf("Error : %v", err)
				}
				dataTxt := string(dataArr)
				for _, v := range itm.search {
					if !strings.Contains(string(dataTxt), v) {
						t.Errorf("Error on %v", itm.url)
					}
				}

			})

		}

		if err := dao.Delete(data.ID); err != nil {
			t.Fatal(err)
		}

		if id, err := dao.Code(data.Code); err != nil {
			t.Fatal(err)
		} else if id > 0 {
			t.Fatal("rec exists")
		}

	}
	cmd.Stop()

	time.Sleep(1 * time.Second)

}
