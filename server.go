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

var prefix string

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./resources/templates/*")
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

			emp, err := db.GetUser(id)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			//* send session id to client here
			sid := misc.GSS.AddSession(&emp)

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
			db.UpdatePayment(method, session.User.ID)
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
		case "mail":
			mail := c.PostForm("mail")
			err = db.UpdatePaymentWIthMail(session.User.ID, mail)
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
		case "deposit":
			bankName := c.PostForm("bankName")
			bankAccount := c.PostForm("bankAccount")
			err = db.UpdatePaymentWithBank(session.User.ID, bankName, bankAccount)
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

		if db.IfCommitted(session.User.ID) {
			c.HTML(http.StatusOK, "timecard.html", gin.H{
				"Committed":    true,
				"showUsername": session.User.Name,
				"prefix":       prefix,
			})
		}

		var arg timecardParam
		err := c.BindQuery(&arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		hours, _ := db.GetHours(session.User.ID)
		if arg.Hours > 24 || hours > int(session.User.HourLimit.Int32) {
			log.Printf("invalid hours %v\n", hours)
			c.HTML(http.StatusOK, "timecard.html", gin.H{
				"Committed":    false,
				"showUsername": session.User.Name,
				"prefix":       prefix,
				"Exceeded":     true,
			})
			return
		}

		err = db.UpdateTimecard(session.User.ID, arg.Charge, arg.Hours, arg.Date)
		log.Printf("%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
	})

	r.Run(":9999")
}
