package main

import (
	"html/template"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var id = 0
type Contact struct {
	Name  string
	Email string
	Id    int
}

func newContact(name, email string) Contact {
	id++
	return Contact{
		Name:  name,
		Email: email,
        Id: id,
	}
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func newData() Data {
	return Data{
		Contacts: Contacts{
			newContact("John", "jd@gmail.com"),
			newContact("Clara", "cd@gmail.com"),
		},
	}
}
func (d Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			println("email already exists", contact.Email, email)
			return true
		}
	}
	return false
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func sameFormDataWithErrMsg(name, email string) FormData {
	formData := newFormData()
	formData.Values["name"] = name
	formData.Values["email"] = email

	formData.Errors["email"] = "Email already exists"
	return formData
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate()
    e.Static("/images", "images")
    e.Static("/css", "css")

	page := newPage()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := sameFormDataWithErrMsg(name, email)
			return c.Render(422, "form", formData)
		}

		newContact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, newContact)

		c.Render(200, "form", newFormData())
		return c.Render(200, "oob-contact", newContact)
	})

    e.DELETE("/contacts/:id", func(c echo.Context) error {
        time.Sleep(3 * time.Second)
        strId := c.Param("id")
        id, err := strconv.Atoi(strId)
        if err != nil {
            return c.String(400, "Invalid id")
        }

        for i, contact := range page.Data.Contacts {
            if contact.Id == int(id) {
                page.Data.Contacts = append(page.Data.Contacts[:i], page.Data.Contacts[i+1:]...)
                break
            }
        }

        return c.NoContent(200)
    })


	e.Logger.Fatal(e.Start(":8080"))
}
