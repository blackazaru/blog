package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
	"net/http"
	"blog/models"
	"crypto/rand"
	"fmt"
	"github.com/martini-contrib/sessions"
	"github.com/russross/blackfriday"
	"html/template"
)
var postsCollection *mgo.Collection

func generateId() string{
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x",b)
}


func homeHandler(rnd render.Render) {
	rnd.HTML(200, "home", nil)
}

func aboutHandler(rnd render.Render)  {

	rnd.HTML(200, "about", nil)
}

func contactsHandler(rnd render.Render)  {

	rnd.HTML(200, "contacts", nil)
}



func blogHandler(rnd render.Render){
	postDocuments := []models.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)
	posts := []models.PostDocument{}
	for _, doc := range postDocuments {
		post := models.PostDocument{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMd}
		posts = append(posts, post)
	}

	rnd.HTML(200, "blog", posts)
}


func editHandler(rnd render.Render, params martini.Params, session sessions.Session){
	if session.Get("auth") == "OK" {
		id := params["id"]
		postDocument := models.PostDocument{}
		err := postsCollection.FindId(id).One(&postDocument)

		if err != nil {
			rnd.Redirect("/")
			return
		}

		post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMd}

		rnd.HTML(200, "write", post)
	}else{
		rnd.Redirect("/admin")
	}


}

func writeHandler(rnd render.Render, session sessions.Session){
	if session.Get("auth") == "OK" {
		rnd.HTML(200, "write", nil)
	}else{
		rnd.Redirect("/admin")
	}
}

func savePostHandler(rnd render.Render, r *http.Request){
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMd := r.FormValue("content")
	contentHtml := string(blackfriday.MarkdownBasic([]byte(contentMd)))


	if id != "" {
		postsCollection.UpdateId(id,&models.PostDocument{id, title, contentHtml, contentMd})
	}else{
		id = generateId()
		postsCollection.Insert(&models.PostDocument{id, title, contentHtml, contentMd})
	}
	rnd.Redirect("/")
}

func deleteHandler(rnd render.Render, params martini.Params , session sessions.Session){
	if session.Get("auth") == "OK" {
		id := params["id"]

		if id == "" {
			rnd.Redirect("/")
			return
		}

		postsCollection.RemoveId(id)

		rnd.Redirect("/")
	}else{
		rnd.Redirect("/admin")
	}
}

func adminHandler(rnd render.Render, session sessions.Session) {
	if session.Get("auth") != "OK" {
		rnd.HTML(200, "login", nil)
	}else{
		rnd.Redirect("/posts")
	}
}

func loginHandler(rnd render.Render, r *http.Request, session sessions.Session){
	login := r.FormValue("login")
	pass := r.FormValue("password")
	if login == "user" && pass == "pass" {
		session.Set("auth","OK")
		rnd.Redirect("/posts")
	}else{
		rnd.Redirect("/admin")
	}
}

func logoutHandler(rnd render.Render, session sessions.Session){
	if session.Get("auth") == "OK" {
		session.Delete("auth")
		rnd.Redirect("/")
	}
}

func postsHandler(rnd render.Render, session sessions.Session){
	if session.Get("auth") == "OK" {
		postDocuments := []models.PostDocument{}
		postsCollection.Find(nil).All(&postDocuments)
		posts := []models.Post{}
		for _, doc := range postDocuments {
			post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMd}
			posts = append(posts, post)
		}

		rnd.HTML(200, "posts", posts)
	}else{
		rnd.Redirect("/admin")
	}
}

func getHtmlHandler(rnd render.Render, r *http.Request){
	md := r.FormValue("md")
	htmlBytes := blackfriday.MarkdownBasic([]byte(md))
	rnd.JSON(200, map[string]interface{}{"html": string(htmlBytes)})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func postHandler(rnd render.Render, params martini.Params){
	id := params["id"]
	postDocument := models.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)

	if err != nil {
		rnd.Redirect("/")
		return
	}

	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMd}

	rnd.HTML(200, "post", post)
}

func main() {

	session, err := mgo.Dial("mongodb://user:pass@ds063439.mongolab.com:63439/heroku_app27487147")
	//session, err := mgo.Dial("localhost")
	if err != nil{
		panic(err)
	}

	postsCollection = session.DB("heroku_app27487147").C("posts")
	store := sessions.NewCookieStore([]byte("secret123"))
	m := martini.Classic()
	//staticOptions := martini.StaticOptions{ Prefix :"assets"}

	m.Use(martini.Static("assets"))

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	m.Use(sessions.Sessions("admin", store))


	m.Get("/",homeHandler)
	m.Get("/about",aboutHandler)
	m.Get("/contacts",contactsHandler)
	m.Get("/write",writeHandler)
	m.Get("/edit/:id",editHandler)
	m.Get("/delete/:id",deleteHandler)
	m.Post("/SavePost", savePostHandler)
	m.Get("/blog",blogHandler)
	m.Get("/admin",adminHandler)
	m.Post("/login",loginHandler)
	m.Get("/posts",postsHandler)
	m.Post("/logout", logoutHandler)
	m.Post("/gethtml", getHtmlHandler)
	m.Get("/post/:id", postHandler)


	m.Run()
}


