package service

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"html/template"
)
type Information struct {
	Name    string
	Phone	string
	Date	string
}
type Temp struct {
	All	[]Information
}

func jsHandler(formatter *render.Render) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			formatter.JSON(w, http.StatusOK, struct {
				ID      string `json:"id"`
				Content string `json:"content"`
			}{ID: "020-88888888", Content: "Guangzhou Xingangdong"})
		}
		
	}

func unknownHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "501 Not Implemented", http.StatusNotImplemented)	
}

func reserve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./template/information.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		r.ParseForm()
	
		issueList := template.Must(template.New("issuelist").Parse(`
			<table>
			<tr style='text-align: left'>
			  <th>Name</th>
				<th>Phone</th>
				<th>Date</th>
			</tr>
			
			{{range .All}}
			<tr>
					<td>{{.Name}}</td>
					<td>{{.Phone}}</td>
					<td>{{.Date}}</td>
			</tr>
			{{end}}
			</table>
		`))

		ns := r.Form["name"]
		ps := r.Form["phone"]
		ds := r.Form["date"]

		var size = len(ns)
		println(size)
		var temp Temp
    		for i := 0; i < size; i++ {
			temp.All = append(temp.All, Information{Name: ns[i], Phone:ps[i], Date:ds[i]})

		}
		/*{ID: r.Form["username"], Content: r.Form["password"]}*/
		if err := issueList.Execute(w, temp); err != nil {
			log.Fatal(err)
		}
	}
}


// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			//println(root)
			panic("Could not retrive working directory")
		} else {
			webRoot = root
			fmt.Println(root)
		}
	}

	mx.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/assets/"))))
	mx.HandleFunc("/", jsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/unknown", unknownHandler)
	mx.HandleFunc("/reserve", reserve)
}
