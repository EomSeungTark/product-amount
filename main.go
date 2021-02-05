package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/eom/product-amount/DBSQL"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "800326"
	DB_NAME     = "postgres"
)

func getWeekAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	cNameCode := c.Param("cname-code")
	decodedValue, err := url.QueryUnescape(cNameCode)
	if err != nil {
		log.Fatal(err)
	}

	text := DBSQL.GetWeekAmount(db, decodedValue)
	return c.String(http.StatusOK, text)
}

func getWeekProductAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	cNameCode := c.Param("cname-code")
	decodedValue, err := url.QueryUnescape(cNameCode)
	if err != nil {
		log.Fatal(err)
	}

	text := DBSQL.GetWeekProductAmount(db, decodedValue)
	return c.String(http.StatusOK, text)
}

func getMonthAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	cNameCode := c.Param("cname-code")
	decodedValue, err := url.QueryUnescape(cNameCode)
	if err != nil {
		log.Fatal(err)
	}

	text := DBSQL.GetMonthAmount(db, decodedValue)
	return c.String(http.StatusOK, text)
}

func getMonthProductAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	cNameCode := c.Param("cname-code")
	decodedValue, err := url.QueryUnescape(cNameCode)
	if err != nil {
		log.Fatal(err)
	}

	text := DBSQL.GetMonthProductAmount(db, decodedValue)
	return c.String(http.StatusOK, text)
}

func getYearAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	cNameCode := c.Param("cname-code")
	decodedValue, err := url.QueryUnescape(cNameCode)
	if err != nil {
		log.Fatal(err)
	}

	text := DBSQL.GetYearAmount(db, decodedValue)
	return c.String(http.StatusOK, text)
}

func getYearProductAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	cNameCode := c.Param("cname-code")
	decodedValue, err := url.QueryUnescape(cNameCode)
	if err != nil {
		log.Fatal(err)
	}

	text := DBSQL.GetYearProductAmount(db, decodedValue)
	return c.String(http.StatusOK, text)
}

func getStartEndAmount(c echo.Context) error {
	defer c.Request().Body.Close()

	StartDayEndDay := DBSQL.StartDayEndDay{}

	byte, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(byte, &StartDayEndDay)
	text := DBSQL.GetStartEndAmount(db, StartDayEndDay)

	return c.String(http.StatusOK, text)
}

func getCompanyInfo(c echo.Context) error {
	defer c.Request().Body.Close()
	companyInfos := DBSQL.GetCompanyInfo(db)
	return c.String(http.StatusOK, companyInfos)
}

func main() {
	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))

	e.GET("/GET/COMPANYID", getCompanyInfo)

	groupCompany := e.Group("/SALES/COMPANY")
	groupPeriod := e.Group("/SALES/PERIOD")

	groupCompany.GET("/GET/WEEK/AMOUNT/:cname-code", getWeekAmount)
	groupCompany.GET("/GET/WEEK/PRODUCTS-AMOUNT/:cname-code", getWeekProductAmount)

	groupCompany.GET("/GET/MONTH/AMOUNT/:cname-code", getMonthAmount)
	groupCompany.GET("/GET/MONTH/PRODUCTS-AMOUNT/:cname-code", getMonthProductAmount)

	groupCompany.GET("/GET/YEAR/AMOUNT/:cname-code", getYearAmount)
	groupCompany.GET("/GET/YEAR/PRODUCTS-AMOUNT/:cname-code", getYearProductAmount)

	groupPeriod.POST("/GET/START-END", getStartEndAmount)

	e.Start(":8000")
}
