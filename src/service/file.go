package service

import (
	"bytes"
	"coding.net/baoquan2017/candy-backend/src/common/constant"
	"coding.net/baoquan2017/candy-backend/src/common/logger"
	"coding.net/baoquan2017/candy-backend/src/common/util"
	"coding.net/baoquan2017/candy-backend/src/common/uuid"
	"coding.net/baoquan2017/candy-backend/src/config"
	"coding.net/baoquan2017/candy-backend/src/model"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/jinzhu/gorm"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex *sync.Mutex

func init() {
	mutex = new(sync.Mutex)
}

func UploadFile(filename string, f multipart.File, db *gorm.DB) (string, error, uint64) {
	tmpPath := filepath.Join(config.GetStorageRoot(), uuid.NewUUID().String()+filepath.Ext(filename))
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		logger.Error(err)
		return "", err, 0
	}
	defer os.Remove(tmpPath)

	bytes, err := io.Copy(tmpFile, f)
	tmpFile.Close()
	if err != nil {
		logger.Error(err)
		return "", err, 0
	}

	return saveDbFile(tmpPath, filename, "", bytes, db)
}

func saveDbFile(tmpPath, filename, sourceUrl string, size int64, db *gorm.DB) (string, error, uint64) {
	var dbFile model.File

	mutex.Lock()
	relDir := time.Now().Format("/2006/01/02/15/0405")
	relPath := filepath.Join(relDir, filename)
	relHash := util.StringHashToUint32(relPath)

	var count int
	if err := db.Model(model.File{}).Where("rel_hash=?", relHash).Count(&count).Error; err != nil {
		mutex.Unlock()
		logger.Error(err)
		return "", err, 0
	}
	if count > 0 {
		relPath = filepath.Join(relDir, fmt.Sprint(rand.Intn(1000))+filename)
		relHash = util.StringHashToUint32(relPath)
	}
	mutex.Unlock()

	if err := os.MkdirAll(filepath.Join(config.GetStorageRoot(), relDir), 0755); err != nil {
		logger.Error(err)
		return "", err, 0
	}

	dbFile.SourceUrl = sourceUrl
	dbFile.SourceHash = util.StringHashToUint32(dbFile.SourceUrl)
	dbFile.RelPath = relPath
	dbFile.RelHash = relHash

	absPath := filepath.Join(config.GetStorageRoot(), dbFile.RelPath)
	if err := os.Rename(tmpPath, absPath); err != nil {
		logger.Error(err)
		return "", err, 0
	}

	dbFile.Ext = filepath.Ext(filename)
	dbFile.Filename = filename
	dbFile.Filesize = size

	ext := strings.ToLower(dbFile.Ext)
	if ext == ".gif" || ext == ".jpeg" || ext == ".jpg" || ext == ".bmp" || ext == ".png" {
		dbFile.ContentType = constant.FileContentTypePicture
	}

	logger.Debug(dbFile)
	if err := db.Create(&dbFile).Error; err != nil {
		logger.Error(err)
		return "", err, 0
	}

	return dbFile.RelPath, nil, dbFile.Id
}

func GetFullUrl(relPath string) string {
	fullUrl := config.GetBindAddr() + "/ims" + relPath
	// 2017-12-22: 返回urlencode的地址，避免ios解析出错
	lastIndex := strings.LastIndex(fullUrl, `/`)
	return fullUrl[:lastIndex+1] + url.QueryEscape(fullUrl[lastIndex+1:])
}

func IsThumb(path string) (string, bool) {
	n := filepath.Base(path)
	yes := strings.Contains(n, "@w") && strings.Contains(n, "_h")
	if yes {
		base := filepath.Base(path)
		sl := strings.Split(base, "@")
		filename := sl[0] + filepath.Ext(base)
		return filepath.Join(filepath.Dir(path), filename), yes
	} else {
		return path, yes
	}
}

func CreateThumb(path string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != "" {
		if ext != ".jpg" && ext != ".jpeg" &&
			ext != ".png" && ext != ".bmp" && ext != ".gif" {
			return nil, errors.New("Unsupport file format " + ext)
		}
	}

	// example path: ../file/o/2016-05-19/91bd48ce-1dae-11e6-9d7f-408d5cdf2c91@w360_h100.jpg
	// take the string after "@" (w360_h100.jpg) to extract the size (360*270)
	sl := strings.Split(path, "@")
	sizeString := sl[1]
	ws := sizeString[strings.Index(sizeString, "@w")+2 : strings.Index(sizeString, "_h")]
	var hs string
	if strings.Index(sizeString, ".") == -1 {
		hs = sizeString[strings.Index(sizeString, "_h")+2:]
	} else {
		hs = sizeString[strings.Index(sizeString, "_h")+2 : strings.Index(sizeString, ".")]
	}
	w, _ := strconv.Atoi(ws)
	h, _ := strconv.Atoi(hs)

	srcPath := sl[0] + ext
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	i, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	m := imaging.Fill(i, w, h, imaging.Center, imaging.Lanczos)

	var buffer bytes.Buffer
	if ext == ".jpg" || ext == ".jpeg" || ext == "" {
		if err := jpeg.Encode(&buffer, m, nil); err != nil {
			return nil, err
		}
	} else if ext == ".png" {
		if err := png.Encode(&buffer, m); err != nil {
			return nil, err
		}
	} else if ext == ".bmp" {
		if err := bmp.Encode(&buffer, m); err != nil {
			return nil, err
		}
	} else if ext == ".gif" {
		if err := gif.Encode(&buffer, m, nil); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}
