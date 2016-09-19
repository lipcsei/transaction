package main

import (
	"strconv"

	"github.com/foosio/api/lib/services/db"
	"github.com/foosio/api/lib/services/env"
	"github.com/wawandco/fako"
	"gopkg.in/mgo.v2/bson"
	redis "gopkg.in/redis.v4"
)

type Redis struct{}

func (r *Redis) connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     env.Get("REDIS_HOST", "localhost:6379"),
		Password: "",
		DB:       0,
	})
	return client
}

var mongo *db.Mongo

const LIMIT = 100

type User struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name" fako:"full_name"`
}

type Game struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Size  int    `json:"size" bson:"size"`
	Users []User `json:"users" bson:"users"`
}

func Join(gameID string, userID int) {
	collection, session := mongo.GetCollection("transaction")
	defer session.Close()

	user := User{ID: strconv.Itoa(userID + 1)}
	fako.Fill(&user)

	collection.Update(bson.M{"_id": gameID}, bson.M{"$push": bson.M{"users": &user}})
}

func main() {

	gameID := "1"
	collection, session := mongo.GetCollection("transaction")
	defer session.Close()
	collection.RemoveId(gameID)
	game := Game{ID: gameID, Size: LIMIT}
	collection.Insert(&game)

	for userID := 0; userID < LIMIT; userID++ {
		Join(gameID, userID)
	}

}
