package controllers

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"../datamodel"
	_ "github.com/go-sql-driver/mysql" //mysqlのドライバーをここでは明示的に呼び出していないのでblank importの形になっている
	"github.com/gocraft/dbr"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

var dbinstance *dbr.Connection

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func getDBInstance() (*dbr.Connection, error) {
	var err error
	if dbinstance == nil {
		conn, err := sql.Open(datamodel.DBTYPE, datamodel.PATH+datamodel.DATABASE_NAME)
		if err != nil {
			log.Errorf("mysql: could not get a connection: %v", err)
			os.Exit(1)
		}
		defer conn.Close()

		// Check the connection.
		if conn.Ping() == driver.ErrBadConn {
			log.Errorf("mysql: could not connect to the database. " +
				"could be bad address, or this address is not whitelisted for access.")
			os.Exit(1)

		}
		dbinstance, err = dbr.Open(datamodel.DBTYPE, datamodel.PATH+datamodel.DATABASE_NAME, nil)
		if err != nil {
			fmt.Println("Success!DBConnection!")
		}
	}
	return dbinstance, err
}

/*
/packet/~
*/

//RootConnect はrootにアクセスした時に呼ばれる。id=1だけ呼び出すお試しのもの
func RootConnect(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	var dbpacket datamodel.DBPacket
	sess.Select("*").From(datamodel.PACKET_TABLENAME).Where("id = ?", 1).Load(&dbpacket)

	return c.JSON(http.StatusCreated, dbpacket)
}

//PacketDataSelectID は id=?に応じたのを返す
func PacketDataSelectID(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	id, _ := strconv.Atoi(c.QueryParam("id"))

	var packet datamodel.DBPacket
	sess.Select("*").From(datamodel.PACKET_TABLENAME).Where("id = ?", id).Load(&packet)
	return c.JSON(http.StatusCreated, packet)
}

//PacketDataSelectNew は最新のidの情報を返す
func PacketDataSelectNew(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	var packet datamodel.DBPacket

	sess.Select("*").From(datamodel.PACKET_TABLENAME).Where("id = (SELECT MAX(id) FROM " + datamodel.PACKET_TABLENAME + ")").Load(&packet)

	return c.JSON(http.StatusCreated, packet)
}

//DistanceSelectMacAddress src_mac=?,dst_mac=?に応じたのを返す
func PacketDataSelectMacAddress(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	srcMac := c.QueryParam("src_macaddress")
	dstMac := c.QueryParam("dst_macaddress")

	var packet []datamodel.DBPacket
	ch := make(chan bool)
	go func(ch chan bool) {
		sess.Select("*").From(datamodel.PACKET_TABLENAME).Where("src_mac=? OR dst_mac=?", srcMac, dstMac).Load(&packet)
		ch <- true
	}(ch)
	<-ch
	fmt.Println(packet)

	return c.JSON(http.StatusCreated, packet)
}

/*
/distance/~
*/

//DistanceSelectID はid=?に応じたのを返す
func DistanceSelectID(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	id, _ := strconv.Atoi(c.QueryParam("id"))

	var distance datamodel.DBDistance
	sess.Select("*").From(datamodel.DISTANCE_TABLENAME).Where("id = ?", id).Load(&distance)
	return c.JSON(http.StatusCreated, distance)
}

//DistanceSelectNew は最新のidの情報を返す
func DistanceSelectNew(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)
	var distance datamodel.DBDistance
	sess.Select("*").From(datamodel.PACKET_TABLENAME).Where("id = (SELECT MAX(id) FROM " + datamodel.DISTANCE_TABLENAME + ")").Load(&distance)
	return c.JSON(http.StatusCreated, distance)
}

//DistanceSelectMacAddress はmacaddress=?に応じたのを返す
func DistanceSelectMacAddress(c echo.Context) error {
	conn, err := getDBInstance()
	if err != nil {
		os.Exit(-1)
	}
	sess := conn.NewSession(nil)

	macaddress := c.QueryParam("macaddress")
	rpimacaddress := c.QueryParam("rpi_macaddress")
	order_one, _ := strconv.Atoi(c.QueryParam("new_order_one"))

	var distance []datamodel.DBDistance

	if order_one == 1 {
		//ex test
		// localhost:3000/api/distance/macaddress?macaddress=84:89:AD:8D:85:F6&rpi_macaddress=b827ebf277a4&order_one=1
		// select * from distance where id = (select MAX(id) from distance where macaddr = "84:89:AD:8D:85:F6" AND rpimac = "b827ebf277a4")
		sess.Select("*").From(datamodel.DISTANCE_TABLENAME).
			Where("id = (select MAX(id) from "+datamodel.DISTANCE_TABLENAME+
				" where macaddr = ? AND rpimac = ?)", macaddress, rpimacaddress).Load(&distance)
	} else {
		sess.Select("*").From(datamodel.DISTANCE_TABLENAME).Where("macaddr = ?", macaddress).Load(&distance)
	}
	return c.JSON(http.StatusCreated, distance)
}
