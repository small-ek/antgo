package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/examples/gin/model"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

type FilterCondition struct {
	And   []FilterCondition           `json:"and"`
	Or    []FilterCondition           `json:"or"`
	Not   map[string]interface{}      `json:"not"`
	Field string                      `json:"field"`
	Op    string                      `json:"op"`
	Value interface{}                 `json:"value"`
	Like  map[string]string           `json:"like"`
	Eq    map[string]interface{}      `json:"eq"`
	In    map[string][]interface{}    `json:"in"`
	Range map[string]map[string]int64 `json:"range"`
}

func buildQuery(db *gorm.DB, table string, queryParams map[string]interface{}) string {
	var queryStrings []string
	for key, value := range queryParams {
		switch key {
		case "and", "or":
			conditions := value.([]interface{})
			var subQueryStrings []string
			for _, condition := range conditions {
				subQuery := buildQuery(db, table, condition.(map[string]interface{}))
				subQueryStrings = append(subQueryStrings, fmt.Sprintf("(%s)", subQuery))
			}
			queryStrings = append(queryStrings, strings.Join(subQueryStrings, fmt.Sprintf(" %s ", key)))
		case "not":
			notQuery := buildQuery(db, table, value.(map[string]interface{}))
			queryStrings = append(queryStrings, fmt.Sprintf("NOT (%s)", notQuery))
		default:
			// Check if the value is a map, if so, it is a nested condition
			if nested, ok := value.(map[string]interface{}); ok {
				nestedQuery := buildQuery(db, table, nested)
				queryStrings = append(queryStrings, fmt.Sprintf("(%s = ? AND %s)", key, nestedQuery))
				db = db.Where(fmt.Sprintf("%s = ?", key), value)
			} else {
				queryStrings = append(queryStrings, fmt.Sprintf("%s = ?", key))
				db = db.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}
	return strings.Join(queryStrings, " ")
}
func main() {
	app := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	gin.DefaultWriter = ioutil.Discard

	app.GET("/", func(c *gin.Context) {
		jsonQuery := `{
		"and":[
			{"name": "John"}
		],
		"or": [
			{"name": "John"},
			{"and": [
				{"age": 25},
				{"or": [
					{"name": "Alice"},
					{"age": 10}
				]}
			]},
			{"name": "Bob"}
		],
		"not": {"status": "inactive"}
	}`
		// Unmarshal JSON string into PostgreSQLQueryCondition struct
		var queryParams map[string]interface{}
		if err := json.Unmarshal([]byte(jsonQuery), &queryParams); err != nil {
			panic(err)
		}

		var list []model.Admin
		ant.Db().Where(buildQuery(ant.Db(), "person", queryParams)).Find(&list)

		c.JSON(200, list)
	})
	app.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	//config := *flag.String("config", "./config.toml", "Configuration file path")
	//eng := ant.New().Etcd([]string{"127.0.0.1:2379"}, "/test.toml", "", "").Serve(app)
	configPath := *flag.String("config", "./config.toml", "Configuration file path")
	eng := ant.New(configPath).Serve(app)
	//result := model.Admin{}
	//page := page.PageParam{}
	////page.Filter=[]string{}{""}
	//ant.Db().Table("admin").Scopes(
	//	sql.Filters(page.Filter),
	//	sql.Order(page.Order),
	//).Find(&result)

	//alog.Info("result", zap.String("12", conv.String(result)))
	defer eng.Close()
}
