package mongo

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// News is for news collection
type News struct {
	ID         bson.ObjectId  `bson:"_id,omitempty"`
	NewsID     int            `bson:"news_id"`
	Name       string         `bson:"name"`
	URL        string         `bson:"url"`
	Categories []NewsCategory `bson:"categories"`
	CreatedAt  time.Time      `bson:"createdAt"`
}

// NewsCategory is categories element of news collection
type NewsCategory struct {
	Name string `bson:"name"`
	URL  string `bson:"url"`
}

// Articles is articles collection
type Articles struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	NewsID int           `bson:"news_id"`
	Title  string        `bson:"title"`
	Link   string        `bson:"link"`
	Date   time.Time     `bson:"lastBuildDate"`
	Items  []Item        `bson:"item"`
}

// Item is categories Items of articles collection
type Item struct {
	Title string    `bson:"title"`
	Link  string    `bson:"link"`
	Date  time.Time `bson:"pubDate"`
}

// Item2 is articles2 collection
type Item2 struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	NewsID int           `bson:"news_id"`
	Title  string        `bson:"title"`
	Link   string        `bson:"link"`
	Date   time.Time     `bson:"pubDate"`
}
