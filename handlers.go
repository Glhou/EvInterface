package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	TokenId        string
	TokenLat       float64
	TokenLon       float64
	TokenPrice     int
	AuctionEndTime sql.NullInt64
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

	w.WriteHeader(http.StatusOK)
}

func auctionHandle(w http.ResponseWriter, r *http.Request) {
	// make bid to winner and remove all other bids linked with this token
	var a Auction
	json.NewDecoder(r.Body).Decode(&a)
	fmt.Printf("%v\n", a)

	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	if a.WinnerCarId == "-1" {
		// remove token, no winner
		stmt, err := db.Prepare("DELETE FROM tokens WHERE tokenId=?")
		checkErr(err)

		res, err := stmt.Exec(a.TokenId)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		fmt.Println(id)

	} else {
		// set time end of auction to now in seconds
		stmt, err := db.Prepare("UPDATE tokens SET auctionEndTime = ? WHERE tokenId = ?")
		checkErr(err)

		res, err := stmt.Exec(time.Now().Unix(), a.TokenId)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		fmt.Println(id)

		// update
		stmt2, err := db.Prepare("UPDATE bids SET win=TRUE WHERE carId=?")
		checkErr(err)

		res2, err := stmt2.Exec(a.WinnerCarId)
		checkErr(err)

		affect, err := res2.RowsAffected()
		checkErr(err)

		fmt.Println("Winner ? : ")
		fmt.Println(affect)

		// delete
		stmt, err = db.Prepare("DELETE FROM bids WHERE tokenId=? AND win=FALSE")
		checkErr(err)

		res, err = stmt.Exec(a.TokenId)
		checkErr(err)

		affect, err = res.RowsAffected()
		checkErr(err)

		fmt.Println("Nb Losers : ")
		fmt.Println(affect)
	}
	db.Close()

	w.WriteHeader(http.StatusOK)
}

func tokenHandle(w http.ResponseWriter, r *http.Request) {
	// Add a new token to db
	var t Token
	json.NewDecoder(r.Body).Decode(&t)
	fmt.Printf("%v\n", t)

	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO tokens(tokenId, tokenLat, tokenLon, tokenPrice, auctionEndTime) values(?, ?, ?, ?, 0)")
	checkErr(err)

	res, err := stmt.Exec(t.TokenId, t.TokenLat, t.TokenLon, t.TokenPrice)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	db.Close()

	w.WriteHeader(http.StatusOK)
}

func dataHandle(w http.ResponseWriter, r *http.Request) {
	// return all bids and tokens in a json
	db, err := sql.Open("sqlite3", "./data.sqlite")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM tokens")
	checkErr(err)

	var removeToken []Token
	var tokens []Token
	for rows.Next() {
		var id int
		var token Token
		err = rows.Scan(&id, &token.TokenId, &token.TokenLat, &token.TokenLon, &token.TokenPrice, &token.AuctionEndTime)
		checkErr(err)
		// if AuctionEndTime is more than 30 seconds remove the token
		if token.AuctionEndTime.Valid && token.AuctionEndTime.Int64 != 0 && token.AuctionEndTime.Int64+30 < time.Now().Unix() {
			// add token to remove list
			removeToken = append(removeToken, token)
		} else {
			tokens = append(tokens, token)
		}
	}
	rows.Close()

	// remove tokens and bid associated
	for _, token := range removeToken {
		// delete
		stmt, err := db.Prepare("DELETE FROM tokens WHERE tokenId=?")
		checkErr(err)

		res, err := stmt.Exec(token.TokenId)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println("remove token : ")
		fmt.Println(affect)

		// delete
		stmt, err = db.Prepare("DELETE FROM bids WHERE tokenId=?")
		checkErr(err)

		res, err = stmt.Exec(token.TokenId)
		checkErr(err)

		affect, err = res.RowsAffected()
		checkErr(err)

		fmt.Println("remove bids : ")
		fmt.Println(affect)
	}

	rows, err = db.Query("SELECT * FROM bids")
	checkErr(err)

	var bids []Bid
	for rows.Next() {
		var id int
		var bid Bid
		err = rows.Scan(&id, &bid.CarId, &bid.CarEnergy, &bid.CarRadius, &bid.CarLat, &bid.CarLon, &bid.Price, &bid.TokenId, &bid.Win)
		checkErr(err)
		bids = append(bids, bid)
	}
	rows.Close()

	db.Close()

	var data struct {
		Bids   []Bid
		Tokens []Token
	}

	data.Bids = bids
	data.Tokens = tokens

	json.NewEncoder(w).Encode(data)
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
	create table tokens (id integer not null primary key, tokenId text, tokenLat real, tokenLon real, tokenPrice integer, auctionEndTime integer);
	delete from tokens;
	`

	_, err = db.Exec(sqlStmt)
	checkErr(err)

	db.Close()

	// send ok http response
	w.WriteHeader(http.StatusOK)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
