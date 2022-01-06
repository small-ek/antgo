package mgo

import (
	"context"
	"errors"
	"github.com/small-ek/antgo/os/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

//Mgo
type Mgo struct {
	Database   *mongo.Database //数据库
	Client     *mongo.Client   //默认链接
	Ctx        context.Context
	Collection *mongo.Collection
	Cancel     context.CancelFunc
	Filter     interface{}
	Pages      Pages
}

//Pages 分页过滤排序
type Pages struct {
	Limit *int64
	Skip  *int64
	Sort  interface{}
}

//GetConfig 获取配置
func GetConfig() (string, string, uint64, int) {
	uri := config.Decode().Get("mgo.uri").String()
	database := config.Decode().Get("mgo.database").String()
	poollimit := config.Decode().Get("mgo.poollimit").Uint64()
	timeout := config.Decode().Get("mgo.timeout").Int()
	return uri, database, poollimit, timeout
}

//Connect Default connection<默认连接>
//参考文档 https://docs.mongodb.com/drivers/go/
func Connect(databaseName ...string) *Mgo {
	uri, db, poollimit, timeout := GetConfig()
	if timeout == 0 {
		timeout = 120
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	opts := options.Client()
	opts.ApplyURI(uri)

	opts.SetMaxPoolSize(poollimit)

	Client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	if err := Client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	if len(databaseName) > 0 {
		db = databaseName[0]
	}
	database := &mongo.Database{}
	if db != "" {
		database = Client.Database(db)
	}
	return &Mgo{
		Client:   Client,
		Cancel:   cancel,
		Database: database,
		Ctx:      ctx,
	}
}

//SetDatabase Modify database switch<修改数据库切换>
func (m *Mgo) SetDatabase(databaseName string) *Mgo {
	m.Database = m.Client.Database(databaseName)
	return m
}

//Table Setting table<表>
func (m *Mgo) Table(tableName string) *Mgo {
	if m.Database.Name() != "" {
		m.Collection = m.Database.Collection(tableName)
	}
	return m
}

//Create Create data<创建数据>
func (m *Mgo) Create(data interface{}) (*mongo.InsertOneResult, error) {

	return m.Collection.InsertOne(m.Ctx, data)
}

//SaveAll save all data<创建多条数据>
func (m *Mgo) SaveAll(data []interface{}) (*mongo.InsertManyResult, error) {

	opts := options.InsertMany().SetOrdered(false)
	return m.Collection.InsertMany(m.Ctx, data, opts)
}

//Update Update data<修改数据>
func (m *Mgo) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {

	return m.Collection.UpdateMany(
		m.Ctx,
		filter,
		bson.M{
			"$set": update,
		},
	)
}

//UpdateById Modify data according to id<根据id修改>
func (m *Mgo) UpdateById(id string, update interface{}) (*mongo.UpdateResult, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return m.Collection.UpdateOne(
		m.Ctx,
		bson.M{"_id": oid},
		bson.M{
			"$set": update,
		},
	)
}

//Delete delete data<删除>
func (m *Mgo) Delete(update interface{}) (*mongo.DeleteResult, error) {

	return m.Collection.DeleteOne(
		m.Ctx,
		update,
	)
}

//DeleteMany Delete multiple data<删除多个数据>
func (m *Mgo) DeleteMany(update interface{}) (*mongo.DeleteResult, error) {

	return m.Collection.DeleteMany(
		m.Ctx,
		update,
	)
}

//Count Get the total quantity<获取总数量>
func (m *Mgo) Count() (int64, error) {
	if m.Filter == nil {
		m.Filter = bson.D{}
	}
	return m.Collection.CountDocuments(m.Ctx, m.Filter)
}

//FindOne Single query<单个查询>
func (m *Mgo) FindOne() (bson.M, error) {
	var result bson.M
	if m.Filter == nil {
		m.Filter = bson.D{}
	}
	if m.Collection != nil {
		if err := m.Collection.FindOne(m.Ctx, m.Filter).Decode(&result); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("Database connection failed")
}

//Find Multiple data search<多条数据查询>
func (m *Mgo) Find() (*mongo.Cursor, error) {
	if m.Filter == nil {
		m.Filter = bson.D{}
	}

	if m.Collection != nil {
		return m.Collection.Find(m.Ctx, m.Filter, &options.FindOptions{Limit: m.Pages.Limit, Skip: m.Pages.Skip, Sort: m.Pages.Sort})
	}
	return nil, errors.New("Database connection failed")

}

//Distinct Query unique data<查询不重复的数据>
func (m *Mgo) Distinct(name string) ([]interface{}, error) {
	if m.Filter == nil {
		m.Filter = bson.D{}
	}

	if m.Collection != nil {
		return m.Collection.Distinct(m.Ctx, name, m.Filter, &options.DistinctOptions{})

	}
	return nil, errors.New("Database connection failed")

}

//Limit Limited number<显示数量>
func (m *Mgo) Limit(limit int64) *Mgo {
	m.Pages.Limit = &limit
	return m
}

//Skip How many pages to jump<跳转多少页>
func (m *Mgo) Skip(skip int64) *Mgo {
	m.Pages.Skip = &skip
	return m
}

//Sort Sort data<排序>
//{KEY:1},{KEY:-1}
func (m *Mgo) Sort(sort map[string]interface{}) *Mgo {
	var bsonM bson.M
	bsonM = sort
	m.Pages.Sort = bsonM
	return m
}

//Where 条件 [][]string{{"test", "name"},{"test", "$lt","1"}}
//map[string]interface{}{
//		"$or": []interface{}{map[string]interface{}{"author": "Nicolas222"}, map[string]interface{}{"author": "Nicolas333"}},
//	}
//等于	{<key>:<value>}
//小于	{<key>:{$lt:<value>}}
//小于或等于	{<key>:{$lte:<value>}}
//大于	{<key>:{$gt:<value>}}
//大于或等于	{<key>:{$gte:<value>}}
//不等于	{<key>:{$ne:<value>}}
func (m *Mgo) Where(filter interface{}) *Mgo {
	where := bson.D{}
	switch filters := filter.(type) {
	case [][]string:
		for i := 0; i < len(filters); i++ {
			value := filters[i]
			if len(value) == 2 {
				where = append(where, bson.E{value[0], value[1]})
			}
			if len(value) == 3 {
				where = append(where, bson.E{value[0], bson.D{{value[1], value[2]}}})
			}
		}
	}

	if len(where) > 0 {
		m.Filter = &where
	} else {
		m.Filter = &filter
	}
	return m
}

//Close Close the connection<关闭连接>
func (m *Mgo) Close() {
	defer m.Cancel()
	defer func() {
		if err := m.Client.Disconnect(m.Ctx); err != nil {
			panic(err)
		}
	}()
}

//Aggregate Aggregate advanced queries<聚合高级查询>
//id, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")
//matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
//groupStage := bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}}
//mongo.Pipeline{matchStage, groupStage}
func (m *Mgo) Aggregate() (*mongo.Cursor, error) {

	if m.Collection != nil {
		return m.Collection.Aggregate(m.Ctx, m.Filter)
	}
	return nil, errors.New("Database connection failed")
}

//GetCtx Get context<获取上下文>
func (m *Mgo) GetCtx() context.Context {
	return m.Ctx
}

//GetClient Get Client<获取连接>
func (m *Mgo) GetClient() *mongo.Client {
	return m.Client
}

//StartTrans 开启事务
//sessionContext.StartTransaction() 开启事务
//sessionContext.AbortTransaction(sessionContext) //终止事务
//sessionContext.CommitTransaction(sessionContext) //提交事务
func (m *Mgo) StartTrans(fn func(mongo.SessionContext) error) {
	m.Client.UseSession(m.Ctx, fn)
}
