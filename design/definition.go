package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("amflow", func() {
	Title("amflow")
	Description("Archivematica workflow editor")
	Docs(func() {
		Description("amflow README")
		URL("https://github.com/sevein/amflow")
	})
	Host("localhost")
	Scheme("http")
	BasePath("/")
	Origin("http://swagger.goa.design", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		MaxAge(600)
		Credentials()
	})
	Consumes("application/json")
	Produces("application/json", "application/xml")
})
