package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type jsonData struct {
	//
	Meeting `json:"Meeting" bson:"Meeting"`
}

//Meeting blablabla
type Meeting struct {
	//Meeting : sldfjwoeijaf
	ID               int           `json:"id" bson:"id"`
	Title            string        `json:"title" bson:"title"`
	Participantsdata []Participant `json:"participant" bson:"participant"`
	StartTime        time.Time     `json:"start_time" bson:"start_time" `
	EndTime          time.Time     `json:"end_time" bson:"end_time"`
	TimeNow          time.Time     `json:"time_now" bson:"time_now"`
}

//Participant blablabla
type Participant struct {
	//Participant : sldjnfoweiaf
	Name  string ` bson:"name"`
	Email string ` bson:"email"`
	Rsvp  string ` bson:"rsvp"`
}

//timeframe function

func inTimeSpan(in, out, check time.Time) bool {
	return check.After(in) && check.Before(out)

}

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:8080")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("test").Collection("meetings")

	// Some dummy data to add to the Database
	// defining a struct instance

	var meetings []Meeting

	// JSON array to be decoded
	// to an array in golang
	//just a representation
	//dummy data
	Data := []byte(` 
    [ 
		 {
			 "id" = 1234
	    "Title":"kj",
	    "Paticipants":[{
	        "name":"ayush",
	        "email":"xyz@gmail.com",
	        "rsvp":"no"
	    },
	    {
	        "name":"piyush",
	        "email":"xyz@gmail.com",
	        "rsvp":"yes"
	    },
	    {
	        "name":"sachin",
	        "email":"xyz@gmail.com",
	        "rsvp":"maybe"
	    }],
	    "Start_time":"2013-10-21T13:28:06.419Z",
	    "End_time":"2013-10-21T13:28:06.419Z",
	    "Time_now":"2013-10-21T13:28:06.419Z"
		}
		]`)

	// decoding JSON array to
	// the country array
	errr := json.Unmarshal(Data, &meetings)

	if errr != nil {

		// if error is not nil
		// print error
		fmt.Println(errr)
	}

	///////https requests

	// get meetings through ID
	getMeetingthroughID := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			filter := bson.M{"id": "234556"}

			// Find a single document
			findOptions := options.Find()

			var result []*Meeting

			// Finding multiple documents returns a cursor
			cur, err := collection.Find(context.TODO(), filter, findOptions)
			if err != nil {
				log.Fatal(err)
			}

			// Iterate through the cursor
			for cur.Next(context.TODO()) {
				var element Meeting
				err := cur.Decode(&element)
				if err != nil {
					log.Fatal(err)
				}

				result = append(result, &element)
			}

			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}

			// Close the cursor once finished
			cur.Close(context.TODO())
			json.NewEncoder(w).Encode(result)
		}
	}

	http.HandleFunc("/meeting/<id here>", getMeetingthroughID)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//POST meeting
	scheduleMeeting := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			//dummy data
			Data := []byte(` 
			[ 
				 {
					 "id" = 1234
				"Title":"kj",
				"Paticipants":[{
					"name":"ayush",
					"email":"xyz@gmail.com",
					"rsvp":"no"
				},
				{
					"name":"piyush",
					"email":"xyz@gmail.com",
					"rsvp":"yes"
				},
				{
					"name":"sachin",
					"email":"xyz@gmail.com",
					"rsvp":"maybe"
				}],
				"Start_time":"2013-10-21T13:28:06.419Z",
				"End_time":"2013-10-21T13:28:06.419Z",
				"Time_now":"2013-10-21T13:28:06.419Z"
				}
				]`)

			// decoding JSON array to
			// the country array
			errr := json.Unmarshal(Data, &meetings)

			if errr != nil {

				// if error is not nil
				// print error
				fmt.Println(errr)
			}

			// Insert a single document
			insertResult, err := collection.InsertOne(context.TODO(), meetings[0])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted a single document: ", insertResult.InsertedID)

		}
	}

	http.HandleFunc("/meetings", scheduleMeeting)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//List all meetings of a participant
	participantAllMeetings := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/articles?participant=<email id>" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		if req.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}
		if req.Method == "GET" {
			filter := bson.M{"email": "xyz@gmail.com"}

			// Find a single document

			findOptions := options.Find()

			var participantAll []*Meeting

			// Finding multiple documents returns a cursor
			cur, err := collection.Find(context.TODO(), filter, findOptions)
			if err != nil {
				log.Fatal(err)
			}

			// Iterate through the cursor
			for cur.Next(context.TODO()) {
				var elem Meeting
				err := cur.Decode(&elem)
				if err != nil {
					log.Fatal(err)
				}

				participantAll = append(participantAll, &elem)
			}

			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}

			// Close the cursor once finished
			cur.Close(context.TODO())

			json.NewEncoder(w).Encode(participantAll)
		}
	}

	http.HandleFunc("/articles?participant=<email id>", participantAllMeetings)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//List all meetings with in a timeframe

	meetingsThroughTimeframe := func(w http.ResponseWriter, req *http.Request) {

		//security
		if req.URL.Path != "/meetings?start=<start time here>&end=<end time here>" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		if req.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		if req.Method == "GET" {

			// Find a single document

			findOptions := options.Find()

			var timeFrameResult []*Meeting

			// Finding multiple documents returns a cursor
			cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
			if err != nil {
				log.Fatal(err)
			}

			// Iterate through the cursor
			for cur.Next(context.TODO()) {

				in, _ := time.Parse(time.RFC822, "01 Jan 15 20:00 UTC")
				out, _ := time.Parse(time.RFC822, "01 Jan 17 10:00 UTC")
				var timeelem Meeting
				err := cur.Decode(&timeelem)
				if err != nil {
					log.Fatal(err)
				}

				if inTimeSpan(in, out, timeelem.StartTime) {
					timeFrameResult = append(timeFrameResult, &timeelem)
				}

			}

			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}

			// Close the cursor once finished
			cur.Close(context.TODO())
			json.NewEncoder(w).Encode(timeFrameResult)
		}
	}

	http.HandleFunc("/meetings?start=<start time here>&end=<end time here>", meetingsThroughTimeframe)
	log.Fatal(http.ListenAndServe(":8080", nil))

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}
}
