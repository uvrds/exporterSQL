package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	nameQuery string = "database"
	help      string = "test"
	outSql           = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: nameQuery,
		Help: help,
	})
)

type Query struct {
	Database    string
	Querycustom string
}

func main() {

	portWeb := flag.String("listen-address", ":9493",
		"--listen-address=< > The address to listen on for HTTP requests.")
	query := flag.String("config", "./query.yml",
		"--query=< > The path to file query.")
	host := flag.String("host", "localhost",
		"--host=< > The server database postgres.")
	port := flag.String("port", "5432",
		"--port=< > The port database postgres.")
	user := flag.String("user", "postgres",
		"--user=< > The user database postgres.")
	password := flag.String("password", "",
		"--password=< > The password for user database postgres.")

	flag.Parse()

	dbname := getInfoQuery(*query, "dbname")
	nameQuery = dbname
	queryCustom := getInfoQuery(*query, "queryCustom")

	prometheus.MustRegister(outSql)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting web server at %s\n", *portWeb)
	go run(*host, *port, *user, *password, dbname, queryCustom)
	log.Fatal(http.ListenAndServe(*portWeb, nil))

}

func connectDB(host string, port string, user string, password string, dbname string, queryCustom string) float64 {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected database")

	var col1 float64
	rows, err := db.Query(queryCustom)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&col1)
		fmt.Println(col1)
	}
	defer rows.Close()

	return col1
}

func run(host string, port string, user string, password string, dbname string, queryCustom string) {
	for {
		outSql.Set(connectDB(host, port, user, password, dbname, queryCustom))
		time.Sleep(time.Duration(5) * time.Second)
	}
}

//получаем бд из query.yml
func getInfoQuery(query string, value string) string {
	var v string
	var config Query
	source, err := ioutil.ReadFile(query)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	if value == "dbname" {
		v = config.Database
	} else if value == "queryCustom" {
		v = config.Querycustom
	}
	return v
}
