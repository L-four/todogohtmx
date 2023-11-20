package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)
import _ "embed"

//go:embed static/css/pico.min.css
var static embed.FS

//go:embed template
var templateFs embed.FS

//go:embed db
var dbFs embed.FS

var homeTpl *template.Template
var todoTmpl *template.Template

type List struct {
	Id     uint64
	Title  string
	SortBy string
}

type TodosPage struct {
	Todos []Todo
	State url.Values
}

func main() {
	CheckAndSetupDb(dbFs)
	buildTemplates()
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/api/todo/", TodoHandler)
	http.Handle("/static/", http.FileServer(http.FS(static)))
	log.Fatal(http.ListenAndServe(":80", nil))
}

func buildTemplates() {
	homeTpl = template.Must(template.New("home").ParseFS(templateFs, "template/base/home.gohtml"))
	todoTmpl = template.Must(template.Must(homeTpl.Clone()).ParseFS(templateFs, "template/pages/todos.gohtml"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	todos := make(TodoSlice, 50)
	count := todos.LoadAll()
	todos = todos[:count]
	todoPage := TodosPage{
		Todos: todos,
		State: getState(r),
	}

	if todoPage.State.Get("list_id") == "" {
		todoPage.State.Set("list_id", "1")
	}

	err := todoTmpl.ExecuteTemplate(w, "home.gohtml", todoPage)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}

const prefix = "s_"

func getState(r *http.Request) url.Values {
	if r.PostForm == nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatalf("Formparse fail %s", err)
		}
	}
	state := url.Values{}
	for i, v := range r.PostForm {
		if i[0:2] == prefix {
			state.Set(i[2:], v[0])
		}
	}
	return state
}

func getIdAt(idx int, str string) (int, error) {
	parts := strings.Split(str, "/")
	id, err := strconv.Atoi(parts[idx])
	if err != nil {
		return 0, err
	}
	return id, nil
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		state := getState(r)
		title := r.PostFormValue("title")
		if len(title) == 0 {
			http.Error(w, "Title is Empty", http.StatusBadRequest)
			return
		}
		list_id := state.Get("list_id")
		if list_id == "" {
			http.Error(w, "No list to add the todo to", http.StatusBadRequest)
			return
		}
		parsed_list_id, err := strconv.Atoi(list_id)
		if err != nil {
			log.Fatalf("Invalid list_id:  %s", err)
		}
		todo := Todo{
			Title:     title,
			ListId:    uint64(parsed_list_id),
			Completed: false,
		}
		err = todo.Insert()
		if err != nil {
			log.Fatalf("db insert: %s", err)
		}

		err = todoTmpl.ExecuteTemplate(w, "todo", todo)
		if err != nil {
			log.Fatalf("template execution: %s", err)
		}
		break
	case http.MethodPatch:
		id, err := getIdAt(3, r.RequestURI)
		log.Printf("PATCH %d", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		completed := r.PostFormValue("completed")
		// if it's not set that means it's false as the html form won't send it.
		todo := Todo{
			Id:        uint64(id),
			Completed: strings.ToLower(completed) == "on",
		}
		err = todo.UpdateCompleted()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		break
	case http.MethodDelete:
		id, err := getIdAt(3, r.RequestURI)
		log.Printf("DELETE %d", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todo := Todo{
			Id: uint64(id),
		}
		err = todo.Delete()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		break
	default:
		http.Error(w, "DID you mean POST, PATCH or DELETE?", http.StatusBadRequest)
		break
	}
}
