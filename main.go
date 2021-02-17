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

func listServe(c echo.Context) error {
	defer c.Request().Body.Close()

	noticeList := DBSQL.ListLoad(db)
	return c.String(http.StatusOK, noticeList)
}

func listContext(c echo.Context) error {
	defer c.Request().Body.Close()

	sid := c.Param("sid")
	noticeContext := DBSQL.ListContext(db, sid)

	var li DBSQL.NoticeInfo
	json.Unmarshal([]byte(noticeContext), &li)
	fmt.Println(li.SID)

	if li.SID != "" {
		return c.String(http.StatusOK, noticeContext)
	} else {
		return c.String(http.StatusBadRequest, noticeContext)
	}
}

func listCreate(c echo.Context) error {
	u := new(DBSQL.NoticeInfo)
	defer c.Request().Body.Close()

	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "return not ok")
	}
	result := DBSQL.ListCreate(db, u)

	if result == "true" {
		return c.String(http.StatusOK, "return ok")
	} else {
		return c.String(http.StatusBadRequest, "create fail")
	}
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

	e.GET("/LIST/GETLISTS", listServe)
	e.GET("/LIST/GETLIST/:sid", listContext)
	e.POST("/LIST/CREATELIST", listCreate)

	// e.POST("/login", LOGINFEATURE.Login)
	// e.GET("/", LOGINFEATURE.Accessible)
	// r := e.Group("/restricted")

	// // Configure middleware with the custom claims type
	// config := middleware.JWTConfig{
	// 	Claims:     &(LOGINFEATURE.JwtCustomClaims{}),
	// 	SigningKey: []byte("secret"),
	// }
	// r.Use(middleware.JWTWithConfig(config))
	// r.GET("", LOGINFEATURE.Restricted)

	e.Start(":8000")
}
