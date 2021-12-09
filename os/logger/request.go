package logger

//Request Middleware logging module(中间件记录日志模块)
//func Request() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		request := c.Request
//		//请求状态不是OPTIONS和日志开启
//		if request.Method != "OPTIONS" {
//			var body []byte
//			requestBody := make(map[string]interface{})
//
//			if c.Request.Body != nil {
//				body, _ = ioutil.ReadAll(c.Request.Body)
//				// 把刚刚读出来的再写进去其他地方使用没有
//				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
//				//解析body
//				if request.Header["Content-Type"] != nil && request.Header["Content-Type"][0] == "application/json" {
//					json.Unmarshal(body, &requestBody)
//				} else if len(body) > 0 {
//					bodyList := strings.Split(string(body), "&")
//					for i := 0; i < len(bodyList); i++ {
//						value := strings.Split(bodyList[i], "=")
//						requestBody[value[0]] = value[1]
//					}
//				}
//			}
//
//			// 开始时间
//			start := time.Now()
//			// 处理请求
//			c.Next()
//			// 结束时间
//			end := time.Now()
//			// 执行时间
//			run_time := end.Sub(start)
//			// 请求URL
//			path, _ := url.QueryUnescape(c.Request.URL.RequestURI())
//			// 请求IP
//			ip := c.ClientIP()
//			// 请求类型
//			method := c.Request.Method
//			// 请求状态
//			status := c.Writer.Status()
//
//			// 写入日志
//			Write.Info("RESTful Request",
//				zap.Any("ip", ip),
//				zap.Any("path", path),
//				zap.Any("request", requestBody),
//				zap.Any("run_time", run_time),
//				zap.Any("status", status),
//				zap.Any("method", method),
//				zap.Any("header", request.Header),
//			)
//		}
//	}
//}
