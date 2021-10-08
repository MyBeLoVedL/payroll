package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	db "payroll/db/sqlc"
	"payroll/misc"
	"time"

	"github.com/gin-gonic/gin"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("./resources/templates/main.html", "./resources/templates/login.html", "./resources/templates/todo.html")
	r.MaxMultipartMemory = 32 << 20

	logFile, err := os.OpenFile("conn.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	r.Static("./resources", "./resources")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
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
		id, validatedRes := db.ValidateUser(info.User, info.Password)
		if validatedRes == nil {
			//* send session id to client here
			sid := misc.GSS.AddSession(id)

			c.SetCookie("sid", sid, 0, "", "", false, false)
			c.HTML(http.StatusOK, "main.html", gin.H{
				"showUsername": info.User,
			})
		} else {
			c.String(http.StatusBadRequest, validatedRes.Error())
		}
	})

	r.Any("/updatePay", func(c *gin.Context) {
		method := c.PostForm("updatePay")
		sid, _ := c.Cookie("sid")
		session, err := misc.GSS.Get(sid)
		switch method {
		case "pick_up":
			db.UpdatePayment(method, session.User)
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
		case "mail":
			mail := c.PostForm("mail")
			err = db.UpdatePaymentWIthMail(session.User, mail)
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
		case "deposit":
			bankName := c.PostForm("bankName")
			bankAccount := c.PostForm("bankAccount")
			err = db.UpdatePaymentWithBank(session.User, bankName, bankAccount)
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})

		}
		c.JSON(http.StatusOK, gin.H{
			"error":  err,
			"method": method,
		})

	})
	r.Any("/timecard", func(c *gin.Context) {
		type timecardParam struct {
			Charge int       `form:"charge"`
			Hours  int       `form:"hours"`
			Date   time.Time `form:"date" time_format:"2006-01-02"`
		}

		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)

		var arg timecardParam
		err := c.BindQuery(&arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		err = db.UpdateTimecard(session.User, arg.Charge, arg.Hours, arg.Date)
		log.Printf("%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
	})

	r.Run(":9999")
}
