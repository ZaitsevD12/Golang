package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/images", GetImages)
	router.GET("/images/:id", ImagesById)
	router.GET("/images/type/:type", ImagesByType)
	router.Static("/downloads", "./downloads")
	//router.POST("/books", Form)
	//router.Run("localhost:8080")
	// router.PATCH("/checkout", checkoutBook)
	// router.PATCH("/return", returnBook)

	//Player, Team, Tournament, Partner

	router.LoadHTMLGlob("./index.html")

	// Обработчик GET-запроса на страницу с формой
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Обработчик POST-запроса на отправку формы
	router.POST("/submit", func(c *gin.Context) {
		// Получаем данные из формы
		imgType := c.PostForm("imgtype")
		file, err := c.FormFile("image")
		if err != nil {
			c.String(http.StatusBadRequest, "Ошибка при загрузке файла: %s", err.Error())
			return
		}

		inputPassword := c.PostForm("password")
		password := "sanyaalewnya"

		if inputPassword != password {
			c.String(http.StatusUnauthorized, "Неправильный пароль")
			return
		}

		AddImg(c, file, imgType)

		c.String(http.StatusOK, "Данные получены: тип=%s, имя=%s, размер=%d", imgType, file.Filename, file.Size)
	})

	// Запускаем сервер на порту 8080
	router.Run(":8080")
}
