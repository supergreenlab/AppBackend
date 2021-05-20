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

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("exportdir", "sgl_export", "Export directory")
	_ = pflag.String("miniodir", "/minio", "Minio fs location")
)

func init() {
	viper.SetDefault("ExportDir", "sgl_export")
	viper.SetDefault("MinioDir", "/minio")
}

type IndexData struct {
	Plants []appbackend.Plant
}

type FeedEntry struct {
	appbackend.FeedEntry
	FeedMedias []appbackend.FeedMedia
}

type PlantData struct {
	Plant       appbackend.Plant
	FeedEntries []FeedEntry
}

func copyTo(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	return err
}

func main() {
	exportDir := viper.GetString("ExportDir")
	minioDir := viper.GetString("MinioDir")

	if err := os.Mkdir(exportDir, 0755); err != nil {
		logrus.Fatalf("os.Mkdir in main %q", err)
	}

	InitConfig()
	db.Init()

	plants := make([]appbackend.Plant, 0)

	if err := db.Sess.Select("*").From("plants").Where("deleted=false").And("is_public=true").OrderBy("uat desc").All(&plants); err != nil {
		logrus.Fatalf("%q", err)
	}

	indexData := IndexData{
		Plants: plants,
	}

	if indexFile, err := os.Create(fmt.Sprintf("%s/index.json", exportDir)); err != nil {
		logrus.Fatalf("os.Open in main %q", err)
	} else if err := json.NewEncoder(indexFile).Encode(indexData); err != nil {
		logrus.Fatalf("json.NewEncoder in main %q", err)
		indexFile.Close()
	}

	for _, plant := range plants {
		fes := []appbackend.FeedEntry{}
		if err := db.Sess.Select("*").From("feedentries").Where("feedid = ?", plant.FeedID).And("deleted=false").OrderBy("cat desc").All(&fes); err != nil {
			logrus.Fatalf("%q", err)
		}

		feedEntries := []FeedEntry{}
		for _, feedEntry := range fes {
			fms := []appbackend.FeedMedia{}
			if err := db.Sess.Select("*").From("feedmedias").Where("feedentryid = ?", feedEntry.ID).And("deleted=false").All(&fms); err != nil {
				logrus.Fatalf("db.Sess.Select in main %q", err)
			}

			for _, feedMedia := range fms {
				if err := copyTo(fmt.Sprintf("%s/feedmedias/%s", minioDir, feedMedia.FilePath), fmt.Sprintf("%s/feedmedias/%s", exportDir, feedMedia.FilePath)); err != nil {
					logrus.Fatalf("copyTo in main %q", err)
				}
				if err := copyTo(fmt.Sprintf("%s/feedmedias/%s", minioDir, feedMedia.ThumbnailPath), fmt.Sprintf("%s/feedmedias/%s", exportDir, feedMedia.ThumbnailPath)); err != nil {
					logrus.Fatalf("copyTo in main %q", err)
				}
			}

			feedEntries = append(feedEntries, FeedEntry{
				FeedEntry:  feedEntry,
				FeedMedias: fms,
			})
		}
		plantData := PlantData{
			Plant:       plant,
			FeedEntries: feedEntries,
		}
		if plantFile, err := os.Create(fmt.Sprintf("%s/plant-%s.json", exportDir, plant.ID.UUID)); err != nil {
			logrus.Fatalf("os.Open in main %q", err)
		} else if err := json.NewEncoder(plantFile).Encode(plantData); err != nil {
			logrus.Fatalf("json.NewEncoder in main %q", err)
			plantFile.Close()
		}
	}
}
