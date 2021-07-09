/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package kv

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	r *redis.Client
)

func HasNumKey(key string) (bool, error) {
	n, err := r.Exists(key).Result()
	return n != 0, err
}

func GetNum(key string, def float64) (float64, error) {
	n, err := r.Get(key).Float64()
	if errors.Is(err, redis.Nil) {
		n = def
		err = nil
	}
	return n, err
}

func SetNum(key string, value float64) error {
	return r.Set(key, value, 0).Err()
}

func GetInt(key string, def int) (int, error) {
	n, err := r.Get(key).Int()
	if errors.Is(err, redis.Nil) {
		n = def
		err = nil
	}
	return n, err
}

func GetBool(key string) (bool, error) {
	i, err := r.Get(key).Int()
	if errors.Is(err, redis.Nil) {
		i = 0
		err = nil
	}
	return i != 0, err
}

func SetBool(key string, value bool, expiration time.Duration) error {
	i := 0
	if value {
		i = 1
	}
	err := r.Set(key, i, expiration).Err()
	return err
}

func GetString(key string) (string, error) {
	return r.Get(key).Result()
}

func SetString(key, value string) error {
	return r.Set(key, value, 0).Err()
}

func SetStringWithExpiration(key, value string, expiration time.Duration) error {
	return r.Set(key, value, expiration).Err()
}

func GetKeys(patterns []string) ([]string, error) {
	keys := []string{}
	for _, p := range patterns {
		if strings.ContainsRune(p, '*') {
			if ks, err := r.Keys(p).Result(); err != nil {
				return nil, err
			} else {
				keys = append(keys, ks...)
			}
		} else {
			keys = append(keys, p)
		}
	}
	return keys, nil
}

func GetValues(keys []string) (map[string]interface{}, error) {
	values, err := r.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}
	results := map[string]interface{}{}
	for i, v := range values {
		results[keys[i]] = v
	}
	return results, nil
}

func Init() {
	viper.SetDefault("RedisURL", "redis:6379")
	r = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("RedisURL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
