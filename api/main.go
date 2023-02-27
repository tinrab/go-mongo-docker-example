package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

type Post struct {
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

var coll *mongo.Collection

func main() {
	log.Println("enter main - connecting to mongo")
	// Connect to mongo
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo"))

	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Get posts collection
	coll = client.Database("app").Collection("posts")

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/posts", createPost).
		Methods("POST")
	r.HandleFunc("/posts", readPosts).
		Methods("GET")

	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
}

func createPost(w http.ResponseWriter, r *http.Request) {
	// Read body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read post
	post := &Post{}
	err = json.Unmarshal(data, post)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.CreatedAt = time.Now().UTC()

	// Insert new post
	if _, err := coll.InsertOne(context.TODO(), post); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, post)
}

func readPosts(w http.ResponseWriter, r *http.Request) {
	result := []Post{}
	if cursor, err := coll.Find(context.TODO(), bson.D{}); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		for cursor.Next(context.TODO()) {
			var elem Post
			err := cursor.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			result = append(result, elem)
		}
		responseJSON(w, result)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
