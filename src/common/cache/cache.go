package cache

import (
	"fmt"
)

var cacheStore *Store

func InitCache(redisAddr, redisPassword string) {
	cacheStore = NewStore(redisAddr, redisPassword, 60, 10, true)
}

func GetCacheKey(entity string, appId, versionId, contentType, contentId uint32, args ...interface{}) string {
	key := fmt.Sprintf("%s_app_id_%d_version_id_%d_content_type_%d_content_id_%d", entity, appId, versionId, contentType, contentId)
	for _, k := range args {
		key += fmt.Sprint(k)
	}
	return key
}

func GetCacheObject(key string, obj interface{}, f StoreLoadFunc) error {
	return cacheStore.GetJsonObject(key, obj, f)
}

func GetCacheObjectWithExpire(key string, obj interface{}, ttl int, f StoreLoadFunc) error {
	return cacheStore.GetJsonObjectWithExpire(key, obj, ttl, f)
}
