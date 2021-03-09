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

package storage

import (
	"fmt"
	"log"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

var (
	Client *minio.Client
)

// SetupBucket - create bucket if not exists
func SetupBucket(name string) {
	exists, err := Client.BucketExists(name)
	if err != nil {
		log.Fatalln(err)
	}
	if exists {
		log.Printf("Already created bucket: %s\n", name)
		return
	}
	err = Client.MakeBucket(name, "")
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateMinioClient - creates an initialized minio client
func createMinioClient() *minio.Client {
	accessKey := viper.GetString("S3AccessKey")
	secretKey := viper.GetString("S3SecretKey")
	host := viper.GetString("S3Host")
	secure := viper.GetString("S3Secure") == "true"
	minioClient, err := minio.New(host, accessKey, secretKey, secure)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return minioClient
}

func Init() {
	Client = createMinioClient()
}
