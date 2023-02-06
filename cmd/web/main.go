package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	//flag name+default value+ help snippet
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	//Custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}
	//go run cmd/web/* >>/tmp/info.log 2>>/tmp/error.log to dump logs

	//custom http server struct to use custom loggers and addresses
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	//Using command line flag for specific port
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

// curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static/

//FOR LOGGING TO FILES
// f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666) 
// if err != nil {
// log.Fatal(err)
// }
// defer f.Close()
// infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)