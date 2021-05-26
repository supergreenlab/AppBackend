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

type FeedEntry struct {
	appbackend.FeedEntry
	FeedMedias []appbackend.FeedMedia
}

type PlantData struct {
	Box         appbackend.Box
	Plant       appbackend.Plant
	FeedEntries []FeedEntry
}

func copyTo(src, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		s, err := os.Open(src)
		if err != nil {
			return err
		}
		defer s.Close()
		d, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer d.Close()
		_, err = io.Copy(d, s)
		return err
	} else if err != nil {
		return err
	}
	logrus.Infof("File %s already exists", dst)
	return nil
}

func main() {
	InitConfig()
	db.Init()

	exportDir := viper.GetString("ExportDir")
	minioDir := viper.GetString("MinioDir")

	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		if err := os.Mkdir(exportDir, 0700); err != nil {
			logrus.Fatalf("os.Mkdir in main %q", err)
		}
	}

	plants := make([]appbackend.Plant, 0)

	if err := db.Sess.Select("*").From("plants").Where("deleted=false").And("is_public=true").OrderBy("uat desc").All(&plants); err != nil {
		logrus.Fatalf("%q", err)
	}

	for _, plant := range plants {
		box := appbackend.Box{}
		if err := db.Sess.Select("*").From("boxes").Where("id = ?", plant.BoxID).One(&box); err != nil {
			logrus.Fatalf("%q", err)
		}

		fes := []appbackend.FeedEntry{}
		if err := db.Sess.Select("*").From("feedentries").Where("feedid = ?", plant.FeedID).And("deleted=false").OrderBy("cat asc").All(&fes); err != nil {
			logrus.Fatalf("%q", err)
		}

		feedEntries := []FeedEntry{}
		for _, feedEntry := range fes {
			fms := []appbackend.FeedMedia{}
			if err := db.Sess.Select("*").From("feedmedias").Where("feedentryid = ?", feedEntry.ID).And("deleted=false").All(&fms); err != nil {
				logrus.Fatalf("db.Sess.Select in main %q", err)
			}

			for _, feedMedia := range fms {
				if err := copyTo(fmt.Sprintf("%s/feedmedias/%s", minioDir, feedMedia.FilePath), fmt.Sprintf("%s/%s", exportDir, feedMedia.FilePath)); err != nil {
					logrus.Errorf("copyTo in main %q", err)
					continue
				}
				if err := copyTo(fmt.Sprintf("%s/feedmedias/%s", minioDir, feedMedia.ThumbnailPath), fmt.Sprintf("%s/%s", exportDir, feedMedia.ThumbnailPath)); err != nil {
					logrus.Errorf("copyTo in main %q", err)
					continue
				}
			}

			feedEntries = append(feedEntries, FeedEntry{
				FeedEntry:  feedEntry,
				FeedMedias: fms,
			})
		}

		plantData := PlantData{
			Box:         box,
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
