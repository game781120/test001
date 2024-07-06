package openai_test

import (
	"context"
	"testing"

	"thundersoft.com/brainos/openai"
	"thundersoft.com/brainos/openai/internal/test/checks"
)

func TestListModel(t *testing.T) {
	client := openai.NewClient("http://10.0.36.13:8888", "KJInf01E1p5Q1zvn65704c7501Ef4e83B85aB44a0128E5Dd")
	resp, err := client.ListModels(context.Background())

	for _, model := range resp.Content {
		t.Log(model.Model)
	}

	checks.NoError(t, err, "TestListModel error")
}
