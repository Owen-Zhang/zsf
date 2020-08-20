package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//Del 删除key
func (r *Redis) Del(key string) {
	r.client.Del(context.Background(), key)
}

//Exists 判断key是否存在
func (r *Redis) Exists(key string) bool {
	value, err := r.client.Exists(context.Background(), key).Result()
	return err == nil || value == 1
}

//Expire 设置过期时间
func (r *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(context.Background(), key, expiration).Result()
}

/////////////////////////////////////////
///////////// string ////////////////////
/////////////////////////////////////////

//Set 保存string
func (r *Redis) Set(key string, value interface{}, expire time.Duration) error {
	return r.client.Set(context.Background(), key, value, expire).Err()
}

//Get 返回内容(不管错误,用时要注意)
func (r *Redis) Get(key string) string {
	cmd := r.client.Get(context.Background(), key)
	return cmd.Val()
}

//GetWithErr 获取string类型的值,返回内容及错误信息
func (r *Redis) GetWithErr(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

//MGet 一次性获取多个string类型的值
//由于在cluster集群下数据可能分布在很多slot中，几乎都会报错
//有一种方案就是客户端包装成多次调用，然后组合(待实现)
//现在最好不要用此方法
func (r *Redis) MGet(keys ...string) ([]string, error) {
	sliceObj := r.client.MGet(context.Background(), keys...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, value := range tmp {
		if value != nil {
			strSlice = append(strSlice, value.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

//Incr 自增+1,返回true表示成功,false表示失败
func (r *Redis) Incr(key string) bool {
	err := r.client.Incr(context.Background(), key).Err()
	return err == nil
}

//IncrWithErr 自增+1 返回自增后的结果,以及错误信息
func (r *Redis) IncrWithErr(key string) (int64, error) {
	return r.client.Incr(context.Background(), key).Result()
}

//IncrBy 将key所储存的值加上增量 increment。
func (r *Redis) IncrBy(key string, increment int64) int64 {
	intObj := r.client.IncrBy(context.Background(), key, increment)
	return intObj.Val()
}

//IncrByWithErr 将 key 所储存的值加上增量 increment 。
func (r *Redis) IncrByWithErr(key string, increment int64) (int64, error) {
	intObj := r.client.IncrBy(context.Background(), key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

//IncrByFloat 将key所储存的值加上增量 increment。
func (r *Redis) IncrByFloat(key string, increment float64) float64 {
	intObj := r.client.IncrByFloat(context.Background(), key, increment)
	return intObj.Val()
}

//IncrByFloatWithErr 将 key 所储存的值加上增量 increment 。
func (r *Redis) IncrByFloatWithErr(key string, increment float64) (float64, error) {
	intObj := r.client.IncrByFloat(context.Background(), key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

//DecrWithErr 自减1 返回自减后的结果,以及错误信息
func (r *Redis) DecrWithErr(key string) (int64, error) {
	return r.client.Decr(context.Background(), key).Result()
}

//DecrByWithErr 自减increment 返回自减后的结果,以及错误信息
func (r *Redis) DecrByWithErr(key string, increment int64) (int64, error) {
	return r.client.DecrBy(context.Background(), key, increment).Result()
}

//HGetAll 获取hash的所有键值对
func (r *Redis) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(context.Background(), key).Result()
}

// HGet 从redis获取hash单个值
func (r *Redis) HGet(key string, fields string) (string, error) {
	strObj := r.client.HGet(context.Background(), key, fields)
	err := strObj.Err()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return strObj.Val(), nil
}

// HMGetMap 批量获取hash值，返回map
func (r *Redis) HMGetMap(key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := r.client.HMGet(context.Background(), key, fields...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return make(map[string]string)
	}

	tmp := sliceObj.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func (r *Redis) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		err := r.client.HMSet(context.Background(), key, hash).Err()
		if err != nil {
			return false
		}
		if expire > 0 {
			r.client.Expire(context.Background(), key, expire)
		}
		return true
	}
	return false
}

//HSet 设置某个字段的值
func (r *Redis) HSet(key string, field string, value interface{}) bool {
	err := r.client.HSet(context.Background(), key, field, value).Err()
	return err == nil
}

//HDel 删除某些字段
func (r *Redis) HDel(key string, field ...string) bool {
	IntObj := r.client.HDel(context.Background(), key, field...)
	err := IntObj.Err()
	return err == nil
}

// HIncrBy 哈希field增加incr
func (r *Redis) HIncrBy(key string, field string, incr int) int64 {
	result, err := r.client.HIncrBy(context.Background(), key, field, int64(incr)).Result()
	if err != nil {
		return 0
	}
	return result
}

// HIncrByWithErr 哈希field自增并且返回错误
func (r *Redis) HIncrByWithErr(key string, field string, incr int) (int64, error) {
	return r.client.HIncrBy(context.Background(), key, field, int64(incr)).Result()
}

/////////////////////////////////////////////
///////////////////// list //////////////////
/////////////////////////////////////////////

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (r *Redis) LPush(key string, values ...interface{}) (int64, error) {
	return r.client.LPush(context.Background(), key, values...).Result()
}

//RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (r *Redis) RPush(key string, values ...interface{}) (int64, error) {
	return r.client.RPush(context.Background(), key, values...).Result()
}

//RPop 移除并返回列表 key 的尾元素。
func (r *Redis) RPop(key string) (string, error) {
	return r.client.RPop(context.Background(), key).Result()
}

//LRange 获取列表指定范围内的元素
func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(context.Background(), key, start, stop).Result()
}

//LLen list长度
func (r *Redis) LLen(key string) int64 {
	IntObj := r.client.LLen(context.Background(), key)
	if err := IntObj.Err(); err != nil {
		return 0
	}
	return IntObj.Val()
}

/////////////////////////////////////////////
//////////////// set(集合) //////////////////
/////////////////////////////////////////////

// SAdd 向set中添加成员
func (r *Redis) SAdd(key string, member ...interface{}) (int64, error) {
	return r.client.SAdd(context.Background(), key, member...).Result()
}

// SMembers 返回set的全部成员
func (r *Redis) SMembers(key string) ([]string, error) {
	return r.client.SMembers(context.Background(), key).Result()
}

//SIsMember 判断某个值是否在此set中
func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(context.Background(), key, member).Result()
}

/////////////////////其它//////////////

//GeoAdd 写入地理位置
func (r *Redis) GeoAdd(key string, location *redis.GeoLocation) (int64, error) {
	return r.client.GeoAdd(context.Background(), key, location).Result()
}

// GeoRadius 根据经纬度查询列表
func (r *Redis) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return r.client.GeoRadius(context.Background(), key, longitude, latitude, query).Result()
}

//Close 关闭连接
func (r *Redis) Close() (err error) {
	c, _ := r.client.(*redis.ClusterClient)
	return c.Close()
}
