package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func connectDB() (*sqlx.DB, error) {
	host := os.Getenv("ISUCONP_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("ISUCONP_DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUCONP_DB_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("ISUCONP_DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("ISUCONP_DB_PASSWORD")
	dbname := os.Getenv("ISUCONP_DB_NAME")
	if dbname == "" {
		dbname = "isuconp"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)

	return sqlx.Open("mysql", dsn)
}

type Post struct {
	ID           int       `db:"id"`
	Mime         string    `db:"mime"`
	Imgdata      []byte    `db:"imgdata"`
}

func main(){
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s.", err.Error())
	}
	defer db.Close()

	query := `select id, mime, imgdata
				from posts
				where id <= 10000
				order by id asc
				limit ?
				offset ?`
	stmt, err := db.Preparex(query)
	if err != nil {
		panic(err)
	}

	postCount := 10000
	pageNum := 500
	posts := []Post{}

	for i := 0; i < postCount/pageNum; i++ {
		stmt.Select(&posts, pageNum, i*pageNum)

		for j := 0; j < pageNum; j++ {
			err := ioutil.WriteFile(imageURL(&posts[j]), posts[j].Imgdata, 0666)
			if err != nil {
				panic(err)
			}
		}
	}
}

func imageURL(p *Post) string {
	ext := ""
	if p.Mime == "image/jpeg" {
		ext = ".jpg"
	} else if p.Mime == "image/png" {
		ext = ".png"
	} else if p.Mime == "image/gif" {
		ext = ".gif"
	}

	return "../../public/image/" + strconv.Itoa(p.ID) + ext
}