package main

import (
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Output struct {
	Timestamp int64   `json:"timestamp"`
	Images    []Image `json:"images"`
}

type OutputNo struct {
	Images []Image `json:"images"`
}

type Image struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

func GetImages(c *gin.Context) {
	dbconn := connect()
	ctx := c.Request.Context()
	var err error
	restimeStr := c.Query("time")
	var restime int64
	if restimeStr != "" {
		restime, err = strconv.ParseInt(restimeStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
	data, err := dbconn.QueryContext(ctx, "SELECT img, type FROM images")
	if err != nil {
		log.Fatal(err)
	}

	timestamp := timme(c)

	var images []Image
	for data.Next() {
		var image Image
		data.Scan(&image.URL, &image.Type)
		images = append(images, image)
	}

	if restime == timestamp {
		c.Status(http.StatusOK)
	} else {
		output := Output{
			Timestamp: timestamp,
			Images:    images,
		}
		c.IndentedJSON(http.StatusOK, output)
	}
	data.Close()
}

func ImagesById(c *gin.Context) {
	dbconn := connect()
	ctx := c.Request.Context()
	var err error
	restimeStr := c.Query("time")
	var restime int64
	if restimeStr != "" {
		restime, err = strconv.ParseInt(restimeStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
	id := c.Param("id")
	data, err := dbconn.QueryContext(ctx, "SELECT img, type FROM images WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	timestamp := timme(c)

	var images []Image
	for data.Next() {
		var image Image
		data.Scan(&image.URL, &image.Type)
		images = append(images, image)
	}
	if restime == timestamp {
		c.Status(http.StatusOK)
	} else {
		output := Output{
			Timestamp: timestamp,
			Images:    images,
		}
		c.IndentedJSON(http.StatusOK, output)
	}
	data.Close()
}

func ImagesByType(c *gin.Context) {
	dbconn := connect()
	ctx := c.Request.Context()
	var err error
	restimeStr := c.Query("time")
	var restime int64
	if restimeStr != "" {
		restime, err = strconv.ParseInt(restimeStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}

	timestamp := timme(c)

	var typ = c.Param("type")
	data, err := dbconn.QueryContext(ctx, "SELECT img, type FROM images WHERE type = ?", typ)
	if err != nil {
		log.Fatal(err)
	}
	var images []Image
	for data.Next() {
		var image Image
		data.Scan(&image.URL, &image.Type)
		images = append(images, image)
	}
	if restime == timestamp {
		c.Status(http.StatusOK)
	} else {
		output := Output{
			Timestamp: timestamp,
			Images:    images,
		}
		c.IndentedJSON(http.StatusOK, output)
	}
	data.Close()
}

func AddImg(c *gin.Context, file *multipart.FileHeader, imgType string) {
	dbconn := connect()
	ctx := c.Request.Context()
	// сохраняем файл на сервере
	err := c.SaveUploadedFile(file, "downloads/"+file.Filename)
	if err != nil {
		log.Fatal(err)
	}

	// сохраняем путь к файлу и тип изображения в базу данных
	_, err = dbconn.ExecContext(ctx, "INSERT INTO images (img, type) VALUES (?, ?)", "downloads/"+file.Filename, imgType)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image added successfully",
	})
}

// func hash(c *gin.Context) string {
// 	dbconn := connect()
// 	ctx := c.Request.Context()
// 	data, err := dbconn.QueryContext(ctx, "SELECT img, type FROM images")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var img string
// 	var imgType string
// 	var alldata string
// 	for data.Next() {
// 		eror := data.Scan(&img, &imgType)
// 		if eror != nil {
// 			log.Fatal(eror)
// 		}
// 		alldata += img + imgType
// 	}

// 	hash := sha256.Sum256([]byte(alldata))
// 	return string(hash[:])
// }

func timme(c *gin.Context) int64 {
	dbconn := connect()
	ctx := c.Request.Context()
	var strTime string
	err := dbconn.QueryRowContext(ctx, "SELECT MAX(updateat) FROM images").Scan(&strTime)
	if err != nil {
		log.Fatal(err)
	}
	layout := "2006-01-02 15:04:05"    // Формат строки даты-времени в вашей базе данных
	loc, _ := time.LoadLocation("UTC") // Установите локальную зону, если требуется
	t, err := time.ParseInLocation(layout, strTime, loc)
	if err != nil {
		log.Fatal(err)
	}
	return t.Unix()
}
