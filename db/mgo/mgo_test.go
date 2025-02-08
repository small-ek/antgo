package mgo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

// TestMongoDBOperations 测试 MongoDB 操作库的主要功能
func TestMongoDBOperations(t *testing.T) {
	// 初始化 MongoDB 连接
	err := InitConnection()
	assert.NoError(t, err, "Failed to initialize MongoDB connection")
	defer CloseConnection() // 确保测试结束后关闭连接

	// 定义测试集合名称
	collectionName := "test_users"
	dbName := "test_db"

	// 创建操作器实例
	operator := NewOperator(collectionName, dbName)

	// 测试数据
	user := bson.M{
		"name":  "Alice",
		"age":   25,
		"email": "alice@example.com",
	}

	// 1. 测试插入单条文档
	insertResult, err := operator.InsertOne(user)
	assert.NoError(t, err, "InsertOne failed")
	assert.NotNil(t, insertResult.InsertedID, "InsertedID should not be nil")

	// 2. 测试查询单条文档
	var foundUser bson.M
	err = operator.Where(bson.M{"name": "Alice"}).FindOne(&foundUser)
	assert.NoError(t, err, "FindOne failed")
	assert.Equal(t, "Alice", foundUser["name"], "Found user name should be Alice")
	assert.Equal(t, int32(25), foundUser["age"], "Found user age should be 25")

	// 3. 测试更新文档
	updateResult, err := operator.Update(bson.M{"name": "Alice"}, bson.M{"$set": bson.M{"age": 30}})
	assert.NoError(t, err, "Update failed")
	assert.Equal(t, int64(1), updateResult.ModifiedCount, "ModifiedCount should be 1")

	// 验证更新结果
	err = operator.Where(bson.M{"name": "Alice"}).FindOne(&foundUser)
	assert.NoError(t, err, "FindOne after update failed")
	assert.Equal(t, int32(30), foundUser["age"], "Updated user age should be 30")

	// 4. 测试删除文档
	deleteResult, err := operator.Delete(bson.M{"name": "Alice"})
	assert.NoError(t, err, "DeleteMany failed")
	assert.Equal(t, int64(1), deleteResult.DeletedCount, "DeletedCount should be 1")

	// 验证删除结果
	err = operator.Where(bson.M{"name": "Alice"}).FindOne(&foundUser)
	assert.Error(t, err, "FindOne after delete should fail")
}

// TestPagination 测试分页功能
func TestPagination(t *testing.T) {
	// 初始化 MongoDB 连接
	err := InitConnection()
	assert.NoError(t, err, "Failed to initialize MongoDB connection")
	defer CloseConnection()

	// 定义测试集合名称
	collectionName := "test_pagination"
	dbName := "test_db"

	// 创建操作器实例
	operator := NewOperator(collectionName, dbName)

	// 插入测试数据
	docs := []interface{}{
		bson.M{"name": "User1", "age": 20},
		bson.M{"name": "User2", "age": 25},
		bson.M{"name": "User3", "age": 30},
	}
	_, err = operator.BatchInsert(docs)
	assert.NoError(t, err, "BatchInsert failed")

	// 测试分页查询
	var results []bson.M
	cursor, err := operator.
		Where(bson.M{"age": bson.M{"$gte": 20}}).
		SetPagination(1, 2).            // 第一页，每页 2 条
		SetSorting(bson.D{{"age", 1}}). // 按 age 升序排序
		FindAll()
	assert.NoError(t, err, "FindAll failed")

	err = cursor.All(context.Background(), &results)
	assert.NoError(t, err, "Cursor decode failed")
	assert.Equal(t, 2, len(results), "Should return 2 documents")
	assert.Equal(t, "User1", results[0]["name"], "First document should be User1")
	assert.Equal(t, "User2", results[1]["name"], "Second document should be User2")
}

// TestTransactions 测试事务功能
func TestTransactions(t *testing.T) {
	// 初始化 MongoDB 连接
	err := InitConnection()
	assert.NoError(t, err, "Failed to initialize MongoDB connection")
	defer CloseConnection()

	// 定义测试集合名称
	collectionName := "test_transactions"
	dbName := "test_db"

	// 创建操作器实例
	operator := NewOperator(collectionName, dbName)

	// 测试事务
	err = operator.ExecuteTransaction(func(sessCtx mongo.SessionContext) error {
		// 插入一条文档
		_, err := operator.InsertOne(bson.M{"name": "Bob", "age": 40})
		if err != nil {
			return err
		}

		// 更新文档
		_, err = operator.Update(bson.M{"name": "Bob"}, bson.M{"$set": bson.M{"age": 45}})
		if err != nil {
			return err
		}

		return nil
	})
	assert.NoError(t, err, "Transaction failed")

	// 验证事务结果
	var result bson.M
	err = operator.Where(bson.M{"name": "Bob"}).FindOne(&result)
	assert.NoError(t, err, "FindOne failed")
	assert.Equal(t, int32(45), result["age"], "Updated age should be 45")
}
