package main

import (
	"database/sql" // New import
	"flag"
"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PrathameshAnwekar/snippets/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	//flag name+default value+ help snippet
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	//Custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//OPEN DATABASE CONNECTION
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	templateCache , err := newTemplateCache("./ui/html")
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	//go run cmd/web/* >>/tmp/info.log 2>>/tmp/error.log to dump logs

	//custom http server struct to use custom loggers and addresses
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	defer db.Close()
	//Using command line flag for specific port
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

//FOR LOGGING TO FILES
// f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
// if err != nil {
// log.Fatal(err)
// }
// defer f.Close()
// infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
