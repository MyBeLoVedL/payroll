package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	db "payroll/db/sqlc"
	"payroll/misc"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

	setRouter(r)

	r.Run(":9999")
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Order struct {
	ID        int64     `form:"empId"`
	Contact   string    `form:"contact"`
	Address   string    `form:"addres"`
	Date      time.Time `form:"date" time_format:"2006-01-02"`
	ProductID int64     `form:"productId"`
	Amount    string    `form:"amount"`
}

func setRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "No error",
		})
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
		if err := c.ShouldBindQuery(&info); err != nil {
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
		c.HTML(http.StatusOK, "change_method.html", nil)
	})

	r.Any("/updatePayAction", func(c *gin.Context) {
		method := c.Query("updatePay")
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
		c.HTML(http.StatusOK, "timecard.html", gin.H{
			"Committed":    false,
			"showUsername": session.User.Name,
			"prefix":       prefix,
			"Exceeded":     false,
			"Projects":     db.GetProjects(),
		})
	})

	r.Any("/showOrder", func(c *gin.Context) {

		type Product struct {
			ProductID   int
			ProductName string
		}

		c.HTML(http.StatusOK, "add_order.html", gin.H{
			"Products": []Product{
				{1111, "RTX 3090"},
				{2222, "GTX 1080"},
				{3333, "Intel i9 9900k"},
			},
		})
	})

	r.GET("/order", func(c *gin.Context) {
		type AddOrderParams struct {
			Contact   string    `form:"contact"`
			Address   string    `form:"address"`
			Date      time.Time `form:"date" time_format:"2006-01-02"`
			ProductID int64     `form:"productId"`
			Amount    string    `form:"amount"`
		}
		var arg AddOrderParams
		err := c.ShouldBindQuery(&arg)
		log.Printf("%+v\n", arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)
		err = db.AddOrder(context.Background(), db.AddOrderParams{
			EmpID:     session.User.ID,
			Contact:   arg.Contact,
			Address:   arg.Address,
			Date:      arg.Date,
			ProductID: arg.ProductID,
			Amount:    arg.Amount,
		})
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		c.HTML(http.StatusOK, "main.html", gin.H{
			"OK": true,
		})

	})

	r.Any("/updateOrder", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order_info.html", gin.H{
			"Orders": []Order{},
		})
	})

	r.Any("/selectOrder", func(c *gin.Context) {
		type Arg struct {
			OrderID   int64  `form:"orderID" binding:"required"`
			Contact   string `form:"contact"`
			ProductID int64  `form:"productID"`
			StartDate string `form:"startDate"`
			EndDate   string `form:"endDate"`
		}
		var arg Arg
		err := c.BindQuery(&arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		log.Printf("select order %+v\n", arg)

		order, _ := db.SelectOrderByID(arg.OrderID)
		c.HTML(http.StatusOK, "order_info.html", gin.H{
			"Orders": []db.Order{order},
		})

	})

	r.Any("/updateOrderAction", func(c *gin.Context) {
		type UpdateOrderParams struct {
			Contact   string    `form:"contact"`
			Address   string    `form:"address"`
			Date      time.Time `form:"date" time_format:"2006-01-02"`
			ProductID int64     `form:"productId"`
			Amount    string    `form:"amount"`
		}
		var arg UpdateOrderParams
		err := c.BindQuery(&arg)
		log.Printf("%+v\n", arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)
		err = db.UpdateOrder(context.Background(), db.UpdateOrderParams{
			ID:        session.User.ID,
			Contact:   arg.Contact,
			Address:   arg.Address,
			Date:      arg.Date,
			ProductID: arg.ProductID,
			Amount:    arg.Amount,
		})
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		c.HTML(http.StatusOK, "main.html", gin.H{
			"OK": true,
		})
	})

	r.Any("/deleteOrder", func(c *gin.Context) {
		// type delOrderArg struct {
		// 	orderID int64 `form:"orderID" binding:"required"`
		// }

		// var arg delOrderArg
		// err := c.BindQuery(&arg)
		arg := c.Query("orderID")
		id, err := strconv.Atoi(arg)
		log.Printf("order ID %+v\n", arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		err = db.DeleteOrder(int64(id))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		c.HTML(http.StatusOK, "order_info.html", gin.H{
			"Orders": []Order{},
		})
	})

}
