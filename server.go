package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	db "payroll/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

func echo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("server handler started\n")
	defer fmt.Printf("server handler ended\n")
	select {
	case <-time.After(time.Second * 5):
		fmt.Fprint(w, "hi\n")
	case <-ctx.Done():
		fmt.Println(ctx.Err().Error())
		http.Error(w, ctx.Err().Error(), http.StatusInternalServerError)
	}
}

func header(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Following is your headers\n")
	for name, headers := range r.Header {
		for _, header := range headers {
			fmt.Fprint(w, fmt.Sprintf("%v %v\n", name, header))
		}
	}
}

type student struct {
	Name string
	Age  int8
}

type index struct {
	Title, Msg string
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type todo struct {
	Name string
	Done bool
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./assets/templates/*")
	r.MaxMultipartMemory = 32 << 20

	logFile, err := os.OpenFile("conn.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", index{"First Web Page", "Hello,world"})
	})

	r.GET("/images/:file", func(c *gin.Context) {
		filename := c.Param("file")
		c.File(fmt.Sprintf("../assets/images/%v", filename))
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

			c.SetCookie("id", "2345", 0, "", "", false, false)
			c.HTML(http.StatusOK, "todo.html", gin.H{
				"User":  info.User,
				"Todos": []todo{{"study", false}, {"work", true}},
			})

		} else {
			c.String(http.StatusBadRequest, validatedRes.Error())
		}
	})

	// r.POST("/loginJson", func(c *gin.Context) {
	// 	log.Printf("%v : %v\n", c.Request.Method, c.Request.RequestURI)
	// 	var info Login
	// 	if err = c.ShouldBindJSON(&info); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	validatedRes := db.ValidateUser(info.User, info.Password)
	// 	if validatedRes == nil {
	// 		c.String(http.StatusOK, "You has successfully authenticated")
	// 	} else {
	// 		c.String(http.StatusBadRequest, validatedRes.Error())
	// 	}
	// })

	// r.POST("/form", func(c *gin.Context) {
	// 	name := c.PostForm("username")
	// 	pass := c.PostForm("password")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"name":     name,
	// 		"password": pass,
	// 	})

	// })

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

	r.Run()
}
