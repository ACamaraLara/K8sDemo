package mongodb

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

func processFilters(filters ...interface{}) bson.M {
	if len(filters) == 1 {
		// If a single filter is provided, convert it to bson.M if needed
		return convertToBSON(filters[0])
	} else if len(filters) > 1 {
		// Combine multiple filters using $and
		andFilters := make([]bson.M, len(filters))
		for i, f := range filters {
			andFilters[i] = convertToBSON(f)
		}
		return bson.M{"$and": andFilters}
	}

	// If no filter is provided, default to an empty filter
	return bson.M{}
}

func convertToBSON(input interface{}) bson.M {
	switch v := input.(type) {
	case bson.M:
		return v
	case map[string]interface{}:
		return bson.M(v)
	case struct{}:
		return convertStructToBSON(v)
	default:
		log.Printf("Unsupported filter type: %T", v)
		return bson.M{}
	}
}

func convertStructToBSON(input interface{}) bson.M {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	bsonDoc := bson.M{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		bsonDoc[field.Name] = value.Interface()
	}

	return bsonDoc
}
