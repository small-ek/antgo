# str - 字符串处理工具库 / String Utilities Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`str` 包提供了一系列字符串处理的工具函数，适用于 Go 项目中的各种字符串操作需求。该库包含了去除引号、首字母大写、字符替换、符号过滤、字符串拆分与修剪、键格式化、转义字符剥离以及数字判断等常用功能，帮助开发者提高代码的可读性和开发效率。

GitHub 地址: [github.com/small-ek/antgo/utils/str](https://github.com/small-ek/antgo/utils/str)

### 📦 安装

```bash
go get github.com/small-ek/antgo/utils/str
```

### 🚀 快速开始

下面的示例展示了如何使用 `str` 包中的各个函数。

#### 1. 移除字符串首尾双引号 —— `ClearQuotes`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	input := `"Hello, World!"`
	output := str.ClearQuotes(input)
	fmt.Println(output) // 输出: Hello, World!
}
```

#### 2. 首字母大写 —— `UcFirst`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	s := "hello, 世界"
	result := str.UcFirst(s)
	fmt.Println(result) // 输出: Hello, 世界
}
```

#### 3. 根据映射替换字符串 —— `ReplaceByMap`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	original := "The quick brown fox jumps over the lazy dog."
	replacements := map[string]string{
		"quick": "slow",
		"brown": "red",
		"dog":   "cat",
	}
	result := str.ReplaceByMap(original, replacements)
	fmt.Println(result)
	// 输出: The slow red fox jumps over the lazy cat.
}
```

#### 4. 移除所有非字母和非数字的字符 —— `RemoveSymbols`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	input := "Go@Lang!2025"
	output := str.RemoveSymbols(input)
	fmt.Println(output) // 输出: GoLang2025
}
```

#### 5. 忽略指定符号后比较字符串 —— `EqualFoldWithoutChars`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	s1 := "Hello-World_2025"
	s2 := "hello world2025"
	equal := str.EqualFoldWithoutChars(s1, s2)
	fmt.Println(equal) // 输出: true
}
```

#### 6. 拆分字符串并修剪 —— `SplitAndTrim`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	text := " apple, banana ,  cherry , "
	// 使用默认的空白字符（以及 DefaultTrimChars 中定义的字符）进行修剪
	parts := str.SplitAndTrim(text, ",")
	for _, part := range parts {
		fmt.Println(part)
	}
	// 输出:
	// apple
	// banana
	// cherry
}
```

#### 7. 去除字符串两端空白或指定字符 —— `Trim`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	s := "  ***Hello, Go!***  "
	// 除了默认剥离字符外，再剥离 '*' 字符
	result := str.Trim(s, "*")
	fmt.Println(result) // 输出: Hello, Go!
}
```

#### 8. 格式化命令键 —— `FormatCmdKey`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	key := "COMMAND_KEY"
	formatted := str.FormatCmdKey(key)
	fmt.Println(formatted) // 输出: command.key
}
```

#### 9. 格式化环境变量键 —— `FormatEnvKey`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	key := "env.variable"
	formatted := str.FormatEnvKey(key)
	fmt.Println(formatted) // 输出: ENV_VARIABLE
}
```

#### 10. 移除转义反斜杠 —— `StripSlashes`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	escaped := "This is a \\\\test string\\\\ with escapes."
	result := str.StripSlashes(escaped)
	fmt.Println(result)
	// 输出: This is a test string with escapes.
}
```

#### 11. 判断字符串是否为数字 —— `IsNumeric`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	numbers := []string{"123", "-456", "78.90", "12a3", ""}
	for _, s := range numbers {
		fmt.Printf("IsNumeric(%q) = %v\n", s, str.IsNumeric(s))
	}
	// 输出:
	// IsNumeric("123") = true
	// IsNumeric("-456") = true
	// IsNumeric("78.90") = true
	// IsNumeric("12a3") = false
	// IsNumeric("") = false
}
```

### ✨ 核心特性

| 特性                             | 描述                                                                         |
|----------------------------------|------------------------------------------------------------------------------|
| **多功能字符串处理**             | 包含去除引号、首字母大写、字符替换、符号过滤、拆分修剪、键格式化、转义剥离以及数字判断等功能。 |
| **Unicode 支持**                 | 使用 `utf8` 和 `unicode` 包处理字符串，确保对 Unicode 字符的良好支持。         |
| **高性能**                       | 采用预分配内存和高效的字符串构建方式，保证在高并发场景下的性能。               |
| **易于使用**                     | 提供简洁的接口和示例，帮助开发者快速集成到项目中。                           |

### ⚠️ 注意事项

1. 在使用 `Trim` 和 `SplitAndTrim` 时，可通过额外的 `characterMask` 参数指定更多需要剥离的字符。
2. 使用 `ReplaceByMap` 进行字符替换时，替换的顺序不确定，适用于不依赖顺序的场景。
3. 在提交任务或字符串处理时，务必确保传入的字符串符合预期格式，以避免意外行为。

### 🤝 参与贡献

- [贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
- [提交 Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

The `str` package provides a collection of utility functions for string manipulation, designed for Go projects. It covers a wide range of common operations including removing quotes, capitalizing the first letter, replacing substrings via a map, filtering out symbols, splitting and trimming strings, formatting keys, stripping escape slashes, and checking if a string represents a numeric value. This library helps improve code readability and efficiency in handling string-related tasks.

GitHub URL: [github.com/small-ek/antgo/utils/str](https://github.com/small-ek/antgo/utils/str)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/utils/str
```

### 🚀 Quick Start

The following examples demonstrate how to use the various functions provided by the `str` package.

#### 1. Remove Leading and Trailing Double Quotes — `ClearQuotes`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	input := `"Hello, World!"`
	output := str.ClearQuotes(input)
	fmt.Println(output) // Output: Hello, World!
}
```

#### 2. Capitalize the First Letter — `UcFirst`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	s := "hello, 世界"
	result := str.UcFirst(s)
	fmt.Println(result) // Output: Hello, 世界
}
```

#### 3. Replace Substrings by Map — `ReplaceByMap`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	original := "The quick brown fox jumps over the lazy dog."
	replacements := map[string]string{
		"quick": "slow",
		"brown": "red",
		"dog":   "cat",
	}
	result := str.ReplaceByMap(original, replacements)
	fmt.Println(result)
	// Output: The slow red fox jumps over the lazy cat.
}
```

#### 4. Remove All Non-Letter and Non-Digit Characters — `RemoveSymbols`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	input := "Go@Lang!2025"
	output := str.RemoveSymbols(input)
	fmt.Println(output) // Output: GoLang2025
}
```

#### 5. Case-Insensitive Comparison After Removing Specific Symbols — `EqualFoldWithoutChars`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	s1 := "Hello-World_2025"
	s2 := "hello world2025"
	equal := str.EqualFoldWithoutChars(s1, s2)
	fmt.Println(equal) // Output: true
}
```

#### 6. Split String and Trim Each Part — `SplitAndTrim`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	text := " apple, banana ,  cherry , "
	// Splitting by comma and trimming each part using default trim characters.
	parts := str.SplitAndTrim(text, ",")
	for _, part := range parts {
		fmt.Println(part)
	}
	// Output:
	// apple
	// banana
	// cherry
}
```

#### 7. Trim Whitespace or Specified Characters from Both Ends — `Trim`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	s := "  ***Hello, Go!***  "
	// Trim default whitespace and additional '*' characters.
	result := str.Trim(s, "*")
	fmt.Println(result) // Output: Hello, Go!
}
```

#### 8. Format Command Key — `FormatCmdKey`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	key := "COMMAND_KEY"
	formatted := str.FormatCmdKey(key)
	fmt.Println(formatted) // Output: command.key
}
```

#### 9. Format Environment Variable Key — `FormatEnvKey`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	key := "env.variable"
	formatted := str.FormatEnvKey(key)
	fmt.Println(formatted) // Output: ENV_VARIABLE
}
```

#### 10. Strip Escape Backslashes — `StripSlashes`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	escaped := "This is a \\\\test string\\\\ with escapes."
	result := str.StripSlashes(escaped)
	fmt.Println(result)
	// Output: This is a test string with escapes.
}
```

#### 11. Check if a String is Numeric — `IsNumeric`

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/str"
)

func main() {
	numbers := []string{"123", "-456", "78.90", "12a3", ""}
	for _, s := range numbers {
		fmt.Printf("IsNumeric(%q) = %v\n", s, str.IsNumeric(s))
	}
	// Output:
	// IsNumeric("123") = true
	// IsNumeric("-456") = true
	// IsNumeric("78.90") = true
	// IsNumeric("12a3") = false
	// IsNumeric("") = false
}
```

### ✨ Key Features

| Feature                                  | Description                                                                |
|------------------------------------------|----------------------------------------------------------------------------|
| **Comprehensive String Operations**      | Functions for removing quotes, capitalizing letters, substring replacement, symbol removal, splitting and trimming, key formatting, escape stripping, and numeric checking. |
| **Unicode Compatibility**                | Utilizes `utf8` and `unicode` packages to support full Unicode character sets. |
| **High Performance**                     | Uses memory preallocation and efficient string building techniques to ensure high performance. |
| **Ease of Use**                          | Provides simple and clear interfaces along with example usage to integrate quickly into your projects. |

### ⚠️ Important Notes

1. When using `Trim` and `SplitAndTrim`, you can supply additional characters via the `characterMask` parameter to further customize trimming.
2. The `ReplaceByMap` function performs replacements in an unordered manner; ensure that the replacement logic does not depend on order.
3. Always validate input strings when processing to prevent unexpected behavior.

### 🤝 Contributing

- [Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
- [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)
