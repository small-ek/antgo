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

// Mgo
type Mgo struct {
	Database     *mongo.Database //数据库
	Ctx          context.Context
	Collection   *mongo.Collection
	Cancel       context.CancelFunc
	DataBaseName string
	TableName    string
	Filter       interface{}
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

// InitEngine 初始化
// 参考文档 https://docs.mongodb.com/drivers/go/
func InitEngine(databaseName ...string) error {
	uri, _, poollimit, _ := GetConfig()
	ctx := context.Background()

	opts := options.Client()
	opts.ApplyURI(uri)

	opts.SetMaxPoolSize(poollimit)

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	client = c
	return err
}

// Connect Default connection<默认连接>
// 参考文档 https://docs.mongodb.com/drivers/go/
func NewMongoDB(tableName string, databaseName ...string) *Mgo {
	_, database, _, _ := GetConfig()
	var databaseNames = ""
	if len(databaseName) > 0 {
		databaseNames = databaseName[0]
	} else {
		databaseNames = database
	}

	return &Mgo{
		DataBaseName: databaseNames,
		TableName:    tableName,
		Ctx:          context.Background(),
	}
}

// SetDatabase Modify database switch<修改数据库切换>
func (m *Mgo) SetDatabase(databaseName string) *Mgo {
	m.DataBaseName = databaseName
	return m
}

// Table Setting table<表>
func (m *Mgo) Table(tableName string) *Mgo {

	m.TableName = tableName

	return m
}

// Create Create data<创建数据>
func (m *Mgo) Create(data interface{}) (*mongo.InsertOneResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).InsertOne(m.Ctx, data)
}

// SaveAll save all data<创建多条数据>
func (m *Mgo) SaveAll(data []interface{}) (*mongo.InsertManyResult, error) {
	opts := options.InsertMany().SetOrdered(false)
	return client.Database(m.DataBaseName).Collection(m.TableName).InsertMany(m.Ctx, data, opts)
}

// Update Update data<修改数据>
func (m *Mgo) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).UpdateMany(
		m.Ctx,
		filter,
		bson.M{
			"$set": update,
		},
	)
}

// UpdateById Modify data according to id<根据id修改>
func (m *Mgo) UpdateById(id string, update interface{}) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return client.Database(m.DataBaseName).Collection(m.TableName).UpdateOne(
		m.Ctx,
		bson.M{"_id": oid},
		bson.M{
			"$set": update,
		},
	)
}

// Delete delete data<删除>
func (m *Mgo) Delete(update interface{}) (*mongo.DeleteResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).DeleteOne(
		m.Ctx,
		update,
	)
}

// DeleteMany Delete multiple data<删除多个数据>
func (m *Mgo) DeleteMany(update interface{}) (*mongo.DeleteResult, error) {
	return client.Database(m.DataBaseName).Collection(m.TableName).DeleteMany(
		m.Ctx,
		update,
	)
}

// Count Get the total quantity<获取总数量>
func (m *Mgo) Count(count *int64) *Mgo {
	if m.Filter == nil {
		m.Filter = bson.D{}
	}
	total, err := client.Database(m.DataBaseName).Collection(m.TableName).CountDocuments(m.Ctx, m.Filter)
	if err != nil {
		fmt.Println(err)
		m.Err = err
	}

	*count = total
	return m
}

// FindOne Single query<单个查询>
func (m *Mgo) FindOne(v interface{}) error {
	m.Lock()
	defer m.Unlock()
	if m.Filter == nil {
		m.Filter = bson.D{}
	}
	if client != nil {
		if err := client.Database(m.DataBaseName).Collection(m.TableName).FindOne(m.Ctx, m.Filter).Decode(v); err != nil {
			return err
		}

		return nil
	}

	return errors.New("Database connection failed")
}

// Find Multiple data search<多条数据查询>
func (m *Mgo) Find() (*mongo.Cursor, error) {
	m.Lock()
	defer m.Unlock()
	if m.Filter == nil {
		m.Filter = bson.D{}
	}

	if client != nil {

		return client.Database(m.DataBaseName).Collection(m.TableName).Find(m.Ctx, m.Filter, &options.FindOptions{Limit: m.Pages.Limit, Skip: m.Pages.Skip, Sort: m.Pages.Sort})
	}
	return nil, errors.New("Database connection failed")

}

// Distinct Query unique data<查询不重复的数据>
func (m *Mgo) Distinct(name string) ([]interface{}, error) {
	m.Lock()
	defer m.Unlock()
	if m.Filter == nil {
		m.Filter = bson.D{}
	}

	if client != nil {
		return client.Database(m.DataBaseName).Collection(m.TableName).Distinct(m.Ctx, name, m.Filter, &options.DistinctOptions{})

	}
	return nil, errors.New("Database connection failed")

}

// Limit Limited number<显示数量>
func (m *Mgo) Limit(limit int64) *Mgo {
	m.Pages.Limit = &limit
	return m
}

// Skip How many pages to jump<跳转多少页>
func (m *Mgo) Skip(skip int64) *Mgo {
	m.Pages.Skip = &skip
	return m
}

// Sort Sort data<排序>
// {KEY:1},{KEY:-1}
func (m *Mgo) Sort(sort map[string]interface{}) *Mgo {
	var bsonM bson.M
	bsonM = sort
	m.Pages.Sort = bsonM
	return m
}

// Where 条件 [][]string{{"test", "name"},{"test", "$lt","1"}}
//
//	map[string]interface{}{
//			"$or": []interface{}{map[string]interface{}{"author": "Nicolas222"}, map[string]interface{}{"author": "Nicolas333"}},
//		}
//
// 等于	{<key>:<value>}
// 小于	{<key>:{$lt:<value>}}
// 小于或等于	{<key>:{$lte:<value>}}
// 大于	{<key>:{$gt:<value>}}
// 大于或等于	{<key>:{$gte:<value>}}
// 不等于	{<key>:{$ne:<value>}}
func (m *Mgo) Where(filter interface{}) *Mgo {
	where := bson.D{}
	switch filters := filter.(type) {
	case [][]interface{}:
		for i := 0; i < len(filters); i++ {
			value := filters[i]
			if len(value) == 2 {
				where = append(where, bson.E{conv.String(value[0]), value[1]})
			}
			if len(value) == 3 {
				where = append(where, bson.E{conv.String(value[0]), bson.D{{conv.String(value[1]), value[2]}}})
			}
		}
	case [][]string:
		for i := 0; i < len(filters); i++ {
			value := filters[i]
			if len(value) == 2 {
				where = append(where, bson.E{conv.String(value[0]), value[1]})
			}
			if len(value) == 3 {
				where = append(where, bson.E{conv.String(value[0]), bson.D{{conv.String(value[1]), value[2]}}})
			}
		}
	}

	if len(where) > 0 {
		m.Filter = &where
	} else {
		m.Filter = &bson.D{}
	}
	return m
}

// Close the connection<关闭连接>
func Close() {

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
}

// Aggregate Aggregate advanced queries<聚合高级查询>
// id, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")
// matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
// groupStage := bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}}
// mongo.Pipeline{matchStage, groupStage}
func (m *Mgo) Aggregate() (*mongo.Cursor, error) {
	if client != nil {
		return client.Database(m.DataBaseName).Collection(m.TableName).Aggregate(m.Ctx, m.Filter)
	}
	return nil, errors.New("Database connection failed")
}

// GetCtx Get context<获取上下文>
func (m *Mgo) GetCtx() context.Context {
	return m.Ctx
}

// StartTrans 开启事务
// sessionContext.StartTransaction() 开启事务
// sessionContext.AbortTransaction(sessionContext) //终止事务
// sessionContext.CommitTransaction(sessionContext) //提交事务
func (m *Mgo) StartTrans(fn func(mongo.SessionContext) error) {

	if err := client.UseSession(m.Ctx, fn); err != nil {
		panic(err)
	}
}

// BuildWhere 构造Where搜索
func BuildWhere(Filter []string) [][]interface{} {
	var where [][]interface{}
	for i := 0; i < len(Filter); i++ {
		var value = Filter[i]
		var binding []interface{}

		if err := json.Unmarshal(conv.Bytes(value), &binding); err != nil {
			panic(err)
		}
		if len(binding) == 3 && binding[1] == "=" {
			where = append(where, []interface{}{binding[0], binding[2]})
		} else if len(binding) == 3 && binding[1] != "=" {
			where = append(where, []interface{}{binding[0], condition(conv.String(binding[1])), binding[2]})
		}

	}
	return where
}

// condition 条件选择器
func condition(condition string) string {
	switch condition {
	case "<":
		return "$lt"
	case "=<", "<=":
		return "$lte"
	case ">":
		return "$gt"
	case ">=", "=>":
		return "$gte"
	case "!=", "<>":
		return "$ne"
	}
	return "="
}
