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
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/config"
	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/services/bot"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	db.Init()
	kv.Init()
	storage.Init()

	id, err := uuid.FromString("8333ec06-9e82-4f02-ac21-f7d1e7ad810c")
	if err != nil {
		logrus.Fatalf("uuid.FromString in main %q", err)
	}

	timelapse, err := db.GetTimelapse(id)
	if err != nil {
		logrus.Fatalf("db.GetTimelapse in main %q", err)
	}

	t := time.Now()
	from := t.Add(-7 * 24 * time.Hour)
	to := t
	if err := bot.SendTimelapseRequest(from, to, timelapse); err != nil {
		logrus.Fatalf("bot.SendTimelapseRequest in main %q", err)
	}
}
