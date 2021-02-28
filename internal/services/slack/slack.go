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

package slack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("slackwebhook", "", "Webhook url for the slack notifications")
)

func listenFeedEntriesAdded() {
	ch := pubsub.SubscribeOject("insert.feedentries")
	for c := range ch {
		fe := c.(middlewares.InsertMessage).Object.(*db.FeedEntry)
		id := c.(middlewares.InsertMessage).ID

		plant, err := db.GetPlantForFeedEntryID(id)
		if err != nil {
			logrus.Errorf("db.GetPlantForFeedEntryID in listenFeedEntriesAdded %q - id: %s fe: %+v", err, id, fe)
			continue
		}
		if !plant.Public {
			continue
		}
		PublicDiaryEntryPosted(id, *fe, plant)
	}
}

func PublicDiaryEntryPosted(id uuid.UUID, fe db.FeedEntry, p db.Plant) {
	params := map[string]interface{}{}
	if err := json.Unmarshal([]byte(fe.Params), &params); err != nil {
		logrus.Errorf("json.Unmarshal in PublicDiaryEntryPosted %q - %+v", err, fe)
	}
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("New public diary entry on the plant %s", p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel> New public diary entry on the plant %s:\n\n%s\n%s\n\n(id: %s)\n<https://supergreenlab.com/public/plant?id=%s&feid=%s>", p.Name, fe.Type, params["message"], id, p.ID.UUID, id),
		Footer:        fmt.Sprintf("on %s", p.Name),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("slack.PostWebhook in PublicDiaryEntryPosted %q - id: %s fe: %+v p: %+v", err, id, fe, p)
	}
}

func CommentPosted(id uuid.UUID, com db.Comment, p db.Plant, u db.User) {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("Comment posted by %s on the plant %s", u.Nickname, p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel> Comment posted by %s on the plant %s:\n\n%s\n\n(id: %s)\n<https://supergreenlab.com/public/plant?id=%s&feid=%s>", u.Nickname, p.Name, com.Text, id, p.ID.UUID, com.FeedEntryID),
		Footer:        fmt.Sprintf("by %s", u.Nickname),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("slack.PostWebhook in CommentPosted %q - com: %+v p: %+v u: %+v", err, com, p, u)
	}
}

func CommentLikeAdded(l db.Like, com db.Comment, p db.Plant, u db.User) {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("%s liked a comment on the plant %s", u.Nickname, p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel> %s liked a comment on the plant %s\n\n%s\n\n<https://supergreenlab.com/public/plant?id=%s&feid=%s>", u.Nickname, p.Name, com.Text, p.ID.UUID, com.FeedEntryID),
		Footer:        fmt.Sprintf("by %s", u.Nickname),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("slack.PostWebhook in CommentLikeAdded %q - l: %+v com: %+v p: %+v u: %+v", err, l, com, p, u)
	}
}

func PostLikeAdded(l db.Like, p db.Plant, u db.User) {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("%s liked a diary entry on the plant %s", u.Nickname, p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel>%s liked a diary entry on the plant %s\n\n<https://supergreenlab.com/public/plant?id=%s&feid=%s>", u.Nickname, p.Name, p.ID.UUID, l.FeedEntryID.UUID),
		Footer:        fmt.Sprintf("by %s", u.Nickname),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("slack.PostWebhook in PostLikeAdded %q - l: %+v p: %+v u: %+v", err, l, p, u)
	}
}

func init() {
	viper.SetDefault("SlackWebhook", "")
}

func Init() {
	go listenFeedEntriesAdded()
}
