package models

import (
	lg "github.com/hiromaily/golibs/log"
	//"time"
)

// News
type News struct {
	//NewsId    int       `column:"news_id"          db:"news_id"`
	Article string `column:"article"       db:"article"`
	//DeleteFlg string    `column:"delete_flg"       db:"delete_flg"`
	//Created   time.Time `column:"create_datetime"  db:"create_datetime"`
	//Updated   time.Time `column:"update_datetime"  db:"update_datetime"`
}

// Get News List
func (self *Models) GetNewsList() ([]map[string]interface{}, error) {
	sql := "SELECT news_id, article, create_datetime FROM t_news WHERE delete_flg='0' ORDER BY create_datetime DESC"
	//TODO:handle correctly even if no parameter
	//TODO:commonalize below code for models package.
	data, _, err := self.Db.SelectSQLAllField(sql, nil)
	if err != nil {
		lg.Errorf("SQL may be wrong. : %s\n", err.Error())
		return nil, err
	} else if len(data) == 0 {
		lg.Info("No data.")
		return nil, nil
	}
	return data, nil
}
