#!/bin/bash

# Copyright (C) 2018  SuperGreenLab <towelie@supergreenlab.com>
# Author: Constantin Clauzel <constantin.clauzel@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

docker build -t timelapse-dev . -f Dockerfile.timelapse.dev
docker run --rm -it -v $(pwd):/app --workdir /app --entrypoint=/usr/local/go/bin/go timelapse-dev build -v -o bin/timelapse cmd/timelapse/*.go
docker build -t supergreenlab/timelapse -f Dockerfile.timelapse .
