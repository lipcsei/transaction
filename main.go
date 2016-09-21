package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/foosio/api/lib/services/db"
	"github.com/wawandco/fako"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongo *db.Mongo
)

const LIMIT = 1000000

type (
	// User struct represent a User
	User struct {
		ID   string `json:"id" bson:"_id,omitempty"`
		Name string `json:"name" bson:"name" fako:"full_name"`
		PID  int    `json:"pid" bson:"pid"`
	}

	// Game struct represent a Game
	Game struct {
		ID          string    `json:"id" bson:"_id,omitempty"`
		Size        int       `json:"size" bson:"size"`
		Users       []User    `json:"users" bson:"users"`
		LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
	}
)

func Join(gameID string, userID int) error {
	collection, session := mongo.GetCollection("transaction")
	defer session.Close()

	user := User{ID: strconv.Itoa(userID + 1)}
	fako.Fill(&user)
	user.PID = os.Getppid()

	var game Game
	err := collection.Find(bson.M{"_id": gameID}).One(&game)

	if err != nil {
		fmt.Println(gameID)
		return err
	}

	err = collection.Update(bson.M{
		"_id":         gameID,
		"lastUpdated": game.LastUpdated,
		"size":        bson.M{"$lt": LIMIT},
	},
		bson.M{
			"$push":        bson.M{"users": &user},
			"$currentDate": bson.M{"lastUpdated": true},
			"$inc":         bson.M{"size": 1},
		})

	return err
}

func main() {
	gameID := "1"
	CreateGame(gameID)

	success := 0
	failed := 0
	for userID := 0; userID < LIMIT; userID++ {
		err := Join(gameID, userID)
		if err != nil {
			failed++
			fmt.Println(err)
			// break
		} else {
			success++
		}
	}

	fmt.Println("success: ", success, "failed: ", failed, os.Getpid())
}

func CreateGame(gameID string) {
	collection, session := mongo.GetCollection("transaction")
	defer session.Close()
	collection.RemoveId(gameID)
	game := Game{ID: gameID, Size: 0, LastUpdated: time.Now()}
	collection.Insert(&game)

}
