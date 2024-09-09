package test

import (
	"context"
	"fmt"
	"github.com/ChongYanOvO/little-blue-book/wire"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type MongoTestSuite struct {
	suite.Suite
	Server *gin.Engine
	mongo  *mongo.Database
}

func (s *MongoTestSuite) SetupSuite() {
	s.mongo, _ = wire.InitMongo()
}

type Student struct {
	Name string `bson:"name,omitempty"`
	Age  int    `bson:"age,omitempty"`
}

func (s *MongoTestSuite) TestCreate() {
	t := s.T()
	testCases := []struct {
		name string
	}{
		{"create"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			collection := s.mongo.Collection("test")
			//collection.InsertOne(ctx, Student{Name: "xixi", Age: 18})
			//collection.InsertOne(ctx, Student{Name: "haha", Age: 20})
			var st Student
			err := collection.FindOne(ctx, Student{Name: "xixi"}).Decode(&st)
			//assert.Error(t, err)
			fmt.Println(st)
			fmt.Println(err)
		})
	}
}

func TestMongo(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
