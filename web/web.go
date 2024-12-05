package web

import (
	"embed"
	"io/fs"
)

// /////////////////////
//
//go:embed go-blog-admin/blog-admin/assets/*
var fsBlogAdminAssetsFS embed.FS

func MustBlogAdminAssetsFS() fs.FS {
	res, err := fs.Sub(fsBlogAdminAssetsFS, "go-blog-admin/blog-admin/assets")
	if err != nil {
		panic(err)
	}
	return res
}

//go:embed  go-blog-admin/index.html
var fsBlogAdminIndexHTML embed.FS

func MustBlogAdminIndexHTML() string {

	data, err := fsBlogAdminIndexHTML.ReadFile("go-blog-admin/index.html")
	if err != nil {
		panic(err)
	}

	return string(data)
}
