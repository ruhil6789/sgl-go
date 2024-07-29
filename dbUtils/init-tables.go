package dbutils

import (
    "database/sql"
    "log"
 "database/sql"
)

func initialize(dbDriver *sql.DB){
	statement,driverError:=dbDriver.Prepare(train)
	if driverError!=nil{
		log.Fatal(driverError)
	}

	// create train table

	_, statementError:= statement.Exec()
	if statementError!=nil{
        log.Println("Table already created")
    }

	statement,_=dbDriver.Prepare(station)
	statement.Exec()
	statement,_=dbDriver.Prepare(schedule)
	statement.Exec()
	log.Println("All tables created/initialized successfully")
}