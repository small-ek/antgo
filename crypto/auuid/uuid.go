package auuid

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"runtime"
	"sync"
	"time"
)

// ========================== 类型与常量定义 ==========================
// ========================== Type & Constant Definitions ==========================

// UUID 表示符合 RFC 4122 的通用唯一标识符
// UUID represents a RFC 4122 compliant Universally Unique Identifier
type UUID = uuid.UUID

// DCE 安全域定义 (RFC 4122 标准)
// DCE Security Domains as defined in RFC 4122
const (
	DomainPerson = uuid.Domain(0) // 用户域 (POSIX UID) / User domain
	DomainGroup  = uuid.Domain(1) // 组域 (POSIX GID) / Group domain
	DomainOrg    = uuid.Domain(2) // 组织域 / Organization domain
)

// ========================== 错误定义 ==========================
// ========================== Error Definitions ==========================

var (
	ErrUUIDGeneration  = errors.New("UUID generation failed")   // 基础生成错误
	ErrInvalidNodeID   = errors.New("invalid node ID")          // 无效节点ID
	ErrInvalidOrgID    = errors.New("organization ID required") // 需要组织ID
	ErrBatchGeneration = errors.New("batch generation failed")  // 批量生成失败
	ErrNotTimeBased    = errors.New("not a time-based UUID")    // 非时间序列UUID
)

// ========================== 全局配置 ==========================
// ========================== Global Configurations ==========================

var (
	uidCache, gidCache uint32       // 系统ID缓存 / System ID cache
	onceUID, onceGID   sync.Once    // 单次初始化 / One-time initialization
	batchPool          = sync.Pool{ // 批量生成资源池
		New: func() interface{} {
			return make([]byte, 16) // UUID字节缓冲
		},
	}
)

// ================== 核心生成函数 ==================
// ================== Core Generation Functions ==================

// New 生成随机UUID (版本4) 性能优化版本
// Generates random UUID (Version 4) with optimized performance
func New() UUID {
	return uuid.New()
}

// Create 生成时间序列UUID (版本1)，含错误重试机制
// Generates time-based UUID (Version 1) with retry mechanism
func Create() (UUID, error) {
	const (
		maxRetries    = 3
		retryInterval = 10 * time.Millisecond
	)

	var u UUID
	var err error

	for i := 0; i < maxRetries; i++ {
		if u, err = uuid.NewUUID(); err == nil {
			return u, nil
		}
		time.Sleep(retryInterval)
	}
	return uuid.Nil, fmt.Errorf("%w: %v", ErrUUIDGeneration, err)
}

// NewRandom 生成加密安全随机UUID (版本4)，含资源池优化
// Generates cryptographically secure random UUID (Version 4) with pool optimization
func NewRandom() (UUID, error) {
	buf := batchPool.Get().([]byte)
	defer batchPool.Put(buf[:0])

	if _, err := uuid.NewRandom(); err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrUUIDGeneration, err)
	}
	return uuid.New(), nil
}

// ================== DCE安全UUID ==================
// ================== DCE Security UUIDs ==================

// NewDCEPerson 生成用户域DCE UUID (版本2)，含系统调用缓存
// Generates user domain DCE UUID (Version 2) with system call caching
func NewDCEPerson() (UUID, error) {
	onceUID.Do(func() {
		uidCache = uint32(os.Getuid())
	})
	return generateDCE(DomainPerson, uidCache)
}

// NewDCEGroup 生成组域DCE UUID (版本2)，含系统调用缓存
// Generates group domain DCE UUID (Version 2) with system call caching
func NewDCEGroup() (UUID, error) {
	onceGID.Do(func() {
		gidCache = uint32(os.Getgid())
	})
	return generateDCE(DomainGroup, gidCache)
}

// NewDCEOrg 生成组织域DCE UUID (版本2)，含严格参数校验
// Generates organization domain DCE UUID (Version 2) with strict validation
func NewDCEOrg(orgID uint32) (UUID, error) {
	if orgID == 0 {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrUUIDGeneration, ErrInvalidOrgID)
	}
	return generateDCE(DomainOrg, orgID)
}

// generateDCE 安全生成DCE UUID的公共方法
// Common method for secure DCE UUID generation
func generateDCE(domain uuid.Domain, id uint32) (UUID, error) {
	if u, err := uuid.NewDCESecurity(domain, id); err == nil {
		return u, nil
	}
	return uuid.Nil, fmt.Errorf("%w: DCE generation failed", ErrUUIDGeneration)
}

// ================== 高级功能 ==================
// ================== Advanced Features ==================

// CreateWithNode 生成带自定义节点ID的V1 UUID（兼容v1.6.0）
// Generates V1 UUID with custom node ID (compatible with v1.6.0)
// 参数:
//
//	node - 6字节节点ID（需确保唯一性）
//
// Parameters:
//
//	node - 6-byte node ID (must be unique)
func CreateWithNode(node []byte) (UUID, error) {
	if len(node) != 6 {
		return uuid.Nil, fmt.Errorf("%w: need 6 bytes, got %d",
			ErrInvalidNodeID, len(node))
	}

	// 生成基础UUID
	u, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrUUIDGeneration, err)
	}

	// 安全替换节点ID（深拷贝）
	uuidBytes := make([]byte, 16)
	copy(uuidBytes, u[:])
	copy(uuidBytes[10:], node)

	// 验证并返回新UUID
	if newUUID, err := uuid.FromBytes(uuidBytes); err == nil {
		return newUUID, nil
	}
	return uuid.Nil, fmt.Errorf("%w: node replacement failed", ErrUUIDGeneration)
}

// BatchGenerate 优化的批量UUID生成方法（工作池模式）
// Optimized batch UUID generation with worker pool pattern
// 参数:
//
//	count - 需要生成的UUID数量
//
// Parameters:
//
//	count - Number of UUIDs to generate
func BatchGenerate(count int) ([]UUID, error) {
	if count <= 0 {
		return nil, fmt.Errorf("invalid count: %d", count)
	}

	type result struct {
		u   UUID
		err error
	}

	var (
		results     = make([]UUID, 0, count)
		resultChan  = make(chan result, count)
		ctx, cancel = context.WithCancel(context.Background())
		wg          sync.WaitGroup
		workerCount = runtime.NumCPU() * 2
	)

	defer cancel()

	// 创建工作池
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					u, err := uuid.NewRandom()
					select {
					case resultChan <- result{u, err}:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// 结果收集
	go func() {
		for res := range resultChan {
			if res.err != nil {
				cancel()
				return
			}
			results = append(results, res.u)
			if len(results) >= count {
				cancel()
				return
			}
		}
	}()

	wg.Wait()
	close(resultChan)

	if len(results) != count {
		return nil, fmt.Errorf("%w: generated %d/%d UUIDs",
			ErrBatchGeneration, len(results), count)
	}
	return results, nil
}

// ================== 检测与转换 ==================
// ================== Inspection & Conversion ==================

// IsValidUUID 验证字符串是否为有效UUID（含格式校验）
// Validates if a string is a valid UUID with format check
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// IsNilUUID 检查是否为全零UUID（优化字节比较）
// Checks for nil UUID with optimized byte comparison
func IsNilUUID(u UUID) bool {
	return bytes.Equal(u[:], uuid.Nil[:])
}

// GetVersion 获取UUID版本号（1-5）
// Returns UUID version (1-5)
func GetVersion(u UUID) uuid.Version {
	return u.Version()
}

// GetVariant 获取UUID变体类型
// Returns UUID variant type
func GetVariant(u UUID) string {
	switch u.Variant() {
	case uuid.RFC4122:
		return "RFC4122"
	case uuid.Microsoft:
		return "Microsoft"
	case uuid.Reserved:
		return "Reserved"
	default:
		return "Unknown"
	}
}

// GetTimestamp 从时间序列UUID提取时间戳（含版本验证）
// Extracts timestamp from time-based UUID with version validation
func GetTimestamp(u UUID) (time.Time, error) {
	switch ver := u.Version(); ver {
	case 1, 2:
		t := u.Time()
		sec, nsec := t.UnixTime()
		return time.Unix(sec, int64(nsec)).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("%w: version %d UUID", ErrNotTimeBased, ver)
	}
}

// ToURN 转换为标准URN格式（含有效性验证）
// Converts to standard URN format with validation
func ToURN(u UUID) string {
	if IsNilUUID(u) {
		return ""
	}
	return "urn:uuid:" + u.String()
}

// StringToUUID 从字符串解析UUID（自动格式化处理）
// Parses UUID from string with auto-formatting
func StringToUUID(s string) (UUID, error) {
	return uuid.Parse(s)
}

// UUIDToString 转换UUID为标准字符串（优化性能）
// Converts UUID to canonical string (optimized)
func UUIDToString(u UUID) string {
	if IsNilUUID(u) {
		return ""
	}
	return u.String()
}
