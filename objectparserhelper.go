package mongodb

import "go.mongodb.org/mongo-driver/bson"

func MarshalBsonDocument(raw bson.Raw, docType interface{}) (interface{}, error) {
	err := bson.Unmarshal(raw, &docType)
	return docType, err
}
