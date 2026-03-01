package retriever

import (
	"context"
	"errors"
	"testing"

	einoRetriever "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
)

func TestShouldFallbackToEmptyDocs(t *testing.T) {
	err := errors.New("[milvus retriever] search result has error: extra output fields [content metadata] found and result does not dynamic field")
	if !shouldFallbackToEmptyDocs(err) {
		t.Fatalf("expected fallback for milvus extra output fields error")
	}

	if shouldFallbackToEmptyDocs(errors.New("random error")) {
		t.Fatalf("unexpected fallback for unrelated error")
	}

	if shouldFallbackToEmptyDocs(nil) {
		t.Fatalf("unexpected fallback for nil error")
	}
}

func TestTolerantRetrieverReturnsEmptyDocsOnKnownMilvusError(t *testing.T) {
	fake := &fakeRetriever{
		err: errors.New("[milvus retriever] search result has error: extra output fields [content metadata] found and result does not dynamic field"),
	}
	r := &tolerantRetriever{inner: fake}
	docs, err := r.Retrieve(context.Background(), "test")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(docs) != 0 {
		t.Fatalf("expected empty docs, got %d", len(docs))
	}
}

func TestTolerantRetrieverPropagatesOtherErrors(t *testing.T) {
	expected := errors.New("boom")
	r := &tolerantRetriever{
		inner: &fakeRetriever{err: expected},
	}
	docs, err := r.Retrieve(context.Background(), "test")
	if err == nil || err.Error() != expected.Error() {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
	if docs != nil {
		t.Fatalf("expected nil docs when error is returned")
	}
}

type fakeRetriever struct {
	docs []*schema.Document
	err  error
}

func (f *fakeRetriever) Retrieve(ctx context.Context, query string, opts ...einoRetriever.Option) ([]*schema.Document, error) {
	return f.docs, f.err
}
