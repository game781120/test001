package chat

import (
	"context"
	"testing"

	"thundersoft.com/brainos/chat/internal/test/checks"
)

func TestUploadFile(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.UploadFile(context.Background(), "/home/ehlxr/codes/rubik-brainos-sdk/go-sdk/chat-sdk/README.md")
	for _, file := range resp {
		t.Log(file.OriginalFilename, file.PreviewUrl, file.Size, file.CustomFileType)
	}

	checks.NoError(t, err, "TestCreateFile error")
}

func TestCreateKnowledge(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.CreateKnowledge(context.Background(), KnowledgeAddRequest{
		KnowledgeName: "test",
		VisibleState:  1,
	})
	t.Log(resp.KnowledgeId, resp.KnowledgeName, resp.VisibleState)

	checks.NoError(t, err, "TestCreateKnowledge error")
}

func TestKnowledgeList(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.KnowledgeList(context.Background(), KnowledgeListRequest{})
	for _, v := range resp.MyList {
		t.Log(v.KnowledgeId, v.KnowledgeName, v.VisibleState)
	}

	checks.NoError(t, err, "TestKnowledgeList error")
}

func TestDeleteKnowledge(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	err := client.DeleteKnowledge(context.Background(), "1767799159370416129")
	checks.NoError(t, err, "TestDeleteKnowledge error")
}

func TestLearning(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	_, err := client.Learning(context.Background(), KnowledgeFileLearningRequest{
		KnowledgeId: "1765698335274041345",
		Files: []KnowledgeFileLearningInfo{{
			OriginUrl:      "http://10.0.36.13:9000/dandelion/test/20240308/d41d8cd98f00b204e9800998ecf8427e.md",
			Size:           "0",
			CustomFileType: "NORMAL",
			Name:           "README.md",
		}},
	})
	checks.NoError(t, err, "TestLearning error")
}

func TestFiles(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.Files(context.Background(), "test-11111")

	for _, v := range resp {
		t.Log(v.KnowledgeId, v.FileId, v.Name, v.Size, v.CustomFileType)
	}

	checks.NoError(t, err, "TestLearning error")
}

func TestDeleteFiles(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	err := client.DeleteFiles(context.Background(), KnowledgeFileDeleteRequest{
		KnowledgeId: "1765698335274041345",
		FileIds:     []string{"1765929183516688385"},
	})
	checks.NoError(t, err, "TestDeleteFiles error")
}
