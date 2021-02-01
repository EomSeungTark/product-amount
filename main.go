package main

import (
	"database/sql"
	"fmt"
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
	DBSQL.GetMonthProductAmount(db)
	return c.String(http.StatusOK, "return ok")
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
	DBSQL.GetYearProductAmount(db)
	return c.String(http.StatusOK, "return ok")
}

func getOneMonth(c echo.Context) error {
	defer c.Request().Body.Close()
	DBSQL.GetOneMonth(db)
	return c.String(http.StatusOK, "return ok")
}

func getThreeMonth(c echo.Context) error {
	defer c.Request().Body.Close()
	DBSQL.GetThreeMonth(db)
	return c.String(http.StatusOK, "return ok")
}

func getStartEnd(c echo.Context) error {
	defer c.Request().Body.Close()
	DBSQL.GetStartEnd(db)
	return c.String(http.StatusOK, "return ok")
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

	groupPeriod.GET("/GET/ONEMONTH", getOneMonth)
	groupPeriod.GET("/GET/THREEMONTH", getThreeMonth)
	groupPeriod.GET("/GET/START-END", getStartEnd)

	e.Start(":8000")
}
