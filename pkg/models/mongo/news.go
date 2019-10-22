package mongo

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/hiromaily/golibs/db/mongodb"
)

const newsCollection string = "news"
const articlesCollection string = "articles"
const articles2Collection string = "articles2"

// MongoModel is extension of mongo.MongoInfo
type MongoModel struct {
	DB *mongodb.MongoInfo
}

// GetNewsData is to get all news data from news collection
func (mg *MongoModel) GetNewsData() ([]News, error) {
	//mg.Db.Session
	mg.DB.GetCol(newsCollection)

	var news []News

	//get
	colQuerier := bson.M{}
	err := mg.DB.C.Find(colQuerier).All(&news)
	if err != nil {
		return nil, err
	}

	return news, nil
}

// GetArticlesData is to get one or all articles from article collection
func (mg *MongoModel) GetArticlesData(newsID int) ([]Articles, error) {
	mg.DB.GetCol(articlesCollection)

	var articles []Articles

	//get
	colQuerier := bson.M{}
	if newsID != 0 {
		colQuerier = bson.M{"news_id": newsID}
	}

	err := mg.DB.C.Find(colQuerier).Sort("news_id").All(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// GetArticlesData2 is to get one or all articles from article2 collection
func (mg *MongoModel) GetArticlesData2(newsID int) ([]Item2, error) {
	mg.DB.GetCol(articles2Collection)

	var items []Item2

	//get
	colQuerier := bson.M{}
	if newsID != 0 {
		colQuerier = bson.M{"news_id": newsID}
	}

	err := mg.DB.C.Find(colQuerier).Sort("news_id").All(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
