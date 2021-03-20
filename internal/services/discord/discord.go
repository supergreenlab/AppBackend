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

package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	"github.com/gofrs/uuid"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	s           *discordgo.Session
	_           = pflag.String("discordtoken", "", "Discord server token")
	_           = pflag.String("discordpublicpostchannel", "", "Public post discord channel")
	_           = pflag.String("discordlinkbookmarkchannel", "", "Link bookmark discord channel")
	sentEntries = map[uuid.UUID]bool{}
)

func init() {
	viper.SetDefault("DiscordToken", "")
	viper.SetDefault("DiscordPublicPostChannel", "")
	viper.SetDefault("DiscordLinkBookmarkChannel", "")
}

func listenFeedMediasAdded() {
	ch := pubsub.SubscribeOject("insert.feedmedias")
	for c := range ch {
		fm := c.(middlewares.InsertMessage).Object.(*db.FeedMedia)
		id := c.(middlewares.InsertMessage).ID

		plant, err := db.GetPlantForFeedEntryID(fm.FeedEntryID)
		if err != nil {
			logrus.Errorf("db.GetPlantForFeedEntryID in listenFeedMediasAdded %q - id: %s fm: %+v", err, id, fm)
			continue
		}
		if !plant.Public {
			continue
		}
		fe, err := db.GetFeedEntry(fm.FeedEntryID)
		if err != nil {
			logrus.Errorf("db.GetFeedEntry in listenFeedMediasAdded %q - id: %s fm: %+v", err, id, fm)
			continue
		}
		filePath := fm.FilePath
		if filePath[len(filePath)-3:] == "mp4" {
			filePath = fm.ThumbnailPath
		}

		obj, err := storage.Client.GetObject("feedmedias", filePath, minio.GetObjectOptions{})
		if err != nil {
			logrus.Errorf("minioClient.GetObject in listenFeedMediasAdded %q - id: %s fm: %+v", err, id, fm)
			continue
		}

		msg := ""
		if sentEntries[fe.ID.UUID] == false {
			msg = fmt.Sprintf("**New diary entry for the Plant \"%s\"**", plant.Name)
			plantSettings := map[string]interface{}{}
			if err := json.Unmarshal([]byte(plant.Settings), &plantSettings); err == nil {
				if strain, ok := plantSettings["strain"].(string); ok && strain != "" {
					msg = fmt.Sprintf("%s\nStrain: **%s**", msg, strain)
					if seedBank, ok := plantSettings["seedBank"].(string); ok && seedBank != "" {
						msg = fmt.Sprintf("%s from **%s**", msg, seedBank)
					}
					if plantType, ok := plantSettings["plantType"].(string); ok && plantType != "" {
						msg = fmt.Sprintf("%s (%s)", msg, plantType)
					}
				}

				if germinationDateStr, ok := plantSettings["germinationDate"].(string); ok && germinationDateStr != "" {
					if germinationDate, err := time.Parse(time.RFC3339, germinationDateStr); err == nil {
						msg = fmt.Sprintf("%s\n**Germinating since**: %d days", msg, int(time.Since(germinationDate).Hours()/24))
					}
				}

				if bloomingDateStr, ok := plantSettings["bloomingStart"].(string); ok && bloomingDateStr != "" {
					if bloomingDate, err := time.Parse(time.RFC3339, bloomingDateStr); err == nil {
						msg = fmt.Sprintf("%s\n**Blooming since**: %d days", msg, int(time.Since(bloomingDate).Hours()/24))
					}
				}
			} else {
				logrus.Errorf("json.Unmarshal in listenFeedMediasAdded %q - %+v", err, plant)
			}

			params := map[string]interface{}{}
			if err := json.Unmarshal([]byte(fe.Params), &params); err == nil {
				paramMsg, _ := params["message"].(string)
				if paramMsg != "" {
					msg = fmt.Sprintf("%s\n\n*%s*\n\n", msg, paramMsg)
				}
			} else {
				logrus.Errorf("json.Unmarshal in listenFeedMediasAdded %q - %+v", err, fe)
			}

			msg = fmt.Sprintf("%s\n\n_Drop a comment/like here: <https://supergreenlab.com/public/plant?id=%s&feid=%s>_ 💚", msg, plant.ID.UUID, fe.ID.UUID)
		}
		sentEntries[fe.ID.UUID] = true

		img, err := imaging.Decode(obj, imaging.AutoOrientation(true))
		if err != nil {
			logrus.Errorf("image.Decode in listenFeedMediasAdded %q - %+v", err, fe)
			continue
		}
		var resized image.Image
		if img.Bounds().Max.X > img.Bounds().Max.Y {
			resized = imaging.Resize(img, 1250, 0, imaging.Lanczos)
		} else {
			resized = imaging.Resize(img, 0, 1250, imaging.Lanczos)
		}

		box, err := db.GetBox(plant.BoxID)
		if err != nil {
			logrus.Errorf("db.GetBox in listenFeedMediasAdded %q - plant: %+v", err, plant)
			continue
		}
		var device *db.Device
		if box.DeviceID.Valid {
			d, err := db.GetDevice(box.DeviceID.UUID)
			if err != nil {
				logrus.Errorf("db.GetDevice in listenFeedMediasAdded %q - device: %+v", err, device)
				continue
			}
			device = &d
		}

		buff := new(bytes.Buffer)
		err = jpeg.Encode(buff, resized, &jpeg.Options{Quality: 80})
		if err != nil {
			fmt.Println("failed to create buffer", err)
			continue
		}

		buff, err = addSGLOverlays(box, plant, device, buff)
		if err != nil {
			logrus.Errorf("addSGLOverlays in listenFeedMediasAdded %q - device: %+v", err, device)
			continue
		}

		jpegimg := bytes.NewReader(buff.Bytes())

		_, err = s.ChannelFileSendWithMessage(viper.GetString("DiscordPublicPostChannel"), msg, "pic.jpg", jpegimg)
		if err != nil {
			logrus.Errorf("s.ChannelFileSendWithMessage in listenFeedMediasAdded %q - id: %s fm: %+v", err, id, fm)
			continue
		}
	}
}

func listenLinkBookmarksAdded() {
	ch := pubsub.SubscribeOject("insert.linkbookmarks")
	for c := range ch {
		lb := c.(middlewares.InsertMessage).Object.(*db.LinkBookmark)
		id := c.(middlewares.InsertMessage).ID

		user, err := db.GetUser(lb.UserID)
		if err != nil {
			logrus.Errorf("db.GetUser in listenLinkBookmarksAdded %q - id: %s lb: %+v", err, id, lb)
			continue
		}

		msg := fmt.Sprintf("**New bookmark posted by %s**\n\n%s", user.Nickname, lb.URL)
		_, err = s.ChannelMessageSend(viper.GetString("DiscordLinkBookmarkChannel"), msg)
		if err != nil {
			logrus.Errorf("s.ChannelMessageSend in listenLinkBookmarksAdded %q - id: %s lb: %+v", err, id, lb)
			continue
		}
	}
}

func Init() {
	var err error
	s, err = discordgo.New("Bot " + viper.GetString("DiscordToken"))
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}

	//s.AddHandler(messageCreate)

	if err = s.Open(); err != nil {
		log.Fatalln("opening websocket failed", err)
	}

	go listenFeedMediasAdded()
	go listenLinkBookmarksAdded()
}
