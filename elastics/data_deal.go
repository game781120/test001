package elastics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/mitchellh/mapstructure"
	"log"

	"strings"
)

func (e *EsData) DeleteRequest(esStruct EsStruct) {
	err := e.GetES()
	if err == nil {
		log.Printf("Error getting Elasticsearch client: %s", err)
		return
	}

	// 准备索引请求
	req := esapi.DeleteRequest{
		Index:      e.Index,
		DocumentID: esStruct.DocumentID,
		Refresh:    "true",
	}
	// 执行索引请求
	res, err := req.Do(context.Background(), e.ES)
	if err != nil {
		log.Printf("Error indexing document: %s", err)
		return
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.IsError() {
		log.Printf("Error indexing document: %s", res.Status())
		return
	}

	// 打印成功消息
	fmt.Println("DeleteRequest successfully!")

}

func (e *EsData) IndexRequest(esStruct EsStruct) error {
	err := e.GetES()
	if err != nil {
		log.Printf("Error getting Elasticsearch client: %s", err)
		return err
	}
	// 将文档转换为 JSON 字节
	docBytes, err0 := json.Marshal(esStruct.BusinessData)
	if err0 != nil {
		log.Printf("Error marshaling document: %s", err0)
		return err0
	}
	docBytesStr := string(docBytes)
	//log.Printf("----------- docBytesStr: %s", docBytesStr)

	// 准备索引请求
	req := esapi.IndexRequest{
		Index:      e.Index,
		DocumentID: esStruct.DocumentID,
		Body:       strings.NewReader(docBytesStr),
		Refresh:    "true",
	}
	// 执行索引请求
	res, err := req.Do(context.Background(), e.ES)
	if err != nil {
		log.Printf("Error indexing document: %s", err)
		return err
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.IsError() {
		log.Printf("Error indexing document: %s", res.Status())
		return err
	}

	// 打印成功消息
	fmt.Println("Document indexed successfully!")
	return nil
}

func (e *EsData) DealReq(req *esapi.SearchRequest) ([]map[string]interface{}, error) {
	//// 执行查询
	res, err := req.Do(context.Background(), e.ES)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	// 检查响应状态
	if res.IsError() {
		log.Printf("Error response: %s", res.String())
		return nil, err
	}
	// 读取响应体
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil, err
	}
	resSlices := make([]map[string]interface{}, 0)
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		resData := make(map[string]interface{})
		err := mapstructure.Decode(source, &resData)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		resSlices = append(resSlices, resData)
	}
	return resSlices, nil
}

func (e *EsData) QueryCommon(queryStr string, size int) ([]map[string]interface{}, error) {

	err := e.GetES()
	if err == nil {
		log.Printf("Error getting Elasticsearch client: %s", err)
		return nil, err
	}
	req := esapi.SearchRequest{
		Index: []string{e.Index},           // 替换为你的索引名
		Body:  strings.NewReader(queryStr), // 替换为你的查询体
		Size:  &size,
	}
	// 读取响应体
	return e.DealReq(&req)
}

func (e *EsData) QueryBySort(queryStr, sortField string, size int) ([]map[string]interface{}, error) {

	err := e.GetES()
	if err == nil {
		log.Printf("Error getting Elasticsearch client: %s", err)
		return nil, err
	}
	req := esapi.SearchRequest{
		Index: []string{e.Index},           // 替换为你的索引名
		Body:  strings.NewReader(queryStr), // 替换为你的查询体
		Size:  &size,
		Sort:  []string{fmt.Sprintf("%s:desc", sortField)},
	}
	// 读取响应体
	return e.DealReq(&req)
}

//func (e *EsData)QueryRequest(ts int64, keys string) ([]map[string]interface{}, error) {
//
//	var esQuery EsQuery
//	esQuery.Query.Bool.Must = make([]map[string]interface{}, 0)
//	var timestamp = map[string]interface{}{"gt": ts}
//	var rangeData = map[string]interface{}{"timestamp": timestamp}
//	esQuery.Query.Bool.Must = append(esQuery.Query.Bool.Must, map[string]interface{}{"range": rangeData})
//	if keys != "" {
//		var content = map[string]interface{}{"content": keys}
//		esQuery.Query.Bool.Must = append(esQuery.Query.Bool.Must, map[string]interface{}{"match": content})
//	}
//	// 将查询条件转换为 JSON 字符串
//	queryBytes, err := json.Marshal(esQuery)
//	if err != nil {
//		fmt.Println("转换查询条件失败:", err)
//		return nil, err
//	}
//
//
//	err = e.GetES()
//	if err == nil {
//		log.Printf("Error getting Elasticsearch client: %s", err)
//		return nil, err
//	}
//	//创建一个查询请求
//	size := 1000
//	req := esapi.SearchRequest{
//		Index: []string{e.Index},                       // 替换为你的索引名
//		Body:  strings.NewReader(string(queryBytes)), // 替换为你的查询体
//		Size:  &size,
//	}
//
//    return e.DealReq(&req)
//}
//
//func QueryByUuids(index, uuidsStr string) ([]utils.EmailVO, error) {
//	emails := make([]utils.EmailVO, 0)
//	es, err := GetES()
//	if es == nil {
//		log.Printf("Error getting Elasticsearch client: %s", err)
//		return nil, err
//	}
//
//	queryData := strings.Replace(utils.QueryByUuids, "UUIDS", uuidsStr, 1)
//	// log.Printf("queryData: %s", queryData)
//	//创建一个查询请求
//	size := 200
//	req := esapi.SearchRequest{
//		Index: []string{index},              // 替换为你的索引名
//		Body:  strings.NewReader(queryData), // 替换为你的查询体
//		Size:  &size,
//	}
//
//	//// 执行查询
//	res, err0 := req.Do(context.Background(), es)
//	if err0 != nil {
//		log.Printf("Error getting response: %s", err0)
//	}
//	defer res.Body.Close()
//
//	// 检查响应状态
//	if res.IsError() {
//		log.Printf("Error response: %s", res.String())
//		return nil, err
//	}
//	// 读取响应体
//	var r map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		log.Printf("Error parsing the response body: %s", err)
//		return nil, err
//	}
//	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
//	for _, hit := range hits {
//		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
//		var emailvo = utils.EmailVO{}
//		err := mapstructure.Decode(source, &emailvo)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return nil, err
//		}
//		emails = append(emails, emailvo)
//		//log.Printf("emailvo: %s", emailvo)
//	}
//	return emails, nil
//}
//
//func QueryAll(index string) ([]utils.EmailVO, error) {
//	emails := make([]utils.EmailVO, 0)
//	es, err := GetES()
//	if es == nil {
//		log.Printf("Error getting Elasticsearch client: %s", err)
//		return nil, err
//	}
//
//	//log.Printf("utils.QueryAall: %s", utils.QueryAall)
//	//创建一个查询请求
//	size := 1000
//	req := esapi.SearchRequest{
//		Index: []string{index},                    // 替换为你的索引名
//		Body:  strings.NewReader(utils.QueryAall), // 替换为你的查询体
//		Size:  &size,
//	}
//
//	//// 执行查询
//	res, err0 := req.Do(context.Background(), es)
//	if err0 != nil {
//		log.Printf("Error getting response: %s", err0)
//	}
//	defer res.Body.Close()
//
//	// 检查响应状态
//	if res.IsError() {
//		log.Printf("Error response: %s", res.String())
//		return nil, err
//	}
//	// 读取响应体
//	var r map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		log.Printf("Error parsing the response body: %s", err)
//		return nil, err
//	}
//	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
//	for _, hit := range hits {
//		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
//		var emailvo = utils.EmailVO{}
//		//log.Printf("emailvo: %s", emailvo)
//		err := mapstructure.Decode(source, &emailvo)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return nil, err
//		}
//		emails = append(emails, emailvo)
//		//log.Printf("emailvo: %s", emailvo)
//	}
//	return emails, nil
//}
//
//func QueryNoRead(index, keys string, level, size int) ([]utils.EmailVO, error) {
//
//	var esQuery utils.EsQuery
//	esQuery.Query.Bool.Must = make([]map[string]interface{}, 0)
//	var term = map[string]interface{}{"isRead": false}
//	esQuery.Query.Bool.Must = append(esQuery.Query.Bool.Must, map[string]interface{}{"term": term})
//	if level != -1 {
//		var levelData = map[string]interface{}{"level": level}
//		esQuery.Query.Bool.Must = append(esQuery.Query.Bool.Must, map[string]interface{}{"term": levelData})
//	}
//	if keys != "" {
//		var subject = map[string]interface{}{"subject": keys}
//		esQuery.Query.Bool.Must = append(esQuery.Query.Bool.Must, map[string]interface{}{"match": subject})
//	}
//
//	// 将查询条件转换为 JSON 字符串
//	queryBytes, err := json.Marshal(esQuery)
//	if err != nil {
//		fmt.Println("转换查询条件失败:", err)
//		return nil, err
//	}
//	log.Printf("queryBytes: %s", string(queryBytes))
//
//	es, err := GetES()
//	if es == nil {
//		log.Printf("Error getting Elasticsearch client: %s", err)
//		return nil, err
//	}
//	req := esapi.SearchRequest{
//		Index: []string{index},                       // 替换为你的索引名
//		Body:  strings.NewReader(string(queryBytes)), // 替换为你的查询体
//		Size:  &size,
//		Sort:  []string{"uid:desc"},
//	}
//	//// 执行查询
//	res, err0 := req.Do(context.Background(), es)
//	if err0 != nil {
//		log.Printf("Error getting response: %s", err0)
//	}
//	defer res.Body.Close()
//
//	// 检查响应状态
//	if res.IsError() {
//		log.Printf("Error response: %s", res.String())
//		return nil, err
//	}
//	// 读取响应体
//	var r map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		log.Printf("Error parsing the response body: %s", err)
//		return nil, err
//	}
//	emails := make([]utils.EmailVO, 0)
//	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
//	for _, hit := range hits {
//		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
//		var emailvo = utils.EmailVO{}
//		err := mapstructure.Decode(source, &emailvo)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return nil, err
//		}
//		emails = append(emails, emailvo)
//	}
//	return emails, nil
//}
//
//func QueryMeeting(index string) ([]utils.EmailVO, error) {
//	var esQuery utils.EsQuery
//	esQuery.Query.Bool.Must = make([]map[string]interface{}, 0)
//	var term = map[string]interface{}{"isMeeting": true}
//	esQuery.Query.Bool.Must = append(esQuery.Query.Bool.Must, map[string]interface{}{"term": term})
//
//	// 将查询条件转换为 JSON 字符串
//	queryBytes, err := json.Marshal(esQuery)
//	if err != nil {
//		fmt.Println("转换查询条件失败:", err)
//		return nil, err
//	}
//	log.Printf("queryBytes: %s", string(queryBytes))
//
//	es, err := GetES()
//	if es == nil {
//		log.Printf("Error getting Elasticsearch client: %s", err)
//		return nil, err
//	}
//	size := 10
//	req := esapi.SearchRequest{
//		Index: []string{index},                       // 替换为你的索引名
//		Body:  strings.NewReader(string(queryBytes)), // 替换为你的查询体
//		Size:  &size,
//		Sort:  []string{"uid:desc"},
//	}
//	//// 执行查询
//	res, err0 := req.Do(context.Background(), es)
//	if err0 != nil {
//		log.Printf("Error getting response: %s", err0)
//	}
//	defer res.Body.Close()
//
//	// 检查响应状态
//	if res.IsError() {
//		log.Printf("Error response: %s", res.String())
//		return nil, err
//	}
//	// 读取响应体
//	var r map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		log.Printf("Error parsing the response body: %s", err)
//		return nil, err
//	}
//	emails := make([]utils.EmailVO, 0)
//	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
//	for _, hit := range hits {
//		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
//		var emailvo = utils.EmailVO{}
//		err := mapstructure.Decode(source, &emailvo)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return nil, err
//		}
//		emails = append(emails, emailvo)
//	}
//	return emails, nil
//}
//
//// CheckMeeting 检查会议是否已经预定,并创建日程
//func QueryMeetingData(isEnglish bool) ([]utils.EmailVO, error) {
//	emails, err := QueryMeeting(conf.ConfigInfo.Elastic.Index)
//	if err != nil {
//		fmt.Printf("查询es 错误 %s:", err)
//		return nil, err
//	}
//	resEmails := make([]utils.EmailVO, 0)
//	for _, email := range emails {
//		if email.IsMeeting == false {
//			continue
//		}
//		// 如果还没有预定会议，则跳过
//		if email.ReserveInfo == "" {
//			continue
//		}
//
//		isEn := utils.IsEnglish(email.Subject)
//		if isEnglish != isEn {
//			continue
//		}
//
//		var reserve utils.Reserve
//		err := json.Unmarshal([]byte(email.ReserveInfo), &reserve)
//		if err != nil {
//			fmt.Println("解析 JSON 失败:", err)
//			continue
//		}
//		sec, _ := strconv.ParseInt(reserve.EndTime, 10, 64)
//		lastTime := time.Unix(sec, 0)
//		if sec <= time.Now().Unix() {
//			fmt.Printf("会议时间已过期 现在时间: %s 会议时间: %s\n", time.Now().Format(utils.Layout01), lastTime.Format(utils.Layout01))
//			continue
//		}
//
//		eventResponse := feishu.CreateEvents(email.MeetingStruct, email.ReserveInfo, isEnglish)
//		if eventResponse == nil {
//			fmt.Println("创建日程失败")
//			continue
//		}
//		if eventResponse.Data.Event.OrganizerCalendarID == "" || eventResponse.Data.Event.EventID == "" {
//			fmt.Println("创建日程失败---")
//			continue
//		}
//		feishu.Attendees(eventResponse.Data.Event.OrganizerCalendarID, eventResponse.Data.Event.EventID)
//		email.IsNotify = true
//		err0 := IndexRequest(conf.ConfigInfo.Elastic.Index, email)
//		if err0 == nil {
//			resEmails = append(resEmails, email)
//		}
//	}
//	return resEmails, nil
//}
//
//func QueryLastOne(index string) (string, error) {
//
//	es, err := GetES()
//	if es == nil {
//		log.Printf("Error getting Elasticsearch client: %s", err)
//		return "", err
//	}
//	size := 1
//	req := esapi.SearchRequest{
//		Index: []string{index}, // 替换为你的索引名
//		Size:  &size,
//		Sort:  []string{"uid:desc"},
//	}
//	//// 执行查询
//	res, err0 := req.Do(context.Background(), es)
//	if err0 != nil {
//		log.Printf("Error getting response: %s", err0)
//		return "", err
//	}
//	defer res.Body.Close()
//
//	// 检查响应状态
//	if res.IsError() {
//		log.Printf("Error response: %s", res.String())
//		return "", err
//	}
//	// 读取响应体
//	var r map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		log.Printf("Error parsing the response body: %s", err)
//		return "", err
//	}
//	var em utils.EmailVO
//	em.UID = "none"
//	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
//	for _, hit := range hits {
//		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
//
//		err := mapstructure.Decode(source, &em)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return "", err
//		}
//	}
//	return em.UID, nil
//}
