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

func GetNum(key string) (float64, error) {
	n, err := r.Get(key).Float64()
	if errors.Is(err, redis.Nil) {
		n = 0
		err = nil
	}
	return n, err
}

func GetInt(key string) (int, error) {
	n, err := r.Get(key).Int()
	if errors.Is(err, redis.Nil) {
		n = 0
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

func init() {
	viper.SetDefault("RedisURL", "redis:6379")
	r = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("RedisURL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
