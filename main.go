package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personalweb/connection"
	"personalweb/middleware"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id              int
	NameProject     string
	StartDate       time.Time
	EndDate         time.Time
	Description     string
	Image           string
	Duration        string
	FormatStartDate string
	FormatEndDate   string
	Author			string
	Html            bool
	Css             bool
	ReactJs         bool
	JavaScript      bool
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var userData = SessionData{}



func main() {
	connection.DatabaseConnect()

	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	e.Static("/public", "public")
	e.Static("/uploads", "uploads")
	
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/form-add-project", formAddProject)
	e.GET("/project-detail/:id", projectDetail)
	e.GET("/testimonial", testimonial)
	e.GET("/update-project/:id", updateProjectEdit)
	e.GET("/form-login", formLogin)
	e.GET("/form-register", formRegister)

	e.POST("/login", login)
	e.POST("/register", register)

	e.POST("/logout", logout)

	e.POST("/add-project", middleware.UploadFile(addProject))
	e.POST("/update-project/:id",  middleware.UploadFile(updateProject))
	e.POST("/delete-project/:id", deleteProject)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, tb_project.name, start_date, end_date, description, html, css, reactjs, javascript, image, tb_user.name AS author FROM tb_project JOIN tb_user ON tb_project.author_id = tb_user.id")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.NameProject, &each.StartDate, &each.EndDate, &each.Description, &each.Html, &each.Css, &each.ReactJs, &each.JavaScript, &each.Image, &each.Author)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}
		each.FormatStartDate = each.StartDate.Format("2006-01-02")
		each.FormatEndDate = each.EndDate.Format("2006-01-02")
		each.Duration = calculateDuration(each.FormatStartDate, each.FormatEndDate)

		result = append(result, each)
	}

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"Projects":     result,
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"DataSession":  userData,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}


func formAddProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func calculateDuration(startDate string, endDate string) string {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	durationTime := int(endTime.Sub(startTime).Hours())
	durationDays := durationTime / 24
	durationWeeks := durationDays / 7
	durationMonths := durationWeeks / 4
	durationYears := durationMonths / 12

	var duration string

	if durationYears > 1 {
		duration = strconv.Itoa(durationYears) + " years"
	} else if durationYears > 0 {
		duration = strconv.Itoa(durationYears) + " year"
	} else {
		if durationMonths > 1 {
			duration = strconv.Itoa(durationMonths) + " months"
		} else if durationMonths > 0 {
			duration = strconv.Itoa(durationMonths) + " month"
		} else {
			if durationWeeks > 1 {
				duration = strconv.Itoa(durationWeeks) + " weeks"
			} else if durationWeeks > 0 {
				duration = strconv.Itoa(durationWeeks) + " week"
			} else {
				if durationDays > 1 {
					duration = strconv.Itoa(durationDays) + " days"
				} else {
					duration = strconv.Itoa(durationDays) + " day"
				}
			}
		}
	}

	return duration
}

func addProject(c echo.Context) error {
	nameProject := c.FormValue("inputProjectName")
	startDate := c.FormValue("inputStartDate")
	endDate := c.FormValue("inputEndDate")
	description := c.FormValue("inputDesc")

	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	var html bool
	if c.FormValue("html") == "checked" {
		html = true
	}

	var css bool
	if c.FormValue("css") == "checked" {
		css = true
	}

	var reactJs bool
	if c.FormValue("react") == "checked" {
		reactJs = true
	}

	var js bool
	if c.FormValue("js") == "checked" {
		js = true
	}

	image := c.Get("dataFile").(string)

	sess, _ := session.Get("session", c)
	author := sess.Values["id"].(int)

	_, err := connection.Conn.Exec(context.Background(), 
	"INSERT INTO tb_project(name, start_date, end_date, description, html, css, reactjs, javascript, image, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
	 nameProject, startTime, endTime, description, html, css, reactJs, js, image, author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func testimonial(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimonial.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, html, css, reactjs, javascript, image FROM tb_project WHERE id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.NameProject, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Html, &ProjectDetail.Css, &ProjectDetail.ReactJs, &ProjectDetail.JavaScript, &ProjectDetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	ProjectDetail.FormatStartDate = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.FormatEndDate = ProjectDetail.EndDate.Format("2006-01-02")
	ProjectDetail.Duration = calculateDuration(ProjectDetail.FormatStartDate, ProjectDetail.FormatEndDate)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errTemplate = template.ParseFiles("views/project-detail.html")

	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Id : ", id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/#my-project")
}



func updateProjectEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var dataProject = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, html, css, reactjs, javascript, image FROM tb_project WHERE id=$1", id).Scan(
		&dataProject.Id, &dataProject.NameProject, &dataProject.StartDate, &dataProject.EndDate, &dataProject.Description, &dataProject.Html, &dataProject.Css, &dataProject.ReactJs, &dataProject.JavaScript, &dataProject.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	data := map[string]interface{}{
		"Project": dataProject,
	}

	var tmpl, errTemplate = template.ParseFiles("views/update-project.html")

	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	nameProject := c.FormValue("inputProjectName")
	startDate := c.FormValue("inputStartDate")
	endDate := c.FormValue("inputEndDate")
	description := c.FormValue("inputDesc")

	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	var html bool
	if c.FormValue("html") == "checked" {
		html = true
	}

	var css bool
	if c.FormValue("css") == "checked" {
		css = true
	}

	var reactJs bool
	if c.FormValue("react") == "checked" {
		reactJs = true
	}

	var js bool
	if c.FormValue("js") == "checked" {
		js = true
	}

	image := c.Get("dataFile").(string)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name=$1, start_date=$2, end_date=$3, description=$4, html=$5, css=$6, reactjs=$7, javascript=$8, image=$9 WHERE id=$10",
		nameProject, startTime, endTime, description, html, css, reactJs, js, image, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/#my-project")
}

func formRegister(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error { 
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("inputName")
	email := c.FormValue("inputEmail")
	password := c.FormValue("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		redirectWithMessage(c, "Register failed, please try again.", false, "/form-register")
	}

	return redirectWithMessage(c, "Register success!", true, "/form-login")
}


func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	email := c.FormValue("inputEmail")
	password := c.FormValue("inputPassword")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "Email Incorrect!", false, "/form-login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Incorrect!", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 JAM
	sess.Values["message"] = "Login success!"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}
