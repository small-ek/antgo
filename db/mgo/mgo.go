package mgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/utils/conv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

var (
	client *mongo.Client
)

// Mgo MongoDB操作结构体
type Mgo struct {
	Database     *mongo.Database // 数据库
	Ctx          context.Context
	Collection   *mongo.Collection
	Cancel       context.CancelFunc
	DataBaseName string
	TableName    string
	Filter       bson.D
	Pages        Pages
	Timeout      int
	Err          error
	sync.Mutex
}

// Pages 分页过滤排序
type Pages struct {
	Limit *int64
	Skip  *int64
	Sort  interface{}
}

// GetConfig 获取配置
func GetConfig() (string, string, uint64, int) {
	uri := config.GetString("mgo_uri")
	database := config.GetString("mgo_database")
	poollimit := config.GetInt64("mgo_poollimit")
	timeout := config.GetInt("mgo_timeout")
	return uri, database, conv.Uint64(poollimit), timeout
}

// InitEngine 初始化MongoDB引擎
func InitEngine(databaseName ...string) error {
	uri, _, poollimit, _ := GetConfig()
	ctx := context.Background()

	opts := options.Client().ApplyURI(uri).SetMaxPoolSize(poollimit)

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	client = c
	return nil
}

// NewMongoDB 创建MongoDB实例
func NewMongoDB(tableName string, databaseName ...string) *Mgo {
	_, database, _, _ := GetConfig()
	dbName := database
	if len(databaseName) > 0 {
		dbName = databaseName[0]
	}

	return &Mgo{
		DataBaseName: dbName,
		TableName:    tableName,
		Ctx:          context.Background(),
		Filter:       bson.D{},
	}
}

// SetDatabase 设置数据库
func (m *Mgo) SetDatabase(databaseName string) *Mgo {
	m.DataBaseName = databaseName
	return m
}

// Table 设置表
func (m *Mgo) Table(tableName string) *Mgo {
	m.TableName = tableName
	return m
}

// Create 创建数据
func (m *Mgo) Create(data interface{}) (*mongo.InsertOneResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).InsertOne(m.Ctx, data)
}

// SaveAll 创建多条数据
func (m *Mgo) SaveAll(data []interface{}) (*mongo.InsertManyResult, error) {
	opts := options.InsertMany().SetOrdered(false)
	return client.Database(m.DataBaseName).Collection(m.TableName).InsertMany(m.Ctx, data, opts)
}

// Update 修改数据
func (m *Mgo) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).UpdateMany(
		m.Ctx,
		filter,
		bson.M{"$set": update},
	)
}

// UpdateById 根据id修改数据
func (m *Mgo) UpdateById(id string, update interface{}) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return client.Database(m.DataBaseName).Collection(m.TableName).UpdateOne(
		m.Ctx,
		bson.M{"_id": oid},
		bson.M{"$set": update},
	)
}

// Delete 删除数据
func (m *Mgo) Delete(filter interface{}) (*mongo.DeleteResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).DeleteOne(m.Ctx, filter)
}

// DeleteMany 删除多个数据
func (m *Mgo) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).DeleteMany(m.Ctx, filter)
}

// Count 获取总数量
func (m *Mgo) Count(count *int64) *Mgo {
	total, err := client.Database(m.DataBaseName).Collection(m.TableName).CountDocuments(m.Ctx, m.Filter)
	if err != nil {
		m.Err = err
	} else {
		*count = total
	}
	return m
}

// FindOne 单个查询
func (m *Mgo) FindOne(v interface{}) error {
	m.Lock()
	defer m.Unlock()

	if client == nil {
		return errors.New("database connection failed")
	}

	return client.Database(m.DataBaseName).Collection(m.TableName).FindOne(m.Ctx, m.Filter).Decode(v)
}

// Find 多条数据查询
func (m *Mgo) Find() (*mongo.Cursor, error) {
	m.Lock()
	defer m.Unlock()

	if client == nil {
		return nil, errors.New("database connection failed")
	}

	return client.Database(m.DataBaseName).Collection(m.TableName).Find(m.Ctx, m.Filter, &options.FindOptions{
		Limit: m.Pages.Limit,
		Skip:  m.Pages.Skip,
		Sort:  m.Pages.Sort,
	})
}

// Distinct 查询不重复的数据
func (m *Mgo) Distinct(name string) ([]interface{}, error) {
	m.Lock()
	defer m.Unlock()

	if client == nil {
		return nil, errors.New("database connection failed")
	}

	return client.Database(m.DataBaseName).Collection(m.TableName).Distinct(m.Ctx, name, m.Filter)
}

// Limit 显示数量
func (m *Mgo) Limit(limit int64) *Mgo {
	m.Pages.Limit = &limit
	return m
}

// Skip 跳转多少页
func (m *Mgo) Skip(skip int64) *Mgo {
	m.Pages.Skip = &skip
	return m
}

// Sort 排序
func (m *Mgo) Sort(sort map[string]interface{}) *Mgo {
	m.Pages.Sort = sort
	return m
}

// Where 条件查询
func (m *Mgo) Where(filter interface{}) *Mgo {
	var where bson.D
	switch filters := filter.(type) {
	case [][]interface{}:
		where = buildWhereFromInterfaceSlice(filters)
	case [][]string:
		where = buildWhereFromStringSlice(filters)
	default:
		where = bson.D{}
	}

	m.Filter = where
	return m
}

// buildWhereFromInterfaceSlice 从[][]interface{}构建查询条件
func buildWhereFromInterfaceSlice(filters [][]interface{}) bson.D {
	var where bson.D
	for _, value := range filters {
		if len(value) == 2 {
			where = append(where, bson.E{Key: conv.String(value[0]), Value: value[1]})
		} else if len(value) == 3 {
			where = append(where, bson.E{Key: conv.String(value[0]), Value: bson.D{{Key: conv.String(value[1]), Value: value[2]}}})
		}
	}
	return where
}

// buildWhereFromStringSlice 从[][]string构建查询条件
func buildWhereFromStringSlice(filters [][]string) bson.D {
	var where bson.D
	for _, value := range filters {
		if len(value) == 2 {
			where = append(where, bson.E{Key: value[0], Value: value[1]})
		} else if len(value) == 3 {
			where = append(where, bson.E{Key: value[0], Value: bson.D{{Key: value[1], Value: value[2]}}})
		}
	}
	return where
}

// Close 关闭连接
func Close() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}
}

// Aggregate 聚合高级查询
func (m *Mgo) Aggregate() (*mongo.Cursor, error) {
	if client == nil {
		return nil, errors.New("database connection failed")
	}
	return client.Database(m.DataBaseName).Collection(m.TableName).Aggregate(m.Ctx, m.Filter)
}

// GetCtx 获取上下文
func (m *Mgo) GetCtx() context.Context {
	return m.Ctx
}

// StartTrans 开启事务
func (m *Mgo) StartTrans(fn func(mongo.SessionContext) error) {
	if err := client.UseSession(m.Ctx, fn); err != nil {
		panic(err)
	}
}

// BuildWhere 构造Where搜索
func BuildWhere(Filter []string) [][]interface{} {
	var where [][]interface{}
	for _, value := range Filter {
		var binding []interface{}
		if err := json.Unmarshal([]byte(value), &binding); err != nil {
			panic(err)
		}
		if len(binding) == 3 {
			op := condition(conv.String(binding[1]))
			where = append(where, []interface{}{binding[0], op, binding[2]})
		}
	}
	return where
}

// condition 条件选择器
func condition(condition string) string {
	switch condition {
	case "<":
		return "$lt"
	case "<=", "=<":
		return "$lte"
	case ">":
		return "$gt"
	case ">=", "=>":
		return "$gte"
	case "!=", "<>":
		return "$ne"
	default:
		return "="
	}
}
