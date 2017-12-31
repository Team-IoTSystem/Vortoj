package lifecycle

import (
	"database/sql"
	"database/sql/driver"

	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"

	"github.com/gocraft/dbr"
)

const (
	PATH          string = "root:root@tcp(127.0.0.1:3306)/"
	DATABASE_NAME string = "vortojpacket"
	TABEL_NAME    string = "packet"
	DBTYPE        string = "mysql"
)

var dbinstance *dbr.Connection

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS ` + DATABASE_NAME + ` CHARACTER SET 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE ` + DATABASE_NAME + `;`,
	`CREATE TABLE IF NOT EXISTS ` + TABEL_NAME + ` (
		id INT UNSIGNED AUTO_INCREMENT NOT NULL,
		deviceid VARCHAR(255), 
		src_mac VARCHAR(255),
		dst_mac VARCHAR(255),
		src_ip VARCHAR(255),
		dst_ip VARCHAR(255),
		src_port VARCHAR(255), 
		dst_port VARCHAR(255), 
		syn INT, 
		ack INT, 
		sequence BIGINT UNSIGNED,
		protocol VARCHAR(255), 
		length INT, 
		datachank BLOB,
		PRIMARY KEY (id)
	);`,
}

type TPacket struct {
	ID        int16  `db:"id"`
	DeviceID  string `db:"deviceid"`
	SrcMAC    string `db:"src_mac"`
	DstMAC    string `db:"dst_mac"`
	SrcIP     string `db:"src_ip"`
	DstIP     string `db:"dst_ip"`
	SrcPort   string `db:"src_port"`
	DstPort   string `db:"dst_port"`
	SYN       bool   `db:"syn"`
	ACK       bool   `db:"ack"`
	Sequence  int64  `db:"sequence"`
	Protocol  string `db:"protocol"`
	Length    int64  `db:"length"`
	DataChank []byte `db:"datachank"`
}

func ensureTablesExist() error {
	conn, err := sql.Open(DBTYPE, PATH)
	if err != nil {
		return fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	defer conn.Close()

	// Check the connection.
	if conn.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql: could not connect to the database. " +
			"could be bad address, or this address is not whitelisted for access.")
	}

	if _, err := conn.Exec("USE " + DATABASE_NAME); err != nil {
		// MySQL error 1049 is "database does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createTable(conn)
		}
	}

	return nil
}

func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		log.Println(stmt)
		_, err := conn.Exec(stmt)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func getDBInstance() (*dbr.Connection, error) {
	var err error = nil
	if dbinstance == nil {
		log.Println("create " + PATH + " DB")
		if err := ensureTablesExist(); err != nil {
			return nil, err
		}

		conn, err := sql.Open(DBTYPE, PATH+DATABASE_NAME)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
		}
		if err := conn.Ping(); err != nil {
			conn.Close()
			return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
		}

		//createAndOpen()
		dbinstance, err = dbr.Open(DBTYPE, PATH+DATABASE_NAME, nil)
	}
	return dbinstance, err
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func InsertPacketData(p *TPacket) bool {
	insertpacket := p

	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	stmt := sess.InsertInto("packet").Columns(
		"deviceid",
		"src_mac",
		"dst_mac",
		"src_ip",
		"dst_ip",
		"src_port",
		"dst_port",
		"syn",
		"ack",
		"sequence",
		"protocol",
		"length",
		"datachank").Record(insertpacket)

	result, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		//result.RowsAffected()
		count, _ := result.RowsAffected()
		fmt.Println(count)
		return true
	}
}
