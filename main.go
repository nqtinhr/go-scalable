package main

import (
	"log"
	"todololist/component/tokenprovider"
	"todololist/middleware"
	ginitem "todololist/module/item/transport/gin"
	"todololist/module/user/storage"
	ginuser "todololist/module/user/transport/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// dsn := os.Getenv("MYSQL_CONN_STRING")
	dsn := "todolist:abc@123@tcp(127.0.0.1:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(db, err)

	// lấy store & middleware
	authStore := storage.NewSQLStore(db)
	tokenProvider := tokenprovider.NewTokenJWTProvider("jwt", "syshero.io")
	middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

	// lấy server
	r := gin.Default()
	r.Use(middleware.Recovery())

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("v1")
	{
		v1.POST("/login", ginuser.Login(db, tokenProvider))
		v1.POST("/register", ginuser.Register(db))
		v1.GET("/profile", middlewareAuth, ginuser.Profile())

		items := v1.Group("/items")
		{
			items.POST("/", middlewareAuth, ginitem.CreateItem(db))
			items.GET("/", ginitem.ListItem(db))
			items.GET("/:id", ginitem.GetItem(db))
			items.PATCH("/:id", middlewareAuth, ginitem.UpdateItem(db))
			items.DELETE("/:id", middlewareAuth, ginitem.DeleteItem(db))
		}
	}

}
