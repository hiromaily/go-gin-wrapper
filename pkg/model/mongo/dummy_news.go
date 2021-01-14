package mongo

// DummyMongo is dummy object
type DummyMongo struct{}

// GetNewsData is to get all news data from news collection
func (d *DummyMongo) GetNewsData() ([]News, error) {
	return nil, nil
}

// GetArticlesData is to get one or all articles from article collection
func (d *DummyMongo) GetArticlesData(newsID int) ([]Articles, error) {
	return nil, nil
}

// GetArticlesData2 is to get one or all articles from article2 collection
func (d *DummyMongo) GetArticlesData2(newsID int) ([]Item2, error) {
	return nil, nil
}
