package elastics

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type businessData interface{}

type EsStruct struct {
	DocumentID   string
	BusinessData businessData
}

type EsData struct {
	Host        string
	Port        string
	Index       string
	Username    string
	Password    string
	ES          *elasticsearch.Client
	TableStruct EsStruct
	Mapping     string
}

var Mapping001 = `{
		"properties": {
			"uid": {
				"type": "keyword",
                "index": true
			},
            "from": {
				"type": "keyword"
			},
            "fromAddr": {
				"type": "keyword"
			},
			"date": {
        		"type":   "date",
        		"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      		},
            "timestamp": {
                "type":   "date",
				"format": "epoch_millis"
            },
            "subject": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"content": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"to": {
				"type": "keyword"
			},
            "cc": {
				"type": "keyword"
			},
            "level": {
                "type": "keyword"
            },
            "isMeeting": {
                "type": "boolean"
            }
        }        
     }`

type EsQuery struct {
	Query struct {
		Bool struct {
			Must   []map[string]interface{} `json:"must"`
			Should []map[string]interface{} `json:"should"`
		} `json:"bool"`
	} `json:"query"`
}
