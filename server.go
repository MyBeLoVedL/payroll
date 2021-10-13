package main

import (
	"context"
	"fmt"
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

	r.GET("/profile", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)

		arg, err := db.SelectEmployee(int64(session.User.ID))
		if err != nil {
			if err != nil {
				c.HTML(http.StatusOK, "error_page.html", gin.H{
					"Msg": "no user when selecting employee",
				})
				return
			}
		}

		c.HTML(http.StatusOK, "profile.html", gin.H{
			"EmpID": arg.ID,
			"Name":  arg.Name,
			"Etype": arg.Type,
			"Mail":  arg.Mail,
			"S":     arg.SocialSecurityNumber,
			"Tax":   arg.StandardTaxDeductions,
			"Other": arg.OtherDeductions,
			"Phone": arg.PhoneNumber,
			"Rate":  arg.SalaryRate,
		})

	})

	r.GET("/updateProfile", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, err := misc.GSS.Get(sid)
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": "NO session",
			})
			return
		}

		emp, err := db.SelectEmployee(session.User.ID)
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": "no user when selecting employee",
			})
			return
		}

		c.HTML(http.StatusOK, "update_profile.html", gin.H{
			"Name":     emp.Name,
			"Password": emp.Password,
			"EmpID":    emp.ID,
		})
	})

	r.GET("/updateProfileAction", func(c *gin.Context) {
		type AddArg struct {
			EmpID    int64  `form:"empID" binding:"required"`
			Name     string `form:"name" binding:"required"`
			Password string `form:"password" binding:"required"`
		}
		var arg AddArg
		err := c.BindQuery(&arg)
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": fmt.Sprintf("update employee not enough arguments: %v", err.Error()),
			})
			return
		}
		log.Printf("update arg %+v\n", arg)
		err = db.UpdateNamePassword(arg.EmpID, arg.Name, arg.Password)
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": fmt.Sprintf("update profile failed due to server error : %v", err.Error()),
			})
			return
		}
		c.HTML(http.StatusOK, "main.html", nil)
	})

	r.GET("/cookie", func(c *gin.Context) {
		cookie, _ := c.Cookie("sid")
		c.JSON(http.StatusOK, gin.H{
			"session id": cookie,
		})

	})

	r.GET("/loginForm", func(c *gin.Context) {
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

	r.GET("/updatePay", func(c *gin.Context) {
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

	r.GET("/timecard", func(c *gin.Context) {

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

		err = db.UpdateTimecard(session.User.ID, arg.Charge, arg.Hours, arg.Date)
		log.Printf("%v\n", err)
		c.HTML(http.StatusOK, "timecard.html", gin.H{
			"Committed":    false,
			"showUsername": session.User.Name,
			"prefix":       prefix,
			"Exceeded":     false,
			"Projects":     misc.GetProjects(),
		})
	})

	r.GET("/showOrder", func(c *gin.Context) {

		c.HTML(http.StatusOK, "add_order.html", gin.H{
			"Products": misc.GetProducts(),
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
		orderID, err := db.AddOrder(context.Background(), db.AddOrderParams{
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

		c.HTML(http.StatusOK, "new_order.html", gin.H{
			"OrderID":     orderID,
			"Contact":     arg.Contact,
			"Address":     arg.Address,
			"ProductName": arg.ProductID,
			"Amount":      arg.Amount,
			"Date":        arg.Date.String(),
		})

	})

	r.GET("/updateOrder", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order_info.html", gin.H{
			"Orders": []Order{},
		})
	})

	r.GET("/selectOrder", func(c *gin.Context) {
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

		order, err := db.SelectOrderByID(arg.OrderID)
		if err != nil {
			c.HTML(http.StatusOK, "order_info.html", gin.H{
				"Display": false,
			})
			return
		}

		var closed string
		if order.Closed == 1 {
			closed = "closed"
		} else {
			closed = "active"
		}

		pros := misc.GetProducts()
		var projectName string
		for _, pro := range pros {
			if pro.ProductID == int(order.ProductID) {
				projectName = pro.ProductName
			}
		}
		log.Printf("%+v \n order ID %+v ", pros, order.ProductID)

		c.HTML(http.StatusOK, "order_info.html", gin.H{
			"Display":     true,
			"ID":          order.ID,
			"Contact":     order.Contact,
			"Address":     order.Address,
			"Amount":      order.Amount,
			"Date":        order.Date.String(),
			"ProductName": projectName,
			"Closed":      closed,
			"Products":    misc.GetProducts(),
		})
	})

	r.GET("/updateOrderAction", func(c *gin.Context) {
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

	r.GET("/deleteOrder", func(c *gin.Context) {
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
			"ID":       id,
			"Contact":  order.Contact,
			"Address":  order.Address,
			"Amount":   order.Amount,
			"Date":     order.Date.String(),
			"Products": misc.GetProducts(),
		})
	})

	r.GET("/manageEmployee", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)
		if session.User.Root.Int32 != 1 {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": "您并非管理员，无法进行此项操作",
			})
			return
		}
		c.HTML(http.StatusOK, "employee_info.html", gin.H{})
	})

	r.GET("/disReport", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)
		if session.User.Root.Int32 != 1 {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": "您并非管理员，无法进行此项操作",
			})
			return
		}
		c.HTML(http.StatusOK, "admin_report.html", nil)
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
			"Name":  emp.Name,
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
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": fmt.Sprintf("update employee not enough arguments: %v", err.Error()),
			})
			return
		}

		arg, err := db.SelectEmployee(int64(id))
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": fmt.Sprintf("update employee not enough arguments: %v", err.Error()),
			})
			return
		}
		c.HTML(http.StatusOK, "update_employee.html", gin.H{
			"empID":    arg.ID,
			"Etype":    arg.Type,
			"Mail":     arg.Mail,
			"Security": arg.SocialSecurityNumber,
			"Tax":      arg.StandardTaxDeductions,
			"Other":    arg.OtherDeductions,
			"Phone":    arg.PhoneNumber,
			"Rate":     arg.SalaryRate,
		})
	})

	r.GET("updateEmployeeAction", func(c *gin.Context) {
		type AddArg struct {
			EmpID    int64  `form:"empID" binding:"required"`
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
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": fmt.Sprintf("update employee not enough arguments: %v", err.Error()),
			})
			return
		}
		log.Printf("update arg %+v\n", arg)
		err = db.UpdateEmployee(arg.EmpID, arg.Etype, arg.Mail, arg.Security, arg.Tax, arg.Other, arg.Phone, arg.Rate)
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": fmt.Sprintf("update failed due to server error : %v", err.Error()),
			})
			return
		}
		c.HTML(http.StatusOK, "main.html", nil)
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

		if session.User.Root.Int32 != 1 {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": "您并非管理员，无法进行此项操作",
			})
			return
		}

		hours, _ := db.GetHoursByEmpID(session.User.ID)

		payYear, _ := db.GetPayYearToDate(session.User.ID)

		c.HTML(http.StatusOK, "employee_report.html", gin.H{
			"Hours":    hours,
			"PayYear":  payYear,
			"Vacation": DEFAULT_VACATION,
			"Projects": misc.GetProjects(),
			"Records":  []ReportRecord{},
		})
	})

	r.GET("chargeProject", func(c *gin.Context) {
		sid, _ := c.Cookie("sid")
		session, _ := misc.GSS.Get(sid)

		hours, _ := db.GetHoursByEmpID(session.User.ID)

		payYear, _ := db.GetPayYearToDate(session.User.ID)

		charge := c.Query("chargeNumber")
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

		projectHours, err := db.GetHoursByProject(session.User.ID, int64(num))
		if err != nil {
			c.HTML(http.StatusOK, "employee_report.html", gin.H{
				"Records":  []ReportRecord{},
				"Hours":    hours,
				"PayYear":  payYear,
				"Vacation": DEFAULT_VACATION,
				"Projects": misc.GetProjects(),
			})
			return
		}

		pros := misc.GetProjects()
		var projectName string
		for _, pro := range pros {
			if pro.ChargeNumber == num {
				projectName = pro.ProjectName
			}
		}
		c.HTML(http.StatusOK, "employee_report.html", gin.H{
			"Records": []ReportRecord{
				{projectHours, projectName},
			},
			"Hours":    hours,
			"PayYear":  payYear,
			"Vacation": DEFAULT_VACATION,
			"Projects": misc.GetProjects(),
		})
	})

	r.GET("/adminReport", func(c *gin.Context) {
		type ReportArg struct {
			Rtype     string `form:"type" binding:"required"`
			EmpName   string `form:"empName" binding:"required"`
			StartDate string `form:"starttime" binding:"required"`
			EndDate   string `form:"endtime" binding:"required"`
		}
		var arg ReportArg
		err := c.BindQuery(&arg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad report arg",
			})
		}

		ids, err := db.GetIDByName(arg.EmpName)
		log.Printf("IDs %v\n", ids)

		if arg.Rtype == "hours" {
			hours := make([]Record, 0)
			for _, id := range ids {
				hour, err := db.GetHoursByEmpID(id)
				if err == nil {
					hours = append(hours, Record{id, arg.EmpName, float64(hour)})
				}
			}
			c.HTML(http.StatusOK, "admin_report.html", gin.H{
				"Employees": hours,
			})

		} else if arg.Rtype == "pay" {
			pays := make([]Record, 0)
			for _, id := range ids {
				pay, err := db.GetPayYearToDate(id)
				if err == nil {
					pays = append(pays, Record{id, arg.EmpName, pay})
				}
			}
			c.HTML(http.StatusOK, "admin_report.html", gin.H{
				"Employees": pays,
			})

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "unknown report type",
			})
		}
	})

	r.GET("/runPayroll", func(c *gin.Context) {
		info, err := db.GetPayInfo()
		if err != nil {
			c.HTML(http.StatusOK, "error_page.html", gin.H{
				"Msg": "no user when selecting employee",
			})
			return
		}
		for _, emp := range info {
			empInfo, _ := db.SelectEmployee(emp.ID)

			switch emp.PaymentMethod {
			case "hour":
				hours, err := db.GetHoursByEmpID(emp.ID)
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "employee info not available",
					})
					return
				}
				rate, err := strconv.ParseFloat(empInfo.SalaryRate, 64)
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "invalid salary rate for employee",
					})
					return
				}
				err = db.CreatePaycheck(emp.ID, fmt.Sprintf("%f", float64(hours)*rate), time.Now(), time.Now())
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "create paycheck failed",
					})
					return
				}
			case "salaried":
				rate, err := strconv.ParseFloat(empInfo.SalaryRate, 64)
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "invalid salary rate for employee",
					})
					return
				}
				err = db.CreatePaycheck(emp.ID, fmt.Sprintf("%f", rate), time.Now(), time.Now())
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "create paycheck failed",
					})
					return
				}

			case "commissioned":
				total, err := db.GetAmountByID(emp.ID)
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "Get amount failed",
					})
					return
				}
				rate, err := strconv.ParseFloat(empInfo.SalaryRate, 64)
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "invalid salary rate for employee",
					})
					return
				}

				err = db.CreatePaycheck(emp.ID, fmt.Sprintf("%f", rate*total), time.Now(), time.Now())
				if err != nil {
					c.HTML(http.StatusOK, "error_page.html", gin.H{
						"Msg": "create paycheck failed",
					})
					return
				}
			}
			err = db.CommitCard(emp.ID)
			if err != nil {
				c.HTML(http.StatusOK, "error_page.html", gin.H{
					"Msg": "commit timecard failed",
				})
				return
			}
		}
	})
}

type ReportRecord struct {
	ProjectHours int
	ProjectName  string
}

const DEFAULT_VACATION = 10

type Record struct {
	EmpID   int64
	Name    string
	Numeric float64
}
