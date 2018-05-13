package front

import (
	"coding.net/baoquan2017/candy-backend/src/common/constant"
	"coding.net/baoquan2017/candy-backend/src/common/logger"
	"coding.net/baoquan2017/candy-backend/src/common/util"
	"coding.net/baoquan2017/candy-backend/src/config"
	"coding.net/baoquan2017/candy-backend/src/model"
	"coding.net/baoquan2017/candy-backend/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func FileUploadHandler(c *gin.Context) {
	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	if err := c.Request.ParseMultipartForm(64 << 20); err != nil { // parse data, set memory size to 64M
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Connection", "close")

	formdata := c.Request.MultipartForm
	files := formdata.File["file"]
	if len(files) <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fh := files[0]
	file, err := fh.Open()
	defer file.Close()
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fileName, _ := url.QueryUnescape(fh.Filename)

	path, err, fileId := service.UploadFile(fileName, file, db)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": service.GetFullUrl(path), "file_id": fileId})
}

func FileServeHandler(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	relPath := strings.Split(c.Request.RequestURI, "?")[0]
	if !strings.HasPrefix(relPath, "/ims") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	relPath = strings.TrimPrefix(relPath, "/ims")

	// url decode
	relPath, _ = url.QueryUnescape(relPath)
	realPath, ok := service.IsThumb(relPath)

	c.Header("Connection", "close")

	realHash := util.StringHashToUint32(realPath)

	var dbFile model.File
	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	if err := db.Where("rel_hash=? and rel_path=?", realHash, realPath).First(&dbFile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	dbFile.AccessCount++
	now := time.Now()
	dbFile.LastAccess = &now
	db.Save(&dbFile)

	absPath := filepath.Join(config.GetStorageRoot(), relPath)

	if ok {
		data, err := service.CreateThumb(absPath)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		me := mime.TypeByExtension(filepath.Ext(filepath.Base(absPath)))
		if me == "" {
			me = "image/jpeg"
		}

		c.Data(http.StatusOK, me, data)
	} else {
		http.ServeFile(c.Writer, c.Request, absPath)
	}
}
