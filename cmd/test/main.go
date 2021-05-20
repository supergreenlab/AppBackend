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

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/sirupsen/logrus"
)

func main() {
	t := time.Now()
	from := t.Add(-24 * time.Hour)
	to := t

	device := appbackend.Device{Identifier: "e4a0c3ab6224"}
	ts, err := appbackend.LoadGraphValue(device, from, to, "BOX", "TEMP", 0)
	if err != nil {
		logrus.Fatalf("%q", err)
	}
	logrus.Fatalf("%+v", ts)
}
