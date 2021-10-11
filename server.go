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

	r.GET("/updatePickupPay", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, err := misc.GSS.Get(sid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "No such session",
			})
		}
		db.UpdatePayment("pick_up", session.User.ID)
	})

	r.GET("/updateMailPay", func(c *gin.Context) {
		mail := c.Query("mail")
		// ! todo : error checking
		sid, _ := c.Cookie("sid")
		session, err := misc.GSS.Get(sid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "No such session",
			})
		}
		db.UpdatePaymentWIthMail(session.User.ID, mail)
	})

	r.GET("/updateDepositPay", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, err := misc.GSS.Get(sid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "No such session",
			})
		}
		bankName := c.Query("bankName")
		bankAccount := c.Query("bankAccount")
		log.Printf("name %v account %v\n", bankName, bankAccount)
		err = db.UpdatePaymentWithBank(session.User.ID, bankName, bankAccount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errror": err.Error(),
			})
		}

	})

	// r.Any("/updatePayAction", func(c *gin.Context) {
	// 	method := c.Query("updatePay")
	// 	sid, _ := c.Cookie("sid")
	// 	session, err := misc.GSS.Get(sid)
	// 	switch method {
	// 	case "pick_up":
	// 		db.UpdatePayment(method, session.User.ID)
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"error": err,
	// 		})
	// 	case "mail":
	// 		mail := c.PostForm("mail")
	// 		err = db.UpdatePaymentWIthMail(session.User.ID, mail)
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"error": err,
	// 		})
	// 	case "deposit":
	// 		bankName := c.PostForm("bankName")
	// 		bankAccount := c.PostForm("bankAccount")
	// 		err = db.UpdatePaymentWithBank(session.User.ID, bankName, bankAccount)
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"error": err,
	// 		})

	// 	}
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"error":  err,
	// 		"method": method,
	// 	})

	// })

	// r.Any("/timecard", func(c *gin.Context) {
	// 	sid, _ := c.Cookie("sid")
	// 	session, _ := misc.GSS.Get(sid)
	// 	card, err := db.SelectTimeCard(session.User.ID)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": "No matched timecard",
	// 		})
	// 		return
	// 	}
	// 	time := card.StartDate.Time.String()
	// 	var com string
	// 	if card.Committed.Int32 == 0 {
	// 		com = "not committed"
	// 	} else {
	// 		com = "committed"
	// 	}
	// 	c.HTML(http.StatusOK, "timecard.html", gin.H{
	// 		"ID":        card.ID,
	// 		"StartDate": time,
	// 		"Committed": com,
	// 	})

	// })

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
				"Committed":    false,
				"showUsername": session.User.Name,
				"prefix":       prefix,
			})
		}

		var arg timecardParam
		err := c.BindQuery(&arg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// hours, _ := db.GetHours(session.User.ID)
		// if arg.Hours > 24 || hours > int(session.User.HourLimit.Int32) {
		// 	log.Printf("invalid hours %v\n", hours)
		// 	c.HTML(http.StatusOK, "timecard.html", gin.H{
		// 		"Committed":    false,
		// 		"showUsername": session.User.Name,
		// 		"prefix":       prefix,
		// 		"Exceeded":     true,
		// 	})
		// 	return
		// }

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
		log.Printf("Order %+v\n", order)
		c.HTML(http.StatusOK, "order_info.html", gin.H{
			"ID":      order.ID,
			"Contact": order.Contact,
			"Address": order.Address,
			"Amount":  order.Amount,
			"Date":    order.Date.String(),
			"Products": []Product{
				{1111, "RTX 3090"},
				{2222, "GTX 1080"},
				{3333, "Intel i9 9900k"},
			},
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

		idStr := c.Query("orderID")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No such order ID",
			})
		}

		var arg UpdateOrderParams
		err = c.BindQuery(&arg)
		log.Printf("%+v\n", arg)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		// todo: verify order id belong to this user
		// sid, _ := c.Cookie("sid")
		// session, _ := misc.GSS.Get(sid)
		err = db.UpdateOrder(context.Background(), db.UpdateOrderParams{
			ID:        int64(id),
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

	r.GET("/updateOldOrder", func(c *gin.Context) {
		arg := c.Query("orderID")
		id, err := strconv.Atoi(arg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No such order",
			})
		}

		order, err := db.SelectOrderByID(int64(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No  order in DB",
			})
		}

		c.HTML(http.StatusOK, "update_order.html", gin.H{
			"ID":      id,
			"Contact": order.Contact,
			"Address": order.Address,
			"Amount":  order.Amount,
			"Date":    order.Date.String(),
			"Products": []Product{
				{1111, "RTX 3090"},
				{2222, "GTX 1080"},
				{3333, "Intel i9 9900k"},
			},
		})
	})

	r.GET("/manageEmployee", func(c *gin.Context) {
		c.HTML(http.StatusOK, "employee_info.html", gin.H{})
	})

	r.GET("/searchEmployee", func(c *gin.Context) {
		idStr := c.Query("empID")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "empty employee id",
			})
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "invalid emp id",
			})
		}
		emp, err := db.SelectEmployee(int64(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "No such employee in DB",
			})
		}
		log.Printf("emp name : %+v\n", emp)
		c.HTML(http.StatusOK, "employee_info.html", gin.H{
			"Name":  emp.Name.String,
			"empID": emp.ID,
		})
	})

	r.GET("/addEmployee", func(c *gin.Context) {
		type AddArg struct {
			Etype    string `form:"etype" binding:"required"`
			Mail     string `form:"mail" binding:"required"`
			Security string `form:"security" binding:"required"`
			Tax      string `form:"tax" binding:"required"`
			Other    string `form:"other" binding:"required"`
			Phone    string `form:"phone" binding:"required"`
			Rate     string `form:"rate" binding:"required"`
		}

		var arg AddArg

		err := c.BindQuery(&arg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "add employee bind fails",
			})
		}

		log.Printf("Server arg %+v\n", arg)
		id, err := db.AddEmployee(arg.Etype, arg.Mail, arg.Security, arg.Tax, arg.Other, arg.Phone, arg.Rate)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})

		}
		c.HTML(http.StatusOK, "new_employee.html", gin.H{
			"EmpID":    id,
			"Mail":     arg.Mail,
			"Etype":    arg.Etype,
			"Phone":    arg.Phone,
			"Tax":      arg.Tax,
			"Rate":     arg.Rate,
			"Other":    arg.Other,
			"Security": arg.Security,
		})
	})

	r.GET("/updateEmployee", func(c *gin.Context) {
		// sid, _ := c.Cookie("sid")
		// session, _ := misc.GSS.Get(sid)

		idStr := c.Query("empID")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No such employee to update",
			})
		}

		arg, err := db.SelectEmployee(int64(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No such employee to delete",
			})
		}

		c.HTML(http.StatusOK, "displayEmpoyee.main", gin.H{
			"ID":       arg.ID,
			"Etype":    arg.Type,
			"Mail":     arg.Mail,
			"Security": arg.SocialSecurityNumber,
			"Tax":      arg.StandardTaxDeductions,
			"Other":    arg.OtherDeductions,
			"Phone":    arg.PhoneNumber,
			"Rate":     arg.SalaryRate,
		})
	})

	r.GET("/deleteEmployee", func(c *gin.Context) {
		idStr := c.Query("empID")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("delete error:%v\n", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No such employee to delete",
			})
		}
		err = db.DeleteEmployee(int64(id))
		if err != nil {
			log.Printf("delete error:%v\n", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "delete employee failed",
			})
		}
	})

	r.GET("/employeeReport", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)

		reportType := c.Query("reportType")
		if reportType == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing report type",
			})
		}

		switch reportType {
		case "0":
			hours, err := db.GetHoursByEmpID(session.User.ID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "no hours found",
				})
			}
			c.String(http.StatusOK, "You have worded %v\n", hours)
		case "1":
			charge := c.Query("charge")
			if charge == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "missing charge number",
				})
			}
			num, err := strconv.Atoi(charge)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid charge number",
				})
			}
			hours, err := db.GetHoursByProject(session.User.ID, int64(num))
			c.String(http.StatusOK, "You have worded %v on project %v\n", hours, charge)

		case "2":
		case "3":
			hours, err := db.GetPayYearToDate(session.User.ID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "no hours found",
				})
			}
			c.String(http.StatusOK, "You have worded %v for this year\n", hours)
		}

	})

}

type Product struct {
	ProductID   int
	ProductName string
}
