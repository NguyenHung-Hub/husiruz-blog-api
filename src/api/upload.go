package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func (server *Server) uploadFile(ctx *gin.Context) {
	image, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("err load image")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	src, err := image.Open()
	if err != nil {
		fmt.Println("err image open")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	defer src.Close()

	fmt.Println(image.Filename, path.Ext(image.Filename))

	path := "assets/" + filepath.Join(filepath.Base(slug.Make(image.Filename)+path.Ext(image.Filename)))

	dst, err := os.Create(path)
	if err != nil {
		fmt.Println("err save image ")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("err save image ")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": "http://" + ctx.Request.Host + "/" + path,
	})
}

func (server *Server) uploadImageToCloud(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("err load image")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	image, err := file.Open()
	if err != nil {
		fmt.Println("err image open")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	defer image.Close()

	res, err := server.cld.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID:       file.Filename,
		Folder:         server.cfg.CloudinaryFolder,
		Format:         "webp",
		ResourceType:   "image",
		Transformation: "q_auto:eco",
		ResponsiveBreakpoints: uploader.ResponsiveBreakpointsParams{
			uploader.SingleResponsiveBreakpointsParams{
				CreateDerived: api.Bool(true),
				MinWidth:      200,
				MaxWidth:      1000,
			},
		},
	})
	if err != nil {
		log.Printf("upload file to cloudinary failed: %s", err)
		ctx.JSON(http.StatusBadRequest, errResponse(errors.New("upload file to cloudinary failed")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": res.URL, "data": res})
}
