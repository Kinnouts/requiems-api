package counter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var incrementIfPresentScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[1]) == 0 then
  return -1
end
local val = redis.call("INCR", KEYS[1])
redis.call("SADD", KEYS[2], ARGV[1])
return val
`)

var incrementWithBootstrapScript = redis.NewScript(`
if redis.call("EXISTS", KEYS[1]) == 0 then
  redis.call("SET", KEYS[1], ARGV[1])
end
local val = redis.call("INCR", KEYS[1])
redis.call("SADD", KEYS[2], ARGV[2])
return val
`)

func incrementIfPresent(ctx context.Context, rdb *redis.Client, namespace string) (int64, bool, error) {
	res, err := incrementIfPresentScript.Run(ctx, rdb, []string{redisKey(namespace), dirtySetKey}, namespace).Result()

	if err != nil {
		return 0, false, err
	}

	val, err := toInt64(res)
	if err != nil {
		return 0, false, err
	}

	if val == -1 {
		return 0, true, nil
	}

	return val, false, nil
}

func incrementWithBootstrap(ctx context.Context, rdb *redis.Client, namespace string, baseline int64) (int64, error) {
	res, err := incrementWithBootstrapScript.Run(
		ctx,
		rdb,
		[]string{redisKey(namespace), dirtySetKey},
		strconv.FormatInt(baseline, 10),
		namespace,
	).Result()

	if err != nil {
		return 0, err
	}

	return toInt64(res)
}

func toInt64(raw any) (int64, error) {
	switch v := raw.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unexpected redis script return type %T", raw)
	}
}
