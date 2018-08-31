package haha


import(

	"net/http"
	"strings"
	"fmt"
	
)



func sayHello(w http.ResponseWriter, r *http.Request){


	message:=r.URL.Path
	message=strings.TrimPrefix(message,"/")
	message="Hello "+message
	w.Write([] byte(message))
}

func ping (w http.ResponseWriter,r *http.Request){
	w.Write([] byte ("pong"))
}


func sayhelloName(w http.ResponseWriter,r *http.Request){
	
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path",r.URL.Path)
	fmt.Println("scheme",r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k,v :=range r.Form{
		fmt.Println("key:",k)
		fmt.Println("val:",strings.Join(v,""))
	}
	fmt.Fprintf(w,
		"Hello")
}
func main(){
	
	http.HandleFunc("/",sayhelloName)
	if err:=http.ListenAndServe(":8080",nil); err!=nil{
		panic(err)
	}
	/*
	http.Handle("/",http.FileServer(http.Dir("./src")))
	http.HandleFunc("/ping",ping)
	if err:=http.ListenAndServe(":8080",nil); err!=nil{
		panic(err)
	}*/
}
