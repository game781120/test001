package elastics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
	"strings"
)

func Start(es *EsData) (bool, error) {
	err := (*es).GetES()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (e *EsData) GetES() error {
	if e.ES != nil {
		// 检查连接是否可用
		res, err := e.ES.Ping()
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Elasticsearch client is connected")
			return nil
		}
	}
	var indexName = e.Index

	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s@%s:%s", e.Username, e.Password,
				e.Host, e.Port),
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("elasticsearch.NewClient ", e.Host, e.Port, err)
		return err
	}

	// 检查索引是否存在，如果存在则删除
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		deleteRes, err := es.Indices.Create(indexName)
		if err != nil {
			return err
		}
		if deleteRes.StatusCode != 200 {
			return err
		}
	}

	var resData map[string]interface{}
	resData, err = e.GetMapping(es)
	existingProperties, ok := resData["mappings"].(map[string]interface{})["properties"].(map[string]interface{})
	if !ok {
		fmt.Println("没有找到mapping数据,重新创建")
		e.PutMapping(e.Mapping, es)
		return nil
	}

	var mapping map[string]interface{}
	err = json.Unmarshal([]byte(e.Mapping), &mapping)
	if err != nil {
		fmt.Println("Error unmarshaling mapping:", err)
		return nil
	}

	new_properties, ok := mapping["properties"].(map[string]interface{})
	for field := range existingProperties {
		if _, ok := new_properties[field]; !ok {
			//if v, ok1 := existingProperties[field].(map[string]interface{}); ok1 {
			//	//v["enabled"] = false
			//	fmt.Printf("--------- field:%s v %s \n\n", field, v)
			//}
			//"enabled": false
			delete(existingProperties, field)
		}
	}

	for field, properties := range new_properties {
		if _, ok := existingProperties[field]; !ok {
			existingProperties[field] = properties
		}
	}
	var properties_temp = make(map[string]interface{})
	properties_temp["properties"] = existingProperties
	mappingJSON, err := json.Marshal(properties_temp)
	if err != nil {
		fmt.Println("Error:", err, mappingJSON)
		return err
	}
	e.PutMapping(string(mappingJSON), es)
	e.ES = es
	return nil
}

func (e *EsData) PutMapping(mappingStr string, es *elasticsearch.Client) {
	index := []string{e.Index}
	req := esapi.IndicesPutMappingRequest{
		Index: index,
		Body:  strings.NewReader(mappingStr),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("PutMapping: %s", err)
	}
	defer res.Body.Close()

	// 3. Check response for errors.
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		if errMsg, ok := e["error"].(map[string]interface{})["root_cause"].([]interface{})[0].(map[string]interface{})["reason"].(string); ok {
			log.Fatalf("Elasticsearch error: %s", errMsg)
		}
	}
}

func (e *EsData) GetMapping(es *elasticsearch.Client) (map[string]interface{}, error) {
	index := []string{e.Index}
	req := esapi.IndicesGetMappingRequest{Index: index}
	res, err := req.Do(context.Background(), es)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("Error getting mapping: %s", err)
		return nil, err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
		return nil, err
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	indexName := fmt.Sprintf("%s", e.Index)
	data, ok := responseBody[indexName].(map[string]interface{})
	if !ok {
		log.Fatalf("Error converting map to MyStruct")
		return nil, err
	} else {
		return data, nil
	}
	return nil, err
}

func (e *EsData) indexes(mappingStr string, es *elasticsearch.Client) {
	index := []string{e.Index}
	req := esapi.IndicesPutMappingRequest{
		Index: index,
		Body:  strings.NewReader(mappingStr),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("PutMapping: %s", err)
	}
	defer res.Body.Close()

	// 3. Check response for errors.
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		if errMsg, ok := e["error"].(map[string]interface{})["root_cause"].([]interface{})[0].(map[string]interface{})["reason"].(string); ok {
			log.Fatalf("Elasticsearch error: %s", errMsg)
		}
	}
}

func (e *EsData) DeleteIndex() error {
	// 检查索引是否存在，如果存在则删除
	res, err := e.ES.Indices.Exists([]string{e.Index})
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		deleteRes, err := e.ES.Indices.Delete([]string{e.Index})
		if err != nil {
			return err
		}
		if deleteRes.StatusCode != 200 {
			return fmt.Errorf("Failed to delete index")
		}
	}
	return nil
}
