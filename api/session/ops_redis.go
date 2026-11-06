package session

import (
	"context"
	"encoding/json"
	"fmt"
	"goVServer/api/dbops"
	"goVServer/api/defs"
	"goVServer/api/utils"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	ctx  = context.Background()
	once sync.Once
)

const (
	// 防缓存雪崩：基础 TTL 30min + 随机偏移 [0, 5min]
	baseTTL  = 30 * time.Minute
	maxJitter = 5 * time.Minute
	// 防缓存穿透：空值标记及短 TTL
	emptyMarker = "__NIL__"
	emptyTTL    = 2 * time.Minute
)

// InitRedis 初始化 Redis 连接（单例，多次调用安全）
func InitRedis(addr, password string, db int) error {
	var initErr error
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     20,
			MinIdleConns: 5,
		})
		initErr = rdb.Ping(ctx).Err()
		if initErr != nil {
			fmt.Printf("[Redis] 连接失败: %v\n", initErr)
		} else {
			fmt.Println("[Redis] 连接成功")
		}
	})
	return initErr
}

// randomTTL 生成随机 TTL，防止缓存雪崩
func randomTTL() time.Duration {
	return baseTTL + time.Duration(rand.Int63n(int64(maxJitter)))
}

func sessionKey(sid string) string {
	return "session:" + sid
}

// GenerateNewSessionId 生成新会话，写入 Redis + MySQL
func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := time.Now().UnixMilli()
	ttl := ct + 30*60*1000
	ss := &defs.SimpleSession{TTL: ttl, Username: un}

	// 1. 写入 Redis（主缓存）
	data, _ := json.Marshal(ss)
	if err := rdb.Set(ctx, sessionKey(id), data, randomTTL()).Err(); err != nil {
		fmt.Printf("[Session] Redis write error: %v\n", err)
	}

	// 2. 持久化到 MySQL（兜底）
	_ = dbops.InsertSession(id, ttl, un)

	return id
}

// IsSessionExpired 从 Redis 查询会话是否过期
// 查询路径：Redis → MySQL（回源）→ 缓存空值（防穿透）
func IsSessionExpired(sid string) (string, bool) {
	key := sessionKey(sid)

	// 1. 查 Redis
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// 缓存未命中，可能是空值（穿透防护）或真的不存在
		// 回源 MySQL
		return loadFromDB(sid)
	}
	if err != nil {
		fmt.Printf("[Session] Redis read error: %v\n", err)
		return loadFromDB(sid)
	}

	// 空值标记 → 缓存穿透防护，直接判定过期
	if val == emptyMarker {
		return "", true
	}

	ss := &defs.SimpleSession{}
	if err := json.Unmarshal([]byte(val), ss); err != nil {
		return "", true
	}

	// 检查 TTL
	if ss.TTL < time.Now().UnixMilli() {
		deleteExpiredSessions(sid)
		return "", true
	}

	// 滑动续期
	_ = rdb.Expire(ctx, key, randomTTL()).Err()

	return ss.Username, false
}

// loadFromDB 回源 MySQL，并回填 Redis
func loadFromDB(sid string) (string, bool) {
	dbSession, err := dbops.RetrieveSession(sid)
	if err != nil || dbSession == nil {
		// 数据库也没有 → 缓存空值，2 分钟内不再回源
		_ = rdb.Set(ctx, sessionKey(sid), emptyMarker, emptyTTL).Err()
		return "", true
	}

	if dbSession.TTL < time.Now().UnixMilli() {
		deleteExpiredSessions(sid)
		return "", true
	}

	// 回填 Redis
	data, _ := json.Marshal(dbSession)
	_ = rdb.Set(ctx, sessionKey(sid), data, randomTTL()).Err()

	return dbSession.Username, false
}

// deleteExpiredSessions 删除 Redis + MySQL 中的过期会话
func deleteExpiredSessions(sid string) {
	_ = rdb.Del(ctx, sessionKey(sid)).Err()
	_ = dbops.DeleteSession(sid)
}

// LoadSessionsFromDB 启动时从 MySQL 预热 Redis
func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		fmt.Printf("[Session] 预热失败: %v\n", err)
		return
	}
	count := 0
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sid, _ := k.(string)
		data, _ := json.Marshal(ss)
		_ = rdb.Set(ctx, sessionKey(sid), data, randomTTL()).Err()
		count++
		return true
	})
	fmt.Printf("[Session] 从 MySQL 预热 %d 条会话到 Redis\n", count)
}

// Close 关闭 Redis 连接
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
