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

type Rule struct {
	Database    string `yaml:"database"`
	Querycustom string `yaml:"querycustom"`
	Timeout     string `yaml:"timeout"`
	Help        string `yaml:"help"`
}

type Query struct {
	Rules_query []Rule `yaml:"rules_query"`
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

	ir := "Test"
	//  db := "ere"

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting web server at %s\n", *portWeb)
	for i := 0; i < getInfoCount(*query); i++ {
		go run(*host, *port, *user, *password, getInfoQuery(*query, "dbname", i), getInfoQuery(*query, "queryCustom", i), creatGauge(getInfoQuery(*query, "dbname", i), ir))
		ir = "test"
		//db = "dsfg"
	}
	log.Fatal(http.ListenAndServe(*portWeb, nil))

}

func run(host string, port string, user string, password string, dbname string, queryCustom string, gauge prometheus.Gauge) {
	for {
		gauge.Set(connectDB(host, port, user, password, dbname, queryCustom))
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func creatGauge(dbname string, help string) prometheus.Gauge {
	outSql := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: dbname,
		Help: help,
	})

	prometheus.MustRegister(outSql)
	return outSql
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

//получаем бд из query.yml

func getInfoQuery(query string, check string, icheck int) string {
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
	if check == "dbname" {
		v = config.Rules_query[icheck].Database
		if icheck > 1 {
			v += string(icheck)
		}
	} else if check == "queryCustom" {
		v = config.Rules_query[icheck].Querycustom
	} else if check == "help" {
		v = config.Rules_query[icheck].Help
	}

	return v
}

func getInfoCount(query string) int {

	var v int
	var config Query
	source, err := ioutil.ReadFile(query)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	v = len(config.Rules_query)

	return v

}
