package DBSQL

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Sales struct {
	SID         string `json:"sid"`
	PRODUCTNAME string `json:"productname"`
	EA          string `json:"ea"`
	SALESDATE   string `json:"salesdate"`
	CNAME       string `json:"cname"`
}

type WeekAmount struct {
	DAY    string `json:"day"`
	AMOUNT string `json:"amount`
}

type ProductsAmountBody struct {
	SALESDATE   string `json:"salesdate"`
	EA          string `json:"ea"`
	PRODUCTNAME string `json:"productname"`
	AMOUNT      string `json:"amount"`
}

type ProductDayAmount struct {
	NAME   string `json:"name"`
	EA     []int  `json:"ea"`
	AMOUNT []int  `json:"amount"`
}

type ProductsAmount struct {
	CATEGORIES []string           `json:"categories"`
	SERIES     []ProductDayAmount `json:"series"`
}

type CompanyInfo struct {
	SID     string `json:"sid"`
	CNAME   string `json:"cname"`
	TEL     string `json:"tel"`
	ADDRESS string `json:"address"`
}

type StartDayEndDay struct {
	START   string `json:"startday"`
	END     string `json:"endday"`
	COMPANY string `json:"company"`
}

type NoticeInfo struct {
	SID       string `json:"sid"`
	TITLE     string `json:"title"`
	CONTEXT   string `json:"context"`
	USERID    string `json:"user_id"`
	DATE      string `json:"date"`
	VIEWCOUNT string `json:"view_count"`
	SECTION   string `json:"section"`
}

func DBToString(rows *sql.Rows, length int, flag string) string {
	var i int = 0
	if flag == "COMPANY" {
		values := make([]CompanyInfo, length)
		for rows.Next() {
			rows.Scan(&values[i].SID, &values[i].CNAME, &values[i].TEL, &values[i].ADDRESS)
			i++
		}
		j, _ := json.Marshal(values)

		return string(j)
	} else if flag == "Amount" {
		values := make([]WeekAmount, length)
		temps := make([]WeekAmount, length)

		for days := 0; days < length; days++ {
			now := time.Now()
			before := now.AddDate(0, 0, -(length-1)+days)
			timebefore := fmt.Sprintf("%d-%02d-%02d", before.Year(), before.Month(), before.Day())
			values[days].DAY = timebefore
		}

		for rows.Next() {
			rows.Scan(&temps[i].DAY, &temps[i].AMOUNT)
			i++
		}

		for _, temp := range temps {
			if temp.DAY == "" {
				continue
			}
			fmt.Println(temp)

			for index, value := range values {
				if temp.DAY == value.DAY {
					values[index].AMOUNT = temp.AMOUNT
				}
			}
		}
		j, _ := json.Marshal(values)

		return string(j)
	} else if flag == "NOTICE" {
		values := make([]NoticeInfo, length)
		for rows.Next() {
			rows.Scan(&values[i].SID, &values[i].TITLE, &values[i].CONTEXT, &values[i].USERID, &values[i].DATE, &values[i].VIEWCOUNT, &values[i].SECTION)
			i++
		}
		j, _ := json.Marshal(values)

		return string(j)
	} else if flag == "NOTICE_ONE" {
		value := NoticeInfo{}
		for rows.Next() {
			rows.Scan(&value.SID, &value.TITLE, &value.CONTEXT, &value.USERID, &value.DATE, &value.VIEWCOUNT, &value.SECTION)
			i++
		}
		j, _ := json.Marshal(value)

		return string(j)
	}

	return "없는 플레그 입니다."
}

func DBToJson(body *sql.Rows, products *sql.Rows, length int, allLength int) string {
	values := ProductsAmount{}
	dayList := make([]string, length)
	names := make([]string, length)

	for days := 0; days < length; days++ {
		now := time.Now()
		before := now.AddDate(0, 0, -(length-1)+days)
		timebefore := fmt.Sprintf("%d-%02d-%02d", before.Year(), before.Month(), before.Day())
		dayList[days] = timebefore
	}
	values.CATEGORIES = dayList

	var i int = 0
	for products.Next() {
		products.Scan(&names[i])
		i++
	}
	temps := make([]ProductDayAmount, i)

	for index, _ := range temps {
		temps[index].NAME = names[index]
		temps[index].EA = make([]int, length)
		temps[index].AMOUNT = make([]int, length)
	}
	// fmt.Println(temps)

	i = 0
	bodies := make([]ProductsAmountBody, allLength)
	for body.Next() {
		body.Scan(&bodies[i].SALESDATE, &bodies[i].EA, &bodies[i].PRODUCTNAME, &bodies[i].AMOUNT)
		i++
	}

	for _, b := range bodies {
		for dayIndex, timeValue := range dayList {
			if timeValue == b.SALESDATE {
				for tt, ttt := range temps {
					if ttt.NAME == b.PRODUCTNAME {
						temps[tt].EA[dayIndex], _ = strconv.Atoi(b.EA)
						temps[tt].AMOUNT[dayIndex], _ = strconv.Atoi(b.AMOUNT)
					}
				}
			}
		}
	}
	values.SERIES = temps

	j, _ := json.Marshal(values)

	return string(j)
}

func GetWeekAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Week Amount")

	getSql := fmt.Sprintf(`select eom.sales_date, sum(eom.amount) from (select sales.sales_date, sum(sales.ea::int) as ea, sales.code,
		(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
		from sales, product
		where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '7 DAYS'
		and sales.code = product.code
		and product.c_sid = '%s'
		GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as eom group by eom.sales_date`, cNameCode)
	rows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	text := DBToString(rows, 7, "Amount")

	return text
}

func GetWeekProductAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Week ProductAmount")

	getSql := fmt.Sprintf(`select sales.sales_date, sum(sales.ea::int) as ea, sales.code, 
	(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
	from sales, product 
	where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '7 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.sales_date, sales.code ORDER BY sales_date`, cNameCode)
	productAmountRows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer productAmountRows.Close()

	getSql = fmt.Sprintf(`SELECT sales.code
	FROM sales, product
	WHERE to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '7 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.code`, cNameCode)
	productRows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer productRows.Close()

	getSql = fmt.Sprintf(`SELECT COUNT(a.*) FROM (select sales.sales_date, sum(sales.ea::int) as ea, sales.code, 
	(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
	from sales, product 
	where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '7 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as a`, cNameCode)
	var productAmountCnt int
	_ = db.QueryRow(getSql).Scan(&productAmountCnt)

	text := DBToJson(productAmountRows, productRows, 7, productAmountCnt)

	return text
}

func GetMonthAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Month Amount")

	getSql := fmt.Sprintf(`select eom.sales_date, sum(eom.amount) from (select sales.sales_date, sum(sales.ea::int) as ea, sales.code,
		(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
		from sales, product
		where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '30 DAYS'
		and sales.code = product.code
		and product.c_sid = '%s'
		GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as eom group by eom.sales_date`, cNameCode)
	rows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	text := DBToString(rows, 30, "Amount")

	return text
}

func GetMonthProductAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Month ProductAmount")

	getSql := fmt.Sprintf(`select sales.sales_date, sum(sales.ea::int) as ea, sales.code, 
	(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
	from sales, product 
	where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '30 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.sales_date, sales.code ORDER BY sales_date`, cNameCode)
	productAmountRows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer productAmountRows.Close()

	getSql = fmt.Sprintf(`SELECT sales.code
	FROM sales, product
	WHERE to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '30 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.code`, cNameCode)
	productRows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer productRows.Close()

	getSql = fmt.Sprintf(`SELECT COUNT(a.*) FROM (select sales.sales_date, sum(sales.ea::int) as ea, sales.code, 
	(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
	from sales, product 
	where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '30 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as a`, cNameCode)
	var productAmountCnt int
	_ = db.QueryRow(getSql).Scan(&productAmountCnt)

	text := DBToJson(productAmountRows, productRows, 30, productAmountCnt)

	return text
}

func GetYearAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Year Amount")

	getSql := fmt.Sprintf(`select eom.sales_date, sum(eom.amount) from (select sales.sales_date, sum(sales.ea::int) as ea, sales.code,
		(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
		from sales, product
		where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '365 DAYS'
		and sales.code = product.code
		and product.c_sid = '%s'
		GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as eom group by eom.sales_date`, cNameCode)
	rows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	text := DBToString(rows, 365, "Amount")

	return text
}

func GetYearProductAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Month ProductAmount")

	getSql := fmt.Sprintf(`select sales.sales_date, sum(sales.ea::int) as ea, sales.code, 
	(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
	from sales, product 
	where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '365 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.sales_date, sales.code ORDER BY sales_date`, cNameCode)
	productAmountRows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer productAmountRows.Close()

	getSql = fmt.Sprintf(`SELECT sales.code
	FROM sales, product
	WHERE to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '365 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.code`, cNameCode)
	productRows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer productRows.Close()

	getSql = fmt.Sprintf(`SELECT COUNT(a.*) FROM (select sales.sales_date, sum(sales.ea::int) as ea, sales.code, 
	(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
	from sales, product 
	where to_date(sales_date, 'YYYY-MM-DD') <= NOW() and to_date(sales_date, 'YYYY-MM-DD') > NOW() - interval '365 DAYS'
	and sales.code = product.code
	and product.c_sid = '%s'
	GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as a`, cNameCode)
	var productAmountCnt int
	_ = db.QueryRow(getSql).Scan(&productAmountCnt)

	text := DBToJson(productAmountRows, productRows, 365, productAmountCnt)

	return text
}

func GetStartEndAmount(db *sql.DB, dayjson StartDayEndDay) string {
	fmt.Println("Get Start End Amount")

	getSql := fmt.Sprintf(`select eom.sales_date, sum(eom.amount) from (select sales.sales_date, sum(sales.ea::int) as ea, sales.code,
		(select price from product where sales.code=product.code) * sum(sales.ea::int) as amount
		from sales, product
		where to_date(sales_date, 'YYYY-MM-DD') <= '%s' and to_date(sales_date, 'YYYY-MM-DD') > '%s'
		and sales.code = product.code
		and product.c_sid = '%s'
		GROUP BY sales.sales_date, sales.code ORDER BY sales_date) as eom group by eom.sales_date`, dayjson.END, dayjson.START, dayjson.COMPANY)
	rows, err := db.Query(getSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	getSql = fmt.Sprintf(`select to_date('%s', 'yyyy-mm-dd') - '%s'`, dayjson.END, dayjson.START)
	var periodAmountCnt int
	_ = db.QueryRow(getSql).Scan(&periodAmountCnt)

	text := DBToString(rows, periodAmountCnt+1, "Amount")

	return text
}

func GetCompanyInfo(db *sql.DB) string {
	fmt.Println("COMPANY info")
	var companyCnt int
	_ = db.QueryRow(`SELECT COUNT(*) FROM company`).Scan(&companyCnt)

	sqlString := fmt.Sprint(`SELECT * FROM company`)
	rows, err := db.Query(sqlString)
	if err != nil {
		return "Error in GetCompanyInfo GetInfo"
	}
	defer rows.Close()

	text := DBToString(rows, companyCnt, "COMPANY")
	return text
}

func ListLoad(db *sql.DB) string {
	getUserSql := fmt.Sprint("SELECT * FROM NOTICE ORDER BY SECTION DESC, DATE DESC")
	rows, err := db.Query(getUserSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cnt int
	_ = db.QueryRow(`select count(*) from NOTICE`).Scan(&cnt)
	text := DBToString(rows, cnt, "NOTICE")

	return text
}

func ListSize(db *sql.DB) string {
	var cnt int
	_ = db.QueryRow(`select count(*) from NOTICE`).Scan(&cnt)

	return strconv.Itoa(cnt)
}

func ListContext(db *sql.DB, sid string) string {
	sqlState := fmt.Sprintf("update notice set VIEW_COUNT=VIEW_COUNT+1 where SID=%s", sid)
	_, _ = db.Query(sqlState)

	getUserSql := fmt.Sprintf("SELECT * FROM NOTICE WHERE SID=%s", sid)
	rows, err := db.Query(getUserSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	text := DBToString(rows, 1, "NOTICE_ONE")

	return text
}

func ListCreate(db *sql.DB, data *NoticeInfo) string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	sqlState := fmt.Sprintf("INSERT INTO NOTICE (TITLE, CONTEXT, USER_ID, DATE, VIEW_COUNT, SECTION) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')", data.TITLE, data.CONTEXT, data.USERID, formatted, "0", strings.ToUpper(data.SECTION))
	rows, err := db.Query(sqlState)
	if err != nil {
		return "fail"
	}

	defer rows.Close()
	return "true"
}