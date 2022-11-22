package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Bid struct {
	CarId     string
	CarEnergy int
	CarRadius float32
	CarLat    float64
	CarLon    float64
	Price     int
	TokenId   string
	Win       bool `default:false`
}

type Auction struct {
	WinnerCarId string
	TokenId     string
}

type Token struct {
	TokenId    string
	TokenLat   float64
	TokenLon   float64
	TokenPrice int
}

func bidHandle(w http.ResponseWriter, r *http.Request) {
	// Add new bid to db
	var b Bid
	json.NewDecoder(r.Body).Decode(&b)
	fmt.Printf("%v\n", b)
	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO bids(carId, carEnergy, carRadius, carLat, carLon, price, tokenId, win) values(?, ?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)

	res, err := stmt.Exec(b.CarId, b.CarEnergy, b.CarRadius, b.CarLat, b.CarLon, b.Price, b.TokenId, b.Win)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	db.Close()
}

func auctionHandle(w http.ResponseWriter, r *http.Request) {
	// make bid to winner and remove all other bids linked with this token
	var a Auction
	json.NewDecoder(r.Body).Decode(&a)
	fmt.Printf("%v\n", a)

	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	// update
	stmt, err := db.Prepare("UPDATE bids SET win=? WHERE carId=?")
	checkErr(err)

	res, err := stmt.Exec(true, a.WinnerCarId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("winner : ")
	fmt.Println(affect)

	// delete
	stmt, err = db.Prepare("DELETE FROM bids WHERE tokenId=? AND win=0")
	checkErr(err)

	res, err = stmt.Exec(a.TokenId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println("losers : ")
	fmt.Println(affect)

	db.Close()
}

func tokenHandle(w http.ResponseWriter, r *http.Request) {
	// Add a new token to db
	var t Token
	json.NewDecoder(r.Body).Decode(&t)
	fmt.Printf("%v\n", t)

	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO tokens(tokenId, tokenLat, tokenLon, tokenPrice) values(?, ?, ?, ?)")
	checkErr(err)

	res, err := stmt.Exec(t.TokenId, t.TokenLat, t.TokenLon, t.TokenPrice)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	db.Close()
}

func initDB(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	// flush all if exists
	stmt, err := db.Prepare("DROP TABLE IF EXISTS bids")
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("DROP TABLE IF EXISTS tokens")
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	// create table
	sqlStmt := `
	create table bids (id integer not null primary key, carId text, carEnergy integer, carRadius real, carLat real, carLon real, price integer, tokenId text, win boolean);
	delete from bids;
	`
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	sqlStmt = `
	create table tokens (id integer not null primary key, tokenId text, tokenLat real, tokenLon real, tokenPrice integer);
	delete from tokens;
	`
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
