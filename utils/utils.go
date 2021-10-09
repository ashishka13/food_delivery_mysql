package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql" // this is mysql driver import
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type mysqlSecret struct {
	Mysql_host     string `json:"mysql_host"`
	Mysql_user     string `json:"mysql_user"`
	Mysql_password string `json:"mysql_password"`
	Mysql_port     string `json:"mysql_port"`
	Mysql_db       string `json:"mysql_db"`
}

// InitializeDatabase will establish database connection
func InitializeDatabase() (err error) {
	idpassword, err := getMysqlVars()
	if err != nil {
		return
	}
	dbUser := idpassword.Mysql_user
	dbPass := idpassword.Mysql_password
	dbName := idpassword.Mysql_db
	dbHost := idpassword.Mysql_host
	dbPort := idpassword.Mysql_port

	dburl1 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, "")

	db, err := sql.Open("mysql", dburl1)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbName)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}

	dburl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dburl)
	log.Println("database connectin error ping db ", err)
	if err != nil {
		log.Println("database connection error")
		return
	}

	_, err = db.Exec("CREATE TABLE students (id int NOT NULL AUTO_INCREMENT, fname varchar(255) NOT NULL, lname varchar(255), PRIMARY KEY (id))")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	defer db.Close()

	return
}

func Database() (db *sql.DB, err error) {
	idpassword, err := getMysqlVars()
	if err != nil {
		return
	}
	dbUser := idpassword.Mysql_user
	dbPass := idpassword.Mysql_password
	dbName := idpassword.Mysql_db
	dbHost := idpassword.Mysql_host
	dbPort := idpassword.Mysql_port

	dburl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dburl)
	log.Println("database connectin error ping db ", err)
	if err != nil {
		log.Println("database connection error")
		return
	}
	return
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

func DatabaseConnect(dbname string) (db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("database connect error", err)
		return
	}
	db = client.Database(dbname)
	return
}
