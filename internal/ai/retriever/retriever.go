package retriever

import (
	"SuperBizAgent/internal/ai/embedder"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"
	"strings"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	einoRetriever "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
)

func NewMilvusRetriever(ctx context.Context) (rtr einoRetriever.Retriever, err error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	r, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  common.MilvusCollectionName,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      1,
		Embedding: eb,
	})
	if err != nil {
		return nil, err
	}
	return &tolerantRetriever{inner: r}, nil
}

type tolerantRetriever struct {
	inner einoRetriever.Retriever
}

func (t *tolerantRetriever) Retrieve(ctx context.Context, query string, opts ...einoRetriever.Option) ([]*schema.Document, error) {
	docs, err := t.inner.Retrieve(ctx, query, opts...)
	if err == nil {
		return docs, nil
	}
	if shouldFallbackToEmptyDocs(err) {
		return []*schema.Document{}, nil
	}
	return nil, err
}

func shouldFallbackToEmptyDocs(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "extra output fields") &&
		strings.Contains(msg, "does not dynamic field")
}
