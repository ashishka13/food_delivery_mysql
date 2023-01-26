package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"food_delivery_mysql/models"
	"io/ioutil"
	"log"
	"net"

	"github.com/go-sql-driver/mysql"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type mysqlSecret struct {
	MysqlHost     string `json:"mysql_host"`
	MysqlUser     string `json:"mysql_user"`
	MysqlPassword string `json:"mysql_password"`
	MysqlPort     string `json:"mysql_port"`
	MysqlDBName   string `json:"mysql_db_name"`
	NetworkType   string `json:"network_type"`
}

func Database() (db *sql.DB, err error) {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "my-secret-password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "food_delivery",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return db, nil
}

func getMysqlVars() (idPassword mysqlSecret, err error) {
	jsonfile, err := ioutil.ReadFile("utils/a.json")
	if err != nil {
		log.Println("file read error ", err)
		return
	}

	err = json.Unmarshal(jsonfile, &idPassword)
	if err != nil {
		log.Println("file unmarshal error ", err)
		return
	}
	return
}

func GetLocalIP() (localip string, err error) {
	localip, err = externalIP()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(localip)
	return localip, err
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	db, err := Database()
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return []Album{}, nil
	}

	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func GetAccount(db *sql.DB, id int, name string) (account models.CustomerAccount, err error) {
	query := `select * from account where id=?, name=?`
	rows, err := db.Query(query, id, name)
	if err != nil {
		log.Println("error occurred while getting account record", err)
		return
	}

	if !rows.Next() {
		log.Println("no records found")
		err = errors.New("no records found")
		return
	}

	err = rows.Scan(&account.ID, &account.Name, &account.Phone, &account.Address,
		&account.OrderID, &account.OrderPlaced, &account.PlacedTime,
		&account.ReceiveTime, &account.BoyName, &account.InProcess)
	if err != nil {
		log.Println("error occurred while decoding rows", err)
		return
	}
	return
}

func GetAccounts(db *sql.DB, query string) (accounts []models.CustomerAccount, err error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Println("error occurred while getting account records", err)
		return
	}

	for rows.Next() {
		singleAccount := models.CustomerAccount{}
		err = rows.Scan(&singleAccount.ID, &singleAccount.Name, &singleAccount.Phone, &singleAccount.Address,
			&singleAccount.OrderID, &singleAccount.OrderPlaced, &singleAccount.PlacedTime,
			&singleAccount.ReceiveTime, &singleAccount.BoyName, &singleAccount.InProcess)
		if err != nil {
			log.Println("error occurred while decoding rows", err)
			return
		}
		accounts = append(accounts, singleAccount)
	}
	return
}

func MyPrint(strs ...string) {
	log.Println("\n======================================================\n", strs,
		"\n======================================================")
}
