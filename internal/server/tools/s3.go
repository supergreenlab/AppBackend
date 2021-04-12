/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

package tools

import (
	"fmt"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/sirupsen/logrus"
)

type S3Path struct {
	Path   *string
	Bucket string
}

type S3FileHolder interface {
	SetURLs(paths []string)
	GetURLs() (paths []S3Path)
}

type S3FileHolders interface {
	AsFeedMediasArray() []S3FileHolder
}

func LoadFeedMediaPublicURLs(fm S3FileHolder) error {
	paths := fm.GetURLs()
	expiry := time.Minute * 60

	results := make([]string, len(paths))
	for i, p := range paths {
		if p.Path == nil {
			results[i] = ""
			continue
		}
		cacheKey := fmt.Sprintf("cache.presigned-%s-%s", p.Bucket, *p.Path)
		if v, err := kv.GetString(cacheKey); err == nil {
			results[i] = v
		} else {
			url1, err := storage.Client.PresignedGetObject(p.Bucket, *p.Path, expiry, nil)
			if err != nil {
				return err
			}
			results[i] = url1.RequestURI()
			if err := kv.SetStringWithExpiration(cacheKey, results[i], time.Minute*45); err != nil {
				logrus.Errorf("kv.SetStringWithExpiration in LoadFeedMediaPublicURLs: %q - cacheKey: %s", err, cacheKey)
			}
		}
	}
	fm.SetURLs(results)
	return nil
}
