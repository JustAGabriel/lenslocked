package controllers

import (
	"html/template"
	"net/http"

	"github.com/justagabriel/lenslocked/models"
)

type TemplateBaseData struct {
	User models.User
}

func StaticHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := TemplateBaseData{}
		user := GetUserFromContext(r.Context())
		if user != nil {
			data.User = *user
		}

		err := tpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

type QuestionAnswer struct {
	Question string
	Answer   template.HTML
}

func FAQ(tpl *template.Template) http.HandlerFunc {
	questions := []QuestionAnswer{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
