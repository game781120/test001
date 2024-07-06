package chat

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
)

// File struct represents an knowledge file.
type File struct {
	Size             string `json:"size"`
	CreateTime       string `json:"createTime"`
	Md5              string `json:"md5"`
	OriginalFilename string `json:"originalFilename"`
	PreviewUrl       string `json:"previewUrl"`
	OssBucket        string `json:"bucket"`
	OssPrefix        string `json:"prefix"`
	OssKey           string `json:"key"`
	CustomFileType   string `json:"customFileType"`

	httpHeader
}

type KnowledgeInfoResult struct {
	KnowledgeId   string `json:"knowledge_id"`
	KnowledgeName string `json:"knowledge_name"`
	VisibleState  uint32 `json:"visible_state"`
}

type KnowledgeAddRequest struct {
	// 空间名称 以字母、数字或中文开头，并且仅包含中文字符、字母、数字、下划线和连字符（减号）。长度为1到20个字符
	KnowledgeName string `json:"knowledge_name"`

	// 可见范围 1 # 个人 ； 2 # 团队 ；3 # 企业 ；
	VisibleState uint32 `json:"visible_state"`
}

type KnowledgeListRequest struct {
	// 过滤类型 必需 1#个人 2#团队 3#企业 4# 全部
	FilterType uint32 `json:"filter_type"`

	// 必需 过滤的关键词，实现模糊查询
	Keyword string `json:"keyword"`

	// 查询范围 必需 1#我的空间 2#共享空间 3#收藏 4# 管理
	SearchType uint32 `json:"search_type"`

	PageNum  uint32 `json:"page"`
	PageSize uint32 `json:"size"`
}

type KnowledgeListResult struct {
	MyList    []KnowledgeInfoResult `json:"knowledge_list"`
	AdminList []KnowledgeInfoResult `json:"admin_list"`
	ShareList []KnowledgeInfoResult `json:"share_list"`
}

type KnowledgeFileLearningRequest struct {
	KnowledgeId string                      `json:"knowledge_id"`
	Type        uint32                      `json:"type"`
	Files       []KnowledgeFileLearningInfo `json:"files"`
}

type KnowledgeFileLearningInfo struct {
	Size           string `json:"size"`
	Name           string `json:"name"`
	OriginUrl      string `json:"originUrl"`
	CustomFileType string `json:"customFileType"`
}

type KnowledgeFileLearningResult struct {
	FileId    string `json:"fileId"`
	Name      string `json:"name"`
	OriginUrl string `json:"originUrl"`
}

type KnowledgeFileListResult struct {
	Remark           string `json:"remark"`
	KnowledgeId      string `json:"datasetId"`
	Name             string `json:"name"`
	CleanStatus      string `json:"cleanStatus"`
	CleanErrorReason string `json:"cleanErrorReason"`
	OriginUrl        string `json:"originUrl"`
	SliceStatus      string `json:"sliceStatus"`
	SliceErrorReason string `json:"sliceErrorReason"`
	CustomFileType   string `json:"customFileType"`
	DirectoryId      string `json:"directoryId"`
	FileSource       string `json:"fileSource"`
	FileExt          string `json:"fileExt"`
	UserId           string `json:"userId"`
	DatasetName      string `json:"datasetName"`
	Size             string `json:"size"`
	FileId           string `json:"fileId"`
}

type KnowledgeFileDeleteRequest struct {
	KnowledgeId string   `json:"knowledge_id"`
	FileIds     []string `json:"file_ids"`
}

type GetKnowledgeInfoRequest struct {
	KnowledgeId string `json:"knowledge_id"`
}

type GetKnowledgeInfoResult struct {
	FileSize int    `json:"fileSize"`
	Name     string `json:"name"`
}

// UploadFile uploads a file to knowledge
// FilePath must be a local file path.
func (c *Client) UploadFile(ctx context.Context,
	filepaths ...string,
) (files []File, err error) {
	var b bytes.Buffer
	builder := c.createFormBuilder(&b)
	for _, filepath := range filepaths {
		var fileData *os.File
		fileData, err = os.Open(filepath)
		if err != nil {
			return
		}

		err = builder.CreateFormFile("files", fileData)
		if err != nil {
			return
		}
	}
	err = builder.Close()
	if err != nil {
		return
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL("/brain/dandelion/api/v1/file"),
		withBody(&b), withContentType(builder.FormDataContentType()))
	if err != nil {
		return
	}

	var resp CommonResponse[[]File, string]
	err = c.sendRequest(req, &resp)

	return resp.Data, err
}

func (c *Client) CreateKnowledge(ctx context.Context,
	request KnowledgeAddRequest,
) (response KnowledgeInfoResult, err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/brain/knowledge/api/knowledge/pub_create"),
		withBody(request),
	)
	if err != nil {
		return
	}

	var resp CommonResponse[KnowledgeInfoResult, uint32]
	err = c.sendRequest(req, &resp)
	return resp.Data, err
}

func (c *Client) KnowledgeList(ctx context.Context,
	request KnowledgeListRequest,
) (response KnowledgeListResult, err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/brain/knowledge/api/knowledge/pub_knowledge_list"),
		withBody(request),
	)
	if err != nil {
		return
	}

	var resp CommonResponse[KnowledgeListResult, uint32]
	err = c.sendRequest(req, &resp)
	return resp.Data, err
}

func (c *Client) DeleteKnowledge(ctx context.Context,
	knowledgeId string,
) (err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(fmt.Sprintf("/brain/knowledge/api/knowledge/pub_delete_knowledge/%s", knowledgeId)),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, nil)
	return
}

func (c *Client) Learning(ctx context.Context,
	request KnowledgeFileLearningRequest,
) (result []KnowledgeFileLearningResult, err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/brain/knowledge/api/file/pub_learning"),
		withBody(request),
	)
	if err != nil {
		return
	}
	var resp CommonResponse[[]KnowledgeFileLearningResult, uint32]
	err = c.sendRequest(req, &resp)
	return resp.Data, err
}

func (c *Client) Files(ctx context.Context,
	userId string,
) (response []KnowledgeFileListResult, err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		c.fullURL("/brain/dandelion/api/v1/dataset/files/listByUserId?userId="+userId),
	)
	if err != nil {
		return
	}

	var resp CommonResponse[[]KnowledgeFileListResult, string]
	err = c.sendRequest(req, &resp)
	return resp.Data, err
}

func (c *Client) DeleteFiles(ctx context.Context,
	request KnowledgeFileDeleteRequest,
) (err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(fmt.Sprintf("/brain/knowledge/api/knowledge/pub_delete_files?knowledge_id=%s&action=delete", request.KnowledgeId)),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, nil)
	return
}

func (c *Client) GetKnowledgeInfo(ctx context.Context,
	knowledgeId string,
) (response GetKnowledgeInfoResult, err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		c.fullURL(fmt.Sprintf("/brain/knowledge/api/search/%s", knowledgeId)),
	)
	if err != nil {
		return
	}

	var resp CommonResponse[GetKnowledgeInfoResult, uint32]
	err = c.sendRequest(req, &resp)
	return resp.Data, err
}
