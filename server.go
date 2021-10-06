package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	db "payroll/db/sqlc"
	"payroll/misc"

	"github.com/gin-gonic/gin"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("./assets/templates/todo.html", "./assets/templates/login_test.html")
	r.MaxMultipartMemory = 32 << 20

	logFile, err := os.OpenFile("conn.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login_test.html", nil)
	})

	r.GET("/cookie", func(c *gin.Context) {
		cookie, _ := c.Cookie("sid")
		c.JSON(http.StatusOK, gin.H{
			"session id": cookie,
		})
	})

	r.Any("/loginForm", func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		log.Printf("%v : %v : %v\n", c.Request.Method, c.Request.RequestURI, string(body))
		var info Login
		if err = c.ShouldBindQuery(&info); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validatedRes := db.ValidateUser(info.User, info.Password)
		if validatedRes == nil {
			//* send session id to client here
			sid := misc.GSS.AddSession(info.User)

			c.SetCookie("sid", sid, 0, "", "", false, false)
			c.HTML(http.StatusOK, "todo.html", nil)
		} else {
			c.String(http.StatusBadRequest, validatedRes.Error())
		}
	})

	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		name := file.Filename
		loc := fmt.Sprintf("./assets/upload/%v", name)
		os.OpenFile(loc, os.O_RDWR|os.O_CREATE, 0644)
		err := c.SaveUploadedFile(file, loc)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusOK, "%v uploaded", name)
	})

	r.Run(":9999")
}
