package myHttp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"thundersoft.com/brain/DigitalVisitor/conf"
	"thundersoft.com/brain/DigitalVisitor/nacos"
	"thundersoft.com/brain/DigitalVisitor/utils"
)

func KnowledgeCallData(question string, knowledgeIds, authorization string) []utils.KnowledgeData {
	reData := make(map[string]interface{})
	reData["question"] = question
	reData["top"] = 200
	knowledge_ids := make([]string, 0)
	knowledge_ids = append(knowledge_ids, knowledgeIds)
	reData["knowledge_ids"] = knowledge_ids
	queryBytes, err := json.Marshal(reData)
	if err != nil {
		fmt.Println("转换查询条件失败:", err)
		return nil
	}
	ip, port := nacos.GetKnowlede()
	if ip == "" || port == 0 {
		return nil
	}
	url := fmt.Sprintf("http://%s:%d/api/search/poc_search", ip, port)
	// 准备 POST 请求的数据
	payload := strings.NewReader(string(queryBytes))

	// 创建一个 HTTP 请求的客户端
	client := &http.Client{}
	// 创建一个 HTTP 请求
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// 设置 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)
	fmt.Println("KnowledgeCallData-------", url)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()
	fmt.Println("KnowledgeCallData-------resp ")
	body, _ := io.ReadAll(resp.Body)
	var response utils.KnowledgeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil
	}
	if response.Code != 1 {
		fmt.Println("未召回数据 KnowledgeCallData :", response.Msg)
		return nil
	} else {
		fmt.Println("知识库召回数据 KnowledgeCallData :", response.Data.DataList)
	}
	return response.Data.DataList
}

// 微软 embedding
func EmbeddingAzure(question string) []float64 {
	reData := make(map[string]interface{})
	reData["input"] = question
	queryBytes, err := json.Marshal(reData)
	if err != nil {
		fmt.Println("转换查询条件失败:", err)
		return nil
	}
	url := fmt.Sprintf("https://xxxxxxxx/openai/deployments/brain-embedding/embeddings?api-version=2024-02-01")
	// 准备 POST 请求的数据
	payload := strings.NewReader(string(queryBytes))
	// 创建一个 HTTP 请求的客户端
	client := &http.Client{}
	// 创建一个 HTTP 请求
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// 设置 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", "sssssssssssss")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error sending request:", resp.StatusCode, " err=", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response utils.EmbeddingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil
	}

	return response.Data[0].Embedding
}

// ThunderSoft embedding
func EmbeddingThunderSoft(question string) ([][]float32, error) {
	reData := make(map[string]interface{})
	reData["texts"] = []string{question}
	reData["model_name"] = conf.ConfigInfo.Embedding.ModelName
	reData["model"] = conf.ConfigInfo.Embedding.ModelName
	queryBytes, err := json.Marshal(reData)
	if err != nil {
		fmt.Println("转换查询条件失败:", err)
		return nil, err
	}
	url := conf.ConfigInfo.Embedding.Address
	// 准备 POST 请求的数据
	payload := strings.NewReader(string(queryBytes))
	// 创建一个 HTTP 请求的客户端
	client := &http.Client{}
	// 创建一个 HTTP 请求
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// 设置 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", conf.ConfigInfo.Embedding.Authorization)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error sending request:", resp.StatusCode, " err=", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response utils.EmbeddingResponseEx
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}

	return response.Data.Value, nil
}

// JunShiQuestion 拒识请求
//func JunShiQuestion(junShiReq utils.JunShiReq) (*utils.JunShiRes, error) {
//	reqBytes, err := json.Marshal(junShiReq)
//	if err != nil {
//		fmt.Println("JunShiQuestion:", err)
//		return nil, err
//	}
//	// 准备 POST 请求的数据
//	payload := strings.NewReader(string(reqBytes))
//	req, err := http.NewRequest("POST", conf.ConfigInfo.JuShi.Url, payload)
//	if err != nil {
//		fmt.Println("Error creating request:", err)
//		return nil, err
//	}
//	// 设置 Header
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	// 发送请求
//	resp, err := client.Do(req)
//	if err != nil || resp.StatusCode != 200 {
//		fmt.Println("Error sending request:", resp.StatusCode, " err=", err)
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, _ := io.ReadAll(resp.Body)
//	var response utils.JunShiRes
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		fmt.Println("解析JSON数据失败:", err)
//		return nil, err
//	}
//
//	return &response, nil
//}

func JuShiQuestion(juShiReq utils.JunShiReq) (*utils.JunShiRes, error) {
	juShiReq.Model = conf.ConfigInfo.JuShi.ModelName
	reqBytes, err := json.Marshal(juShiReq)
	if err != nil {
		fmt.Println("JunShiQuestion:", err)
		return nil, err
	}
	// 准备 POST 请求的数据
	payload := strings.NewReader(string(reqBytes))
	req, err := http.NewRequest("POST", conf.ConfigInfo.JuShi.Url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	// 设置 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", conf.ConfigInfo.JuShi.Authorization)
	client := &http.Client{}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error sending request:", resp.StatusCode, " err=", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var response utils.JunShiRes
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}

	return &response, nil
}
func SoftQuestion(safeReq utils.SafeReq) (*utils.SafeRes, error) {
	u := uuid.New()
	safeReq.RequestID = u.String()
	safeReq.Model = conf.ConfigInfo.JuShi.ModelName
	reqBytes, err := json.Marshal(safeReq)
	if err != nil {
		fmt.Println("JunShiQuestion:", err)
		return nil, err
	}
	// 准备 POST 请求的数据
	payload := strings.NewReader(string(reqBytes))
	req, err := http.NewRequest("POST", conf.ConfigInfo.JuShi.Url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	// 设置 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", conf.ConfigInfo.JuShi.Authorization)
	client := &http.Client{}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error sending request:", resp.StatusCode, " err=", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response utils.SafeRes
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}

	return &response, nil
}

// NluQuestion 意图识别
func NluQuestion(content string) (*utils.NluRes, error) {
	var nluReq utils.NluReq
	nluReq.Question = content
	u := uuid.New()
	nluReq.RequestId = u.String()

	reqBytes, err := json.Marshal(nluReq)
	if err != nil {
		fmt.Println("JunShiQuestion:", err)
		return nil, err
	}
	// 准备 POST 请求的数据
	payload := strings.NewReader(string(reqBytes))
	req, err := http.NewRequest("POST", conf.ConfigInfo.Nlu.Host, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	// 设置 Header
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error sending request:", resp.StatusCode, " err=", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var response utils.NluRes
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return nil, err
	}

	return &response, nil
}
