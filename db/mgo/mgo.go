package mgo

import (
	"context"
	"fmt"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/utils/conv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"
)

// MongoDB客户端（线程安全）
// Global MongoDB Client (Thread-Safe)
var (
	globalClient     *mongo.Client
	clientMutex      sync.RWMutex
	defaultTimeout   = 30 * time.Second
	maxPoolSizeLimit = 100
)

// MongoDB操作器（包含数据库操作配置）
// MongoDB Operator (Contains DB operation settings)
type MongoOperator struct {
	databaseName    string             // 数据库名称 / Database Name
	collectionName  string             // 集合名称 / Collection Name
	queryFilter     bson.D             // 查询过滤器 / Query Filter
	paginationOpts  PaginationOpt      // 分页配置 / Pagination Settings
	operationCtx    context.Context    // 操作上下文 / Operation Context
	operationCancel context.CancelFunc // 上下文取消函数 / Context Cancel Function
	lastError       error              // 最后操作错误 / Last Operation Error
	sync.RWMutex                       // 读写锁 / Read-Write Lock
}

// PaginationOpt 分页配置参数 / Pagination Options
type PaginationOpt struct {
	Limit int64       // 返回文档数量限制 / Limit the number of documents to return
	Skip  int64       // 跳过文档数量 / Skip the number of documents
	Sort  interface{} // 排序规则 / Sorting Rules
}

// initConfig 初始化配置 / Initialize Config
func initConfig() (uri string, dbName string, poolSize uint64, timeout time.Duration) {
	uri = config.GetString("mgo_uri")
	dbName = config.GetString("mgo_database")
	poolSize = conv.Uint64(config.GetInt64("mgo_poollimit"))
	timeout = time.Duration(config.GetInt("mgo_timeout")) * time.Second
	return
}

// InitConnection 初始化全局MongoDB连接（线程安全）
// Initialize Global MongoDB Connection (Thread-Safe)
func InitConnection(databaseNames ...string) error {
	uri, _, poolSize, _ := initConfig()

	// 连接配置 / Connection Options
	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(poolSize).
		SetServerSelectionTimeout(defaultTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	clientMutex.Lock()
	defer clientMutex.Unlock()

	// 创建新客户端 / Create New Client
	newClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("MongoDB connection failed: %w", err)
	}

	// 健康检查 / Health Check
	if err = newClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	globalClient = newClient
	return nil
}

// NewOperator 创建新的MongoDB操作实例 / Create New MongoDB Operator Instance
func NewOperator(collectionName string, databaseNames ...string) *MongoOperator {
	_, defaultDB, _, _ := initConfig()
	targetDB := defaultDB
	if len(databaseNames) > 0 {
		targetDB = databaseNames[0]
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &MongoOperator{
		databaseName:    targetDB,
		collectionName:  collectionName,
		operationCtx:    ctx,
		operationCancel: cancel,
		queryFilter:     bson.D{},
	}
}

// Clone 克隆操作器实例（用于链式调用） / Clone Operator Instance (For Chain Calling)
func (mo *MongoOperator) Clone() *MongoOperator {
	return &MongoOperator{
		databaseName:    mo.databaseName,
		collectionName:  mo.collectionName,
		queryFilter:     mo.queryFilter,
		paginationOpts:  mo.paginationOpts,
		operationCtx:    mo.operationCtx,
		operationCancel: mo.operationCancel,
	}
}

// SetDatabase 切换目标数据库 / Switch Target Database
func (mo *MongoOperator) SetDatabase(dbName string) *MongoOperator {
	mo.Lock()
	defer mo.Unlock()
	mo.databaseName = dbName
	return mo
}

// Collection 切换目标集合 / Switch Target Collection
func (mo *MongoOperator) Collection(colName string) *MongoOperator {
	mo.Lock()
	defer mo.Unlock()
	mo.collectionName = colName
	return mo
}

// InsertOne 插入单条文档 / Insert a Single Document
func (mo *MongoOperator) InsertOne(document interface{}) (*mongo.InsertOneResult, error) {
	client := mo.getClient()
	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		InsertOne(mo.operationCtx, document)
}

// BatchInsert 批量插入文档 / Insert Multiple Documents
func (mo *MongoOperator) BatchInsert(documents []interface{}) (*mongo.InsertManyResult, error) {
	client := mo.getClient()
	insertOpts := options.InsertMany().SetOrdered(false)
	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		InsertMany(mo.operationCtx, documents, insertOpts)
}

// Update 更新文档（支持批量更新）/ Update Documents (Supports Bulk Update)
func (mo *MongoOperator) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	client := mo.getClient()
	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		UpdateMany(mo.operationCtx, filter, bson.M{"$set": update})
}

// UpdateByID 通过ID更新文档 / Update Document by ID
func (mo *MongoOperator) UpdateByID(id string, update interface{}) (*mongo.UpdateResult, error) {
	client := mo.getClient()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID format: %w", err)
	}

	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		UpdateOne(mo.operationCtx, bson.M{"_id": oid}, bson.M{"$set": update})
}

// DeleteOne 删除单个文档 / Delete a Single Document
func (mo *MongoOperator) DeleteOne(filter interface{}) (*mongo.DeleteResult, error) {
	client := mo.getClient()
	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		DeleteOne(mo.operationCtx, filter)
}

// Delete 批量删除文档 / Delete Multiple Documents
func (mo *MongoOperator) Delete(filter interface{}) (*mongo.DeleteResult, error) {
	client := mo.getClient()
	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		DeleteMany(mo.operationCtx, filter)
}

// QueryCount 获取匹配文档数量 / Get Count of Matching Documents
func (mo *MongoOperator) QueryCount() (int64, error) {
	client := mo.getClient()
	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		CountDocuments(mo.operationCtx, mo.queryFilter)
}

// FindOne 查询单条文档 / Find One Document
func (mo *MongoOperator) FindOne(result interface{}) error {
	client := mo.getClient()
	mo.RLock()
	defer mo.RUnlock()

	err := client.Database(mo.databaseName).
		Collection(mo.collectionName).
		FindOne(mo.operationCtx, mo.queryFilter).
		Decode(result)
	if err != nil {
		return fmt.Errorf("document decode failed: %w", err)
	}
	return nil
}

// FindAll 查询多条文档（带分页）/ Find Multiple Documents (With Pagination)
func (mo *MongoOperator) FindAll() (*mongo.Cursor, error) {
	client := mo.getClient()
	mo.RLock()
	defer mo.RUnlock()

	findOpts := options.Find().
		SetLimit(mo.paginationOpts.Limit).
		SetSkip(mo.paginationOpts.Skip).
		SetSort(mo.paginationOpts.Sort)

	return client.Database(mo.databaseName).
		Collection(mo.collectionName).
		Find(mo.operationCtx, mo.queryFilter, findOpts)
}

// SetPagination 设置分页参数 / Set Pagination Parameters
func (mo *MongoOperator) SetPagination(page, pageSize int64) *MongoOperator {
	skip := page * pageSize
	mo.paginationOpts.Limit = pageSize
	mo.paginationOpts.Skip = skip
	return mo
}

// SetSorting 设置排序规则 / Set Sorting Rules
func (mo *MongoOperator) SetSorting(sortRules bson.D) *MongoOperator {
	mo.paginationOpts.Sort = sortRules
	return mo
}

// Where 构建查询条件 / Build Query Conditions
func (mo *MongoOperator) Where(conditions interface{}) *MongoOperator {
	var query bson.D
	switch cond := conditions.(type) {
	case [][]interface{}:
		query = buildDynamicQuery(cond)
	case map[string]interface{}:
		query = buildMapQuery(cond)
	case bson.D:
		query = cond
	default:
		query = bson.D{}
	}

	mo.queryFilter = query
	return mo
}

// ExecuteTransaction 执行事务操作 / Execute Transaction Operations
func (mo *MongoOperator) ExecuteTransaction(txnFunc func(mongo.SessionContext) error) error {
	client := mo.getClient()
	session, err := client.StartSession()
	if err != nil {
		return fmt.Errorf("session start failed: %w", err)
	}
	defer session.EndSession(mo.operationCtx)

	return mongo.WithSession(mo.operationCtx, session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("transaction start failed: %w", err)
		}

		if err := txnFunc(sessCtx); err != nil {
			_ = session.AbortTransaction(sessCtx)
			return err
		}

		if err := session.CommitTransaction(sessCtx); err != nil {
			return fmt.Errorf("transaction commit failed: %w", err)
		}
		return nil
	})
}

// CloseConnection 安全关闭全局连接 / Safely Close Global Connection
func CloseConnection() {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	if globalClient != nil {
		_ = globalClient.Disconnect(context.Background())
	}
}

// getClient 获取线程安全客户端 / Get Thread-Safe Client
func (mo *MongoOperator) getClient() *mongo.Client {
	clientMutex.RLock()
	defer clientMutex.RUnlock()
	return globalClient
}

// buildDynamicQuery 构建动态查询条件 / Build Dynamic Query Conditions
func buildDynamicQuery(conditions [][]interface{}) bson.D {
	var query bson.D
	for _, cond := range conditions {
		if len(cond) == 2 {
			query = append(query, bson.E{Key: conv.String(cond[0]), Value: cond[1]})
		} else if len(cond) == 3 {
			op := convertOperator(conv.String(cond[1]))
			query = append(query, bson.E{
				Key: conv.String(cond[0]),
				Value: bson.D{{
					Key:   op,
					Value: cond[2],
				}}})
		}
	}
	return query
}

// convertOperator 转换操作符 / Convert Operator
func convertOperator(op string) string {
	switch op {
	case ">":
		return "$gt"
	case ">=":
		return "$gte"
	case "<":
		return "$lt"
	case "<=":
		return "$lte"
	case "!=":
		return "$ne"
	default:
		return "$eq"
	}
}

// buildMapQuery 构建Map类型查询条件 / Build Map-Type Query Conditions
func buildMapQuery(conditions map[string]interface{}) bson.D {
	var query bson.D
	for k, v := range conditions {
		query = append(query, bson.E{Key: k, Value: v})
	}
	return query
}
