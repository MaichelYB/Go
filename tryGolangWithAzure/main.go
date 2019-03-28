package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

var server = "maichelsqldata.database.windows.net"
var port = 1433
var user = "maichelyb"
var password = "Mai6chel"
var database = "maicheldata"

var templ *template.Template

type Biodata struct {
	ID   int
	Nama string
	Job  string
	Date time.Time
}

func main() {
	logFile := "testlogfile"
	port := "3001"
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		logFile = "D:\\home\\site\\wwwroot\\testlogfile"
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	DBConnection()

	http.HandleFunc("/", index)

	http.HandleFunc("/read", ReadUsers)
	if err != nil {
		http.ListenAndServe(":"+port, nil)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("--->   UP @ " + port + "  <------")
	}
	//http.ListenAndServe(":8800", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	// tmpl := template.Must(template.ParseFiles("template/welcome.html"))
	tmpl := template.Must(template.ParseFiles("D:\\home\\site\\wwwroot\\template\\welcome.html"))
	if req.Method != "POST" {
		tmpl.Execute(res, nil)
		return
	}

	details := Biodata{
		Nama: req.FormValue("nama"),
		Job:  req.FormValue("job"),
	}
	_ = details

	createID, err := CreateUsers(req.FormValue("name"), req.FormValue("job"))

	if err != nil {
		log.Fatal("Error Creating Employees", err.Error())
	}
	fmt.Printf("Inserted ID: %d successfully.\n", createID)
	tmpl.ExecuteTemplate(res, "welcome.html", nil)
}

//for connecting to DB Azure
func DBConnection() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func get2DigitNumber(number int) string {
	ret := ""
	if number < 10 {
		ret = "0" + strconv.Itoa(number)
	} else {
		ret = strconv.Itoa(number)
	}
	return ret
}

//For Insert to employee
func CreateUsers(name string, job string) (int64, error) {
	ctx := context.Background()
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	dateString := strconv.Itoa(year) + "-" + get2DigitNumber(int(month)) + "-" + get2DigitNumber(day)
	fmt.Println(dateString)
	var err error

	if db == nil {
		err = errors.New("Create Employee: db is null")
		return -1, err
	}

	//check database if alive
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := "INSERT INTO TestSchema.Users(Name, Job, Date) VALUES (@name, @job, '" + dateString + "');select convert(bigint, SCOPE_IDENTITY());"

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Name", name),
		sql.Named("Job", job),
		sql.Named("Date", dateString))

	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}

// ReadEmployees reads all employee records
func ReadUsers(res http.ResponseWriter, req *http.Request) {
	//tmpl := template.Must(template.ParseFiles("template/tableTemplate.html"))
	tmpl := template.Must(template.ParseFiles("D:\\home\\site\\wwwroot\\template\\tableTemplate.html"))
	// Execute query
	rows, err := db.Query("SELECT * FROM TestSchema.Users;")
	if err != nil {

	}

	defer rows.Close()

	var name string
	var job string
	var date time.Time
	var id int
	var b []Biodata

	// Iterate through the result set.
	for rows.Next() {

		// Get values from row.
		err = rows.Scan(&id, &name, &job, &date)
		b = append(b, Biodata{ID: id, Nama: name, Job: job, Date: date})

	}
	type PageData struct {
		PageTitle string
		Biodatas  []Biodata
	}
	tmpl.ExecuteTemplate(res, "tableTemplate.html", PageData{
		PageTitle: "People Who Already Registered",
		Biodatas:  b,
	})
}
