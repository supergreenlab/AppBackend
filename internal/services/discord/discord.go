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
	"encoding/json"
	"fmt"
	"log"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/bwmarrin/discordgo"
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
	sentEntries = map[uuid.UUID]bool{}
)

func init() {
	viper.SetDefault("DiscordToken", "")
	viper.SetDefault("DiscordPublicPostChannel", "")
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

		minioClient := storage.CreateMinioClient()
		obj, err := minioClient.GetObject("feedmedias", filePath, minio.GetObjectOptions{})
		if err != nil {
			logrus.Errorf("minioClient.GetObject in listenFeedMediasAdded %q - id: %s fm: %+v", err, id, fm)
			continue
		}

		params := map[string]interface{}{}
		if err := json.Unmarshal([]byte(fe.Params), &params); err != nil {
			logrus.Errorf("json.Unmarshal in listenFeedMediasAdded %q - %+v", err, fe)
		}
		msg := ""
		if sentEntries[fe.ID.UUID] == false {
			paramMsg, _ := params["message"].(string)
			msg = fmt.Sprintf("**New diary entry for the Plant \"%s\"**", plant.Name)
			if paramMsg != "" {
				msg = fmt.Sprintf("%s\n\n*%s*\n\n", msg, paramMsg)
			}
			msg = fmt.Sprintf("%s\nCheck it out here: https://supergreenlab.com/public/plant?id=%s&feid=%s", msg, plant.ID.UUID, fe.ID.UUID)
		}
		sentEntries[fe.ID.UUID] = true
		_, err = s.ChannelFileSendWithMessage(viper.GetString("DiscordPublicPostChannel"), msg, "pic.jpg", obj)
		if err != nil {
			logrus.Errorf("s.ChannelFileSendWithMessage in listenFeedMediasAdded %q - id: %s fm: %+v", err, id, fm)
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
}
