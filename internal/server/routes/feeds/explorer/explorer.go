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

package feeds

type publicPlantResult struct {
	ID            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	FilePath      string `db:"filepath" json:"filePath"`
	ThumbnailPath string `db:"thumbnailpath" json:"thumbnailPath"`

	Followed bool `db:"followed" json:"followed"`

	Settings    string `db:"settings" json:"settings"`
	BoxSettings string `db:"boxsettings" json:"boxSettings"`
}

func (r *publicPlantResult) SetURLs(filePath string, thumbnailPath string) {
	r.FilePath = filePath
	r.ThumbnailPath = thumbnailPath
}

func (r publicPlantResult) GetURLs() (filePath string, thumbnailPath string) {
	filePath, thumbnailPath = r.FilePath, r.ThumbnailPath
	return
}

type publicPlantsResult struct {
	Plants []publicPlantResult `json:"plants"`
}
