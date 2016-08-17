package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const NEWS_COLLECTION string = "news"
const ARTICLES_COLLECTION string = "articles"
const ARTICLES2_COLLECTION string = "articles2"

type News struct {
	ID         bson.ObjectId  `bson:"_id,omitempty"`
	NewsId     int            `bson:"news_id"`
	Name       string         `bson:"name"`
	Url        string         `bson:"url"`
	Categories []NewsCategory `bson:"categories"`
	CreatedAt  time.Time      `bson:"createdAt"`
}

type NewsCategory struct {
	Name string `bson:"name"`
	Url  string `bson:"url"`
}

type Articles struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	NewsId int           `bson:"news_id"`
	Title  string        `bson:"title"`
	Link   string        `bson:"link"`
	Date   time.Time     `bson:"lastBuildDate"`
	Items  []Item        `bson:"item"`
}

type Item struct {
	Title string    `bson:"title"`
	Link  string    `bson:"link"`
	Date  time.Time `bson:"pubDate"`
}

type Item2 struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	NewsId int           `bson:"news_id"`
	Title  string        `bson:"title"`
	Link   string        `bson:"link"`
	Date   time.Time     `bson:"pubDate"`
}

func (mg *Models) GetNewsData() ([]News, error) {
	mg.Db.GetCol(NEWS_COLLECTION)

	var news []News

	//get
	colQuerier := bson.M{}
	err := mg.Db.C.Find(colQuerier).All(&news)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (mg *Models) GetArticlesData(newsId int) ([]Item2, error) {
	mg.Db.GetCol(ARTICLES2_COLLECTION)

	var items []Item2

	//get
	colQuerier := bson.M{}
	if newsId != 0 {
		colQuerier = bson.M{"news_id": newsId}
	}

	err := mg.Db.C.Find(colQuerier).Sort("news_id").All(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
