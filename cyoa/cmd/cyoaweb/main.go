package main

import (
	"cyoa"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// brief server welcome text
	fmt.Println("Welcome to Choose Your Own Adventure")

	// get flags
	filename := flag.String("f", "gopher.json", "name of CYOA JSON file")
	port := flag.Int("p", 3000, "desired HTTP server port")
	flag.Parse()

	// read in JSON file containing story data
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cyoa: error opening file %s: %v\n", *filename, err)
	}

	// parse JSON
	s, _ := cyoa.ParseJson(data)

	tmpl := template.Must(template.New("").Parse(storyTmpl))
	// Create new HTTP handler and serve
	h := cyoa.NewHandler(s,
		cyoa.WithTemplate(tmpl),
		cyoa.WithPathFunc(pathFn),
	)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </head>
  <body>
    <section class="page">
     <h1>{{.Title}}</h1>
     {{range .Story}}
       <p>{{.}}</p>
     {{end}}

     <ul>
     {{range .Options}}
       <li><a href="/story/{{.Arc}}">{{.Text}}</a></li> 
     {{end}}
     </ul>
   </body>
  </section>
</html>
`
