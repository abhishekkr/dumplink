package main

import (
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"
  "runtime"
  "time"
  "io/ioutil"
)

type Links struct {
  Urls []string
}

var linkfile = "./public/static/links.json"

func not_in_list(new_item string, list []string) bool{
  for _, item := range list {
    if item == new_item {
      return false
    }
  }
  return true
}

func save_json(new_url string) bool{
  contents,_ := ioutil.ReadFile(linkfile)
  var dump Links
  err := json.Unmarshal([]byte(contents), &dump)
  if err != nil {
    fmt.Println("Error:", err)
    return false
  }

  if not_in_list(new_url, dump.Urls) {
    dump.Urls = append(dump.Urls, new_url)
  } else {
    return true
  }
  b, _ := json.Marshal(&dump)
  ioutil.WriteFile(linkfile, b, 0x644)
  return true
}

func DumpLink(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  t, _ := template.ParseFiles("public/bookmarklet.html")
  t.Execute(w, nil)
}

func Dump(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  req.ParseForm()
  fmt.Println(req.FormValue("url"))
  fmt.Println(req.FormValue("title"))
  save_json(req.FormValue("url"))

  t, _ := template.ParseFiles("public/status.html")
  t.Execute(w, nil)
}

func GetLinks(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  t, _ := template.ParseFiles("public/status.html")
  t.Execute(w, nil)
}


func main(){
  runtime.GOMAXPROCS(runtime.NumCPU())
  httpport := 8000

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))
  http.HandleFunc("/", DumpLink)
  http.HandleFunc("/dumplink", DumpLink)

  http.HandleFunc("/dump", Dump)
  http.HandleFunc("/links", GetLinks)

  srv := &http.Server{
    Addr:        fmt.Sprintf(":%d", httpport),
    Handler:     http.DefaultServeMux,
    ReadTimeout: time.Duration(5) * time.Second,
  }

  fmt.Printf("access your goshare at http://<IP>:%d\n", httpport)
  srv.ListenAndServe()
}
