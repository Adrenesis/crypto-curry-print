package model

import "encoding/json"

type CoinDataMapUnmarshaler struct {
	CoinDataMap map[string]json.RawMessage `json:"data"`
}

type CoinDataMap struct {
	CoinDataMap []CoinDatumMap
}

type CoinDatumMap struct {
	Id int64 `json:"id"`
	//Name		string   	`json:"name"`
	//Symbol		string   	`json:"symbol"`
	//DateAdded	string   	`json:"date_added"`
	//Category	string		`json:"category"`
	//Description	string		`json:"description"`
	//Slug		string		`json:"slug"`
	//Logo		string		`json:"logo"`
	//Subreddit	string		`json:"subreddit"`
	//Tags		[]string	`json:"tags"`
	////TagsName	[]string	`json:"tags-names"`
	////TagsGroups	[]string	`json:"tags-groups"`
	URLs Urls `json:"urls"`
}

type Urls struct {
	//Chat		 []string	`json:"chat"`
	//Facebook	 []string	`json:"facebook"`
	//MessageBoard []string	`json:"message_board"`
	Explorer []string `json:"explorer"`
	//Twitter		 []string	`json:"twitter"`
	//Website 	 []string	`json:"website"`
	//Technical	 []string	`json:"technical_doc"`
	//SourceCode	 []string	`json:"source_code"`
	//Announcement []string	`json:"announcement"`
}
