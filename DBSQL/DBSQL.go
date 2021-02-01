package DBSQL

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

type ProductsAmount struct {
	PRODUCTNAME string `json:"productname"`
	EA          string `json:"ea"`
	AMOUNT      string `json:"amount"`
	PRICE       string `json:"price"`
}

type CompanyInfo struct {
	SID     string `json:"sid"`
	CNAME   string `json:"cname"`
	TEL     string `json:"tel"`
	ADDRESS string `json:"address"`
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
	} else if flag == "WeekAmount" {
		values := make([]WeekAmount, length)
		temps := make([]WeekAmount, length)

		for days := 0; days < 7; days++ {
			now := time.Now()
			before := now.AddDate(0, 0, -6+days)
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
	} else if flag == "MonthAmount" {
		values := make([]WeekAmount, length)
		temps := make([]WeekAmount, length)

		for days := 0; days < 30; days++ {
			now := time.Now()
			before := now.AddDate(0, 0, -29+days)
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
	} else if flag == "YearAmount" {
		values := make([]WeekAmount, length)
		temps := make([]WeekAmount, length)

		for days := 0; days < 365; days++ {
			now := time.Now()
			before := now.AddDate(0, 0, -364+days)
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
	}

	return "없는 플레그 입니다."
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
	text := DBToString(rows, 7, "WeekAmount")

	return text
}

func GetWeekProductAmount(db *sql.DB, cNameCode string) string {
	fmt.Println("Get Week ProductAmount")

	return ""
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
	text := DBToString(rows, 30, "MonthAmount")

	return text
}

func GetMonthProductAmount(db *sql.DB) {
	fmt.Println("cccc")
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
	text := DBToString(rows, 365, "YearAmount")

	return text
}

func GetYearProductAmount(db *sql.DB) {
	fmt.Println("gggg")
}

func GetOneMonth(db *sql.DB) {
	fmt.Println("mmmm")
}

func GetThreeMonth(db *sql.DB) {
	fmt.Println("nnnn")
}

func GetStartEnd(db *sql.DB) {
	fmt.Println("bbbb")
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
