package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/salopensource/salversion/pkg/common"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SetupClient(ctx context.Context) (*firestore.Client, error) {
	projectID := common.GetEnv("PROJECT_ID", "notset")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func QueryDocuments(ctx context.Context, collection string, field string, operator string, value string) ([]*firestore.DocumentSnapshot, error) {
	var docs []*firestore.DocumentSnapshot
	client, err := SetupClient(ctx)
	if err != nil {
		return docs, err
	}

	defer client.Close()

	iter := client.Collection(collection).Where(field, operator, value).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return docs, err
		}
		docs = append(docs, doc)
		// fmt.Println(doc.Data())
	}

	return docs, nil
}

func SetDocument(ctx context.Context, collection string, document string, data map[string]interface{}) error {
	client, err := SetupClient(ctx)
	if err != nil {
		return err
	}

	defer client.Close()
	_, err = client.Collection(collection).Doc(document).Set(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDocument(ctx context.Context, collection string, document string, data map[string]interface{}) error {
	client, err := SetupClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Collection(collection).Doc(document).Set(ctx, data, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

func GetDocument(ctx context.Context, collection string, document string) (*firestore.DocumentSnapshot, bool, error) {
	client, err := SetupClient(ctx)
	if err != nil {
		return nil, false, err
	}
	defer client.Close()
	dsnap, err := client.Collection(collection).Doc(document).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return dsnap, true, nil
}

func DeleteDocument(ctx context.Context, collection string, document string) error {
	client, err := SetupClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Collection(collection).Doc(document).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
