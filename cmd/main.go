package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	Templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
func newTemplate() *Templates {
	return &Templates{
		Templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
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

	e.Logger.Fatal(e.Start(":8080"))
}
