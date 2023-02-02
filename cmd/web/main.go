package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	//flag name+default value+ help snippet
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	//Custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//For serveing static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//go run cmd/web/* >>/tmp/info.log 2>>/tmp/error.log to dump logs

	//custom http server struct to use custom loggers and addresses
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	//Using command line flag for specific port
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

// curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static/
