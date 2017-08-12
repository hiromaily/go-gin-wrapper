package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const newsCollection string = "news"
const articlesCollection string = "articles"
const articles2Collection string = "articles2"

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

// GetNewsData is to get all news data from news collection
func (mg *Models) GetNewsData() ([]News, error) {
	//mg.Db.Session
	mg.Db.GetCol(newsCollection)

	var news []News

	//get
	colQuerier := bson.M{}
	err := mg.Db.C.Find(colQuerier).All(&news)
	if err != nil {
		return nil, err
	}

	return news, nil
}

// GetArticlesData is to get one or all articles from article collection
func (mg *Models) GetArticlesData(newsID int) ([]Articles, error) {
	mg.Db.GetCol(articlesCollection)

	var articles []Articles

	//get
	colQuerier := bson.M{}
	if newsID != 0 {
		colQuerier = bson.M{"news_id": newsID}
	}

	err := mg.Db.C.Find(colQuerier).Sort("news_id").All(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// GetArticlesData2 is to get one or all articles from article2 collection
func (mg *Models) GetArticlesData2(newsID int) ([]Item2, error) {
	mg.Db.GetCol(articles2Collection)

	var items []Item2

	//get
	colQuerier := bson.M{}
	if newsID != 0 {
		colQuerier = bson.M{"news_id": newsID}
	}

	err := mg.Db.C.Find(colQuerier).Sort("news_id").All(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
