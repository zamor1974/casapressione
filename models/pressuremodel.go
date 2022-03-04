package models

import (
	"casapressione/lang"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// swagger:model Pressure
type Pressure struct {
	// Id of Pressure value
	// in: int64
	Id int64 `json:"id"`
	// Value of Pressure
	// in: float
	Value float32 `json:"valore"`
	// Timestamp of insert
	// in: time
	DateInsert time.Time `json:"data_inserimento"`
}

type Pressures []Pressure

type ReqAddPressure struct {
	// Value of the Pressure
	// in: float
	Value float32 `json:"valore" validate:"required"`
}

// swagger:parameters add Pressure
type ReqPressureBody struct {
	// - name: body
	//  in: body
	//  description: Pressure
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddPressure"
	//  required: true
	Body ReqAddPressure `json:"body"`
}

// ErrHandler returns error message bassed on env debug
func ErrHandler(err error) string {
	var errmessage string
	if os.Getenv("DEBUG") == "true" {
		errmessage = err.Error()
	} else {
		errmessage = lang.Get("something_went_wrong")
	}
	return errmessage
}

func GetPressuresSqlx(db *sql.DB) *Pressures {
	pressures := Pressures{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM pressione order by id desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Pressure
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		pressures = append(pressures, p)
	}
	return &pressures
}
func GetLastAltitudeSqlx(db *sql.DB) *Pressures {
	pressures := Pressures{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM pressione where id = (select max(id) from pioggia)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Pressure
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		pressures = append(pressures, p)
	}
	return &pressures
}
func GetLastHourSqlx(db *sql.DB) *Pressures {
	pressures := Pressures{}

	tFine := time.Now()
	dataFine := tFine.Format("2006-01-02 15:04:05")

	tInizio := time.Now().Add(time.Duration(-1) * time.Hour)
	dataInizio := tInizio.Format("2006-01-02 15:04:05")

	sqlStatement := fmt.Sprintf("SELECT id,valore,data_inserimento FROM pressione where data_inserimento  >= '%s' AND data_inserimento <= '%s'", dataInizio, dataFine)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Pressure
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		pressures = append(pressures, p)
	}

	if len(pressures) == 0 {
		elemento := GetLastAltitudeSqlx(db)
		pressures = append(pressures, *elemento...)
	}
	return &pressures
}

// PostPressureSqlx insert Pressure value
func PostPressureSqlx(db *sql.DB, reqPressure *ReqAddPressure) (*Pressure, string) {

	value := reqPressure.Value

	var pressure Pressure

	lastInsertId := 0

	//sqlStatement := fmt.Sprintf("insert into 'pioggia' ('valore','data_inserimento') values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	sqlStatement := fmt.Sprintf("insert into pressione (valore,data_inserimento) values (%.2f,CURRENT_TIMESTAMP) RETURNING id", value)
	log.Println(sqlStatement)

	err := db.QueryRow(sqlStatement).Scan(&lastInsertId)

	if err != nil {
		return &pressure, ErrHandler(err)
	}

	sqlStatement1 := fmt.Sprintf("SELECT id,valore,data_inserimento FROM pressione where id = %d", lastInsertId)
	rows, err := db.Query(sqlStatement1)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Pressure
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		pressure = p
	}
	if err != nil {
		return &pressure, lang.Get("no_result")
	}
	return &pressure, ""
}
