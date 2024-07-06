package myMilvus

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"thundersoft.com/brain/DigitalVisitor/utils"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"thundersoft.com/brain/DigitalVisitor/conf"
	"thundersoft.com/brain/DigitalVisitor/myHttp"
)

const (
	idCol, contentCol, answerCol, contentVector, qaCol, imageUrlsCol, videoUrlsCol, relatedQuestionsCol = "ID", "content", "answer", "contentVector", "qa", "imageUrls", "videoUrls", "relatedQuestions"
)

var milvusClient client.Client

func Init() error {

	collName := conf.ConfigInfo.Milvus.CollectionName
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()

	var err error
	milvusClient, err = client.NewClient(ctx, client.Config{
		Address:  conf.ConfigInfo.Milvus.Address,
		Username: conf.ConfigInfo.Milvus.User,
		Password: conf.ConfigInfo.Milvus.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to milvus: %w", err)
	}
	slog.Info("connect to milvus success", "addr", conf.ConfigInfo.Milvus.Address)

	dbs, err := milvusClient.ListDatabases(ctx)
	if err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}
	haveDb := false
	for _, v := range dbs {
		if v.Name == conf.ConfigInfo.Milvus.DbName {
			log.Println("database have already exist", "db", conf.ConfigInfo.Milvus.DbName)
			haveDb = true
			break
		}
	}
	if !haveDb {
		if err := milvusClient.CreateDatabase(ctx, conf.ConfigInfo.Milvus.DbName); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		slog.Info("create database success", "db", conf.ConfigInfo.Milvus.DbName)
	}

	if err := milvusClient.UsingDatabase(ctx, conf.ConfigInfo.Milvus.DbName); err != nil {
		return fmt.Errorf("failed to useing database %s: %w", conf.ConfigInfo.Milvus.DbName, err)
	}
	slog.Info("useing database success", "db", conf.ConfigInfo.Milvus.DbName)

	collExists, err := milvusClient.HasCollection(ctx, collName)
	if err != nil {
		return fmt.Errorf("failed to check collection exists: %w", err)
	}
	if !collExists {
		schema := entity.NewSchema().WithName(collName).WithDescription("this is the basic example collection").
			// currently primary key field is compulsory, and only int64 is allowed
			WithField(entity.NewField().WithName(idCol).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
			WithField(entity.NewField().WithName(contentVector).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(conf.ConfigInfo.Milvus.Dim))).
			WithField(entity.NewField().WithName(contentCol).WithDataType(entity.FieldTypeVarChar).WithMaxLength(65535))

		err = milvusClient.CreateCollection(ctx, schema, entity.DefaultShardNumber)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
		slog.Info("create collection success", "name", collName)

		idx, err := entity.NewIndexHNSW(entity.L2, 8, 64)
		if err != nil {
			return fmt.Errorf("fail to create hnsw index: %w", err)
		}
		err = milvusClient.CreateIndex(ctx, collName, contentVector, idx, false)
		if err != nil {
			return fmt.Errorf("fail to create index: %w", err)
		}
		slog.Info("create index success", "idx", idx)

		err = milvusClient.LoadCollection(ctx, collName, false)
		if err != nil {
			return fmt.Errorf("failed to load collection: %w", err)
		}
		slog.Info("load collection completed")
	}
	slog.Info("check collection success", "name", collName)
	return nil
}

func Insert(emailId int64, content string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	vectors, err := myHttp.EmbeddingThunderSoft(content)
	if err != nil {
		return fmt.Errorf("embeding %s error: %w", content, err)
	}
	columns := []entity.Column{
		entity.NewColumnInt64("ID", []int64{emailId}),
		entity.NewColumnVarChar("content", []string{content}),
		entity.NewColumnFloatVector("embeddings", 1024, vectors),
	}

	collName := conf.ConfigInfo.Milvus.CollectionName
	r, err := milvusClient.Insert(ctx, collName, "", columns...)
	if err != nil {
		return fmt.Errorf("failed to insert film data: %w", err)
	}

	//err = milvusClient.Flush(ctx, collName, false)
	//if err != nil {
	//	return fmt.Errorf("failed to flush collection: %w", err)
	//}
	slog.Info("flush completed")
	slog.Info("insert completed", "id", r.(*entity.ColumnInt64).Data()[0])

	return nil
}

func Query(question string, datasetIds []string) ([]utils.QaRes, error) {
	resList := make([]utils.QaRes, 0)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()

	vectors, err := myHttp.EmbeddingThunderSoft(question)
	if err != nil {
		fmt.Printf("embeding %s error: %w", question, err)
		return resList, err
	}
	expr := ""
	if len(datasetIds) > 0 {
		expr = fmt.Sprintf("datasetId in %s", datasetIds)
	}

	vector := entity.FloatVector(vectors[0])
	sp, _ := entity.NewIndexFlatSearchParam()
	sr, err := milvusClient.Search(ctx, conf.ConfigInfo.Milvus.CollectionName, []string{}, expr,
		[]string{contentCol, answerCol, qaCol, imageUrlsCol, videoUrlsCol, relatedQuestionsCol}, []entity.Vector{vector}, contentVector,
		entity.L2, 10, sp)
	if err != nil {
		fmt.Printf("fail to search collection: %s", err.Error())
		return resList, err
	}

	for _, result := range sr {
		var contentColumn *entity.ColumnVarChar
		var answerColumn *entity.ColumnVarChar
		var qaColumn *entity.ColumnInt32
		var imageUrls *entity.ColumnVarChar
		var videoUrls *entity.ColumnVarChar
		var relatedQuestions *entity.ColumnVarChar
		for _, field := range result.Fields {
			if field.Name() == contentCol {
				contentColumn = field.(*entity.ColumnVarChar)
			}
			if field.Name() == answerCol {
				answerColumn = field.(*entity.ColumnVarChar)
			}
			if field.Name() == qaCol {
				qaColumn = field.(*entity.ColumnInt32)
			}
			if field.Name() == imageUrlsCol {
				imageUrls = field.(*entity.ColumnVarChar)
			}
			if field.Name() == videoUrlsCol {
				videoUrls = field.(*entity.ColumnVarChar)
			}
			if field.Name() == relatedQuestionsCol {
				relatedQuestions = field.(*entity.ColumnVarChar)
			}
		}
		for i := 0; i < result.ResultCount; i++ {
			var know utils.QaRes
			content, err0 := contentColumn.ValueByIdx(i)
			if err0 == nil {
				know.Content = content
			}
			image, err1 := imageUrls.ValueByIdx(i)
			if err1 == nil {
				know.ImageUrls = image
			}
			video, err2 := videoUrls.ValueByIdx(i)
			if err2 == nil {
				know.VideoUrls = video
			}
			relatedQ, err3 := relatedQuestions.ValueByIdx(i)
			if err3 == nil {
				know.RelatedQuestions = relatedQ
			}
			qa, err4 := qaColumn.ValueByIdx(i)
			if err4 == nil {
				know.Qa = int(qa)
			}
			answer, err5 := answerColumn.ValueByIdx(i)
			if err5 == nil {
				know.Answer = answer
			}
			know.Scores = result.Scores[i]
			resList = append(resList, know)

		}
	}
	return resList, nil
}

func QueryQaListKeyWords(req utils.QaReq) ([]utils.QaRes, error) {
	resList := make([]utils.QaRes, 0)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	vectors, err := myHttp.EmbeddingThunderSoft(req.KeyWords)
	if err != nil {
		fmt.Printf("embeding %s error: %w", req.KeyWords, err)
		return resList, err
	}
	vector := entity.FloatVector(vectors[0])
	expr := ""
	if len(req.DatasetIds) > 0 {
		joinedString := strings.Join(req.DatasetIds, ",")
		expr = fmt.Sprintf("datasetId in [%s] and qa == 1", joinedString)
	} else {
		expr = fmt.Sprintf(" qa == 1")
	}
	sp, _ := entity.NewIndexFlatSearchParam()
	sr, err := milvusClient.Search(ctx, conf.ConfigInfo.Milvus.CollectionName, []string{}, expr, []string{contentCol, imageUrlsCol, videoUrlsCol, relatedQuestionsCol}, []entity.Vector{vector}, contentVector,
		entity.L2, req.Top, sp)
	if err != nil {
		fmt.Printf("fail to search collection: %s", err.Error())
		return resList, err
	}

	for _, result := range sr {
		var contentColumn *entity.ColumnVarChar
		var imageUrls *entity.ColumnVarChar
		var videoUrls *entity.ColumnVarChar
		var relatedQuestions *entity.ColumnVarChar
		for _, field := range result.Fields {
			if field.Name() == contentCol {
				contentColumn = field.(*entity.ColumnVarChar)
			}
			if field.Name() == imageUrlsCol {
				imageUrls = field.(*entity.ColumnVarChar)
			}
			if field.Name() == videoUrlsCol {
				videoUrls = field.(*entity.ColumnVarChar)
			}
			if field.Name() == relatedQuestionsCol {
				relatedQuestions = field.(*entity.ColumnVarChar)
			}
		}
		for i := 0; i < result.ResultCount; i++ {
			var know utils.QaRes
			know.Qa = 1
			content, err0 := contentColumn.ValueByIdx(i)
			if err0 == nil {
				know.Content = content
			}

			image, err1 := imageUrls.ValueByIdx(i)
			if err1 == nil {
				know.ImageUrls = image
			}
			video, err2 := videoUrls.ValueByIdx(i)
			if err2 == nil {
				know.VideoUrls = video
			}
			relatedQ, err3 := relatedQuestions.ValueByIdx(i)
			if err3 == nil {
				know.RelatedQuestions = relatedQ
			}
			resList = append(resList, know)

		}
	}
	return resList, nil
}

func QueryQaListAll(req utils.QaReq) ([]utils.QaRes, error) {
	resList := make([]utils.QaRes, 0)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	expr := ""
	if len(req.DatasetIds) > 0 {
		joinedString := strings.Join(req.DatasetIds, ",")
		expr = fmt.Sprintf("datasetId in [%s] and qa == 1", joinedString)
	} else {
		expr = fmt.Sprintf(" qa == 1")
	}

	result, err := milvusClient.Query(ctx, conf.ConfigInfo.Milvus.CollectionName,
		[]string{}, expr, []string{contentCol, imageUrlsCol, videoUrlsCol, relatedQuestionsCol},
	)
	if err != nil {
		fmt.Printf("fail to search collection: %s", err.Error())
		return resList, err
	}

	fmt.Println("sr=", result)
	var contentColumn *entity.ColumnVarChar
	var imageUrls *entity.ColumnVarChar
	var videoUrls *entity.ColumnVarChar
	var relatedQuestions *entity.ColumnVarChar
	for _, field := range result {
		if field.Name() == contentCol {
			contentColumn = field.(*entity.ColumnVarChar)
		}
		if field.Name() == imageUrlsCol {
			imageUrls = field.(*entity.ColumnVarChar)
		}
		if field.Name() == videoUrlsCol {
			videoUrls = field.(*entity.ColumnVarChar)
		}
		if field.Name() == relatedQuestionsCol {
			relatedQuestions = field.(*entity.ColumnVarChar)
		}
	}
	for i := 0; i < contentColumn.Len(); i++ {
		if i >= req.Top {
			break
		}
		var know utils.QaRes
		know.Qa = 1
		content, err0 := contentColumn.ValueByIdx(i)
		if err0 == nil {
			know.Content = content
		}
		image, err1 := imageUrls.ValueByIdx(i)
		if err1 == nil {
			know.ImageUrls = image
		}
		video, err2 := videoUrls.ValueByIdx(i)
		if err2 == nil {
			know.VideoUrls = video
		}
		relatedQ, err3 := relatedQuestions.ValueByIdx(i)
		if err3 == nil {
			know.RelatedQuestions = relatedQ
		}

		resList = append(resList, know)
	}
	return resList, nil
}
