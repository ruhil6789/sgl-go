package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gin-gonic/gin.v1"
)

// db driver visible to whole programmer
var DB *sql.DB

// train resources is there to hold the  rain info
type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

// station resource hold info about the station
type StationResource struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
}

// schedule resource hold schedule info
type schedule struct {
	Id          int
	TrainId     int
	StationId   int
	ArrivalTime time.Time
}

// register adds a paths and routes contains
func (t *TrainResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	// api will contain only content-type as application/JSON
	ws.Path("/v1/trains").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))
	container.Add(ws)

}

// get http://localhost:8000/v1/trains/1
func (t TrainResource) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	err := DB.QueryRow("select ID, DRIVER_NAME,OPERATING_STATUS FROM train where id =?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Train could not be found")
	} else {
		response.WriteEntity(t)
	}
}

// post http://localhost:8000/v1/trains
func (t TrainResource) createTrain(request *restful.Request, response *restful.Response) {

	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b TrainResource

	err := decoder.Decode(&b)
	log.Println(b.DriverName, b.OperatingStatus)
	// if err!= nil {
	//     log.Println(err)
	//     response.WriteErrorString(http.StatusBadRequest,"Invalid JSON")
	//     return
	// }
	statement, _ := DB.Prepare("insert into train(DRIVER_NAME,OPERATING_STATUS) values(?,?)")
	result, err := statement.Exec(b.DriverName, b.OperatingStatus)
	if err == nil {
		newId, _ := result.LastInsertId()
		b.ID = int(newId)
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}

}

// delete
func (t TrainResource) removeTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	statement, _ := DB.Prepare("delete from train where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}



func getStation(c *gin.Context){
	var station StationResource
  id:=c.Param("station_id")
  err:=DB.QueryRow("select ID,NAME,CAST(OPENING_TIME as CHAR),
           CAST(CLOSING_TIME as  CHAR) from station where id=?",id).Scan(&station.id,&station.Name,&station.OpeningTime,&station.ClosingTime)

   if err!=nil{
	log.Println(err)
	c.JSON(500,gin.H(
		"error:"err.Error(),
	))
   }else{
	c.JSON(200,gin.H(
		"result":station
	))
   }

}


//  create handle post
func createStation(c *gin.Context){
   
	var station StationResource
	
	if err:=c.BindJSON(&station);err==nil{
		statement,_:=DB.Prepare("insert into station(NAME,OPENING_TIME,CLOSING_TIME) values (? ? ?)")
		result,_:=statement.</span>Exec(station.Name,station.OpeningTime,station.closingTime)

		if err!==nil{
			newId,_:=result.LastInsertId()
			station.Id=int(newId)
			c.JSON(http.StatusOK,gin.h(
				"result":station
			))
		}else{
		c.String(http.StatusInternalServerError,err.Error())
	}
	}else{
		c.String(http.StatusInternalServerError,err.Error())
	}


}



func main() {
	var err error
	// connect to database
	db, err := sql.Open("sql", "./railApi.db")
	if err != nil {
		// log.Fatal(err)
		log.Println("Driver creation failed!")
	}

	// create tables
	dbUtils.initialize(db)
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	t := TrainResource{}
	t.register(wsContainer)
	log.Printf("start listening on localhost:8080")

	server := &http.Server{Addr: ":8000", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())



	//  add routes to the add verbs
	 r:=gin.Default()
	 r.GET("/v1/stations/:station_id",getStation)
	 r.POST("/v1/stations",createStation)
	 r.DELETE("/v1/stations/:station_id",removeStation)
	 r.Run(":8000")


}
