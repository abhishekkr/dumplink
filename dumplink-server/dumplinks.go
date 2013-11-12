package main

import (
  "fmt"
  "flag"
  "net/http"
  "html/template"
  "encoding/json"
  "runtime"
  "time"
  "io/ioutil"
)

type Link map[string]string
type Links struct {
  Urls []Link
}

var (
  httpuri   = flag.String("uri", "127.0.0.1", "what PORT to run the DumpLink daemon")
  httpport  = flag.String("port", "8100", "what PORT to run the DumpLink daemon")
  linkfile  = flag.String("json", "./public/static/links.json", "path to json file to save links, need {'url':[]}")
)

func not_in_list(new_link string, list []Link) bool{
  for _, item := range list {
    for link, _ := range item {
      if link == new_link { return false }
    }
  }
  return true
}

func read_json() string{
  contents, err := ioutil.ReadFile(*linkfile)
  if err != nil {
    fmt.Println("Error:", err)
    return "Error: JSON gave error while reading."
  }
  return string(contents)
}

func save_json(new_url string, title string) bool{
  contents,_ := ioutil.ReadFile(*linkfile)
  var dump Links
  err := json.Unmarshal([]byte(contents), &dump)
  if err != nil {
    fmt.Println("Error:", err)
    return false
  }
  var _link Link
  _link = make(Link)
  _link[new_url] = title

  if not_in_list(new_url, dump.Urls) {
    dump.Urls = append(dump.Urls, _link)
  } else {
    return true
  }
  b, _ := json.Marshal(&dump)
  ioutil.WriteFile(*linkfile, b, 0x644)
  return true
}

func DumpLink(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  t, _ := template.ParseFiles("public/bookmarklet.html")
  t.Execute(w, map[string]string {"HTTPURI": *httpuri, "HTTPPort": *httpport})
}

func Dump(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  req.ParseForm()
  fmt.Println(req.FormValue("url"))
  fmt.Println(req.FormValue("title"))
  save_json(req.FormValue("url"), req.FormValue("title"))

  t, _ := template.ParseFiles("public/dump200.html")
  t.Execute(w, nil)
}

func GetLinks(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  t, _ := template.ParseFiles("public/status.html")
  t.Execute(w, map[string]string {"BODY": read_json()})
}


func main(){
  runtime.GOMAXPROCS(runtime.NumCPU())
  flag.Parse()

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))
  http.HandleFunc("/", DumpLink)
  http.HandleFunc("/dumplink", DumpLink)

  http.HandleFunc("/dump", Dump)
  http.HandleFunc("/links", GetLinks)

  srv := &http.Server{
    Addr:        fmt.Sprintf("%s:%s", *httpuri, *httpport),
    Handler:     http.DefaultServeMux,
    ReadTimeout: time.Duration(5) * time.Second,
  }

  fmt.Printf("access your goshare at http://%s:%s\n", *httpuri, *httpport)
  err := srv.ListenAndServe()
  fmt.Println("Game Over:", err)
}
