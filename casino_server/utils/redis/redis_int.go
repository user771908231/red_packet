package data

import (
	"casino_server/common/cfg"
	"casino_server/common/log"
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/golang/protobuf/proto"
)

const (
	SERVERADDR = "127.0.0.1:6379"
)

var (
	Redis_svr string
)

func InitRedis() {
	config := cfg.Get()

	Redis_svr = SERVERADDR
	if len(config["redis_svr"]) != 0 {
		Redis_svr = config["redis_svr"]
	}

	log.T("redis_svr:%v", Redis_svr)
}

func TestRedis() error {
	db := new(redis.Conn)
	log.T("TestRedis db: %v", db)

	return nil
}

type Data struct {
	conn redis.Conn
}

func (t *Data) Open(table string) error {
	var err error
	t.conn, err = redis.DialTimeout("tcp", Redis_svr, 3*time.Second, 10*time.Second, 10*time.Second) //timeout=10s

	if err != nil {
		log.Error("ERR: redis Open(table:%v) ret err:%v t.conn:%v", table, err, t.conn)
		return err
	}

	//	if table != "" {
	//		_, err = t.conn.Do("SELECT", table)
	//		if err != nil {
	//			log.Error("[ERROR] redis.Select(%v) ret err:%v", table, err)
	//		}
	//		log.T("[TRACE] redis.Open(%v) ok...", table)
	//	}

	return err
}

func (t *Data) Close() (err error) {
	if t.conn != nil {
		log.T("[TRACE] t.conn.Close().")
		err = t.conn.Close()
		t.conn = nil
		return err
	} else {
		log.Fatal("FATAL ERR: try redis.Close() BUT t.conn=nil")
	}
	return err
}

func (t *Data) Select(table string) (err error) {
	//	_, err = t.conn.Do("SELECT", table)
	//	if err != nil {
	//		log.Error("[ERROR] redis.Select(%v) ret err:%v", table, err)
	//		return err
	//	}
	//	log.T("db.select(%v) ok.", table)
	//
	//	return err
	return nil
}

//================= String ==================

func (t *Data) GetKeys(keyword string) (keys []string, err error) {
	if t.conn != nil {
		if keyword == "" {
			keyword = "*"
		}

		keys, err = redis.Strings(t.conn.Do("KEYS", keyword))
		if err == redis.ErrNil {
			err = nil
		}
		return keys, err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}

	return keys, err
}

func (t *Data) Get(key string) (value string, err error) {
	if t.conn != nil {
		value, err := redis.String(t.conn.Do("GET", key))
		//log.T("redis.GET(%v) ret err:%v value:%v", key, err, value)
		if err == redis.ErrNil {
			err = nil
		}

		return value, err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}

	return "", err
}

//
func (t *Data) MGet(args []interface{}) (values []interface{}, err error) {
	if t.conn != nil {
		values, err := redis.Values(t.conn.Do("MGET", args...))
		if err == redis.ErrNil {
			err = nil
		}

		//log.T("redis.GET(%v) ret err:%v value:%v", args, err, values)
		return values, err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}

	return values, err
}

//return []byte
func (t *Data) Gets(key string) (value []byte, err error) {
	if t.conn != nil {
		value, err := redis.Bytes(t.conn.Do("GET", key))
		if err == redis.ErrNil {
			err = nil
		}

		return value, err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}

	return nil, err
}

func (t *Data) GetInt(key string) (value int, err error) {
	if t.conn != nil {
		value, err := redis.Int(t.conn.Do("GET", key))
		if err == redis.ErrNil {
			err = nil
		}

		return value, err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}

	return 0, err
}

func (t *Data) Set(key string, value []byte) error {
	if t.conn != nil {
		//log.T("[TRACE] try redis.Set(%v) value:%v", key, value)
		_, err := redis.String(t.conn.Do("SET", key, value))
		//log.T("[TRACE] after redis.Set(%v) ret err:%v", key, err)
		return err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}
	return nil
}

func (t *Data) SetInt(key string, value int32) error {
	if t.conn != nil {
		_, err := redis.String(t.conn.Do("SET", key, value))
		return err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}
	return nil
}

func (t *Data) SetUInt(key string, value uint32) error {
	if t.conn != nil {
		_, err := redis.String(t.conn.Do("SET", key, value))
		return err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}
	return nil
}

func (t *Data) SetObj(key string,pb proto.Message) error{
	d,err :=proto.Marshal(pb)
	if err != nil {
		return err
	}
	return t.Set(key,d)
}

func (t *Data) GetObj(key string,pb proto.Message) error{
	d,err := t.Gets(key)
	if err != nil {
		return err
	}
	return proto.Unmarshal(d,pb)
}

func (t *Data) Del(key string) error {
	if t.conn != nil {
		_, err := redis.String(t.conn.Do("DEL", key))
		return err
	} else {
		log.Fatal("invalid redis conn:%v", t.conn)
	}
	return nil
}

//================= List ==================
func (t *Data) ListGetAll(key string) (values []interface{}, err error) {

	values, err = redis.Values(t.conn.Do("LRANGE", key, 0, -1))
	if err == redis.ErrNil {
		err = nil
	}

	return values, err
}

func (t *Data) LIndex(key string, index int) (value []byte, err error) {

	value, err = redis.Bytes(t.conn.Do("LINDEX", key, index))
	if err == redis.ErrNil {
		err = nil
	}

	return value, err
}

func (t *Data) LRange(key string, start int, stop int) (values []interface{}, err error) {

	values, err = redis.Values(t.conn.Do("LRANGE", key, start, stop))
	if err == redis.ErrNil {
		err = nil
	}

	return values, err
}

func (t *Data) LLen(key string) (count int, err error) {

	count, err = redis.Int(t.conn.Do("LLEN", key))
	if err == redis.ErrNil {
		count = 0
		err = nil
	}

	return count, err
}

func (t *Data) LPush(key string, value []byte) (err error) {
	_, err = t.conn.Do("LPUSH", key, value)

	return err
}

func (t *Data) LPushInt(key string, value int32) (err error) {
	_, err = t.conn.Do("LPUSH", key, value)

	return err
}

func (t *Data) LRem(key string, removeValue []byte) (err error) {

	_, err = redis.Values(t.conn.Do("LREM", 0, removeValue))
	return err
}

//================= HASH ==================
func (t *Data) HGetAll(key string) (values []interface{}, err error) {
	values, err = redis.Values(t.conn.Do("HGETALL", key))
	if err == redis.ErrNil {
		err = nil
	}

	return values, err
}

func (t *Data) HGet(key string, field string) (value []byte, err error) {
	value, err = redis.Bytes(t.conn.Do("HGET", key, field))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func (t *Data) HLen(key string) (len int, err error) {
	len, err = redis.Int(t.conn.Do("HLEN", key))
	if err == redis.ErrNil {
		err = nil
	}
	return len, err
}

func (t *Data) HSet(key string, field string, value []byte) (err error) {
	_, err = t.conn.Do("HSET", key, field, value)
	return
}

func (t *Data) HMSet(key string, fields ...[]byte) (err error) {
	_, err = t.conn.Do("HMSET", key, fields)
	return err
}

func (t *Data) HDel(key string, field string) (num int, err error) {
	num, err = redis.Int(t.conn.Do("HDEL", key, field))
	return num, err
}

//func (t *Data) HDel(key string, field ...interface {}) (num int, err error) {
//TODO: HMDel still has problems
func (t *Data) HMDel(key string, field []string) (num int, err error) {
	num, err = redis.Int(t.conn.Do("HDEL", key, field))
	return num, err
}

//================= ZSET ==================
func (t *Data) ZAdd(key string, member string, score int32) (err error) {
	_, err = t.conn.Do("ZADD", key, score, member)
	return err
}

func (t *Data) ZRem(key string, member string) (err error) {
	_, err = t.conn.Do("ZREM", key, member)
	return err
}

func (t *Data) ZRangeByScore(key string, min int, max int, offset int, count int) (values []interface{}, err error) {
	//log.T("[TRACE] ZRangeByScore key:%v %v %v limit %v %v", key, min, max, offset, count)
	values, err = redis.Values(t.conn.Do("ZRANGEBYSCORE", key, min, max, "limit", offset, count))
	return values, err
}

func (t *Data) ZCount(key string, min int, max int) (count int, err error) {
	count, err = redis.Int(t.conn.Do("ZCOUNT", key, min, max))
	return count, err
}

func (t *Data) ZRemRangeByRank(key string, start, stop int) (err error) {
	_, err = t.conn.Do("ZREMRANGEBYRANK", key, start, stop)
	return err
}