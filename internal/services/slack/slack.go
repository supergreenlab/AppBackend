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
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	slackWebhook = pflag.String("slackwebhook", "", "Webhook url for the slack notifications")
)

func CommentPosted(com db.Comment, p db.Plant, u db.User) {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("Comment posted on the plant %s", p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel> %s\n(id: %s)\n<https://supergreenlab.com/public/plant?id=%s&feid=%s>", com.Text, com.ID.UUID, p.ID.UUID, com.FeedEntryID),
		Footer:        fmt.Sprintf("by %s", u.Nickname),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("%q", err)
	}
}

func CommentLikeAdded(l db.Like, com db.Comment, p db.Plant, u db.User) {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("Liked a comment on the plant %s", p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel> %s\n<https://supergreenlab.com/public/plant?id=%s&feid=%s>", com.Text, p.ID.UUID, com.FeedEntryID),
		Footer:        fmt.Sprintf("by %s", u.Nickname),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("%q", err)
	}
}

func PostLikeAdded(l db.Like, p db.Plant, u db.User) {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      fmt.Sprintf("Liked a diary entry on the plant %s", p.Name),
		AuthorName:    "SuperGreenApp",
		AuthorSubname: "supergreenlab.com/app",
		AuthorLink:    "https://www.supergreenlab.com",
		AuthorIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Text:          fmt.Sprintf("<!channel><https://supergreenlab.com/public/plant?id=%s&feid=%s>", p.ID.UUID, l.FeedEntryID.UUID),
		Footer:        fmt.Sprintf("by %s", u.Nickname),
		FooterIcon:    "https://www.supergreenlab.com/_nuxt/img/icon_sgl_basics.709180a.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(viper.GetString("SlackWebhook"), &msg)
	if err != nil {
		logrus.Errorf("%q", err)
	}
}

func init() {
	viper.SetDefault("SlackWebhook", "")
}
