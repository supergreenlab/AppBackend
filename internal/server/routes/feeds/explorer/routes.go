/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

package explorer

import (
	cmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
)

// Init -
func Init(router *httprouter.Router) {
	//anon := cmiddlewares.AnonStack()
	auth := cmiddlewares.AuthStack()
	optionalAuth := cmiddlewares.OptionalAuthStack()

	router.GET("/public/plants/followed", auth.Wrap(fetchLatestUpdatedFollowedPublicPlants))
	router.GET("/public/feedEntries/followed", auth.Wrap(fetchLatestFollowedFeedEntries))

	router.GET("/public/plants", optionalAuth.Wrap(fetchLatestUpdatedPublicPlants))
	router.GET("/public/plants/search", optionalAuth.Wrap(searchPublicPlants))
	router.GET("/public/feedEntries", optionalAuth.Wrap(fetchLatestPublicFeedEntries))
	router.GET("/public/plant/:id", optionalAuth.Wrap(fetchPublicPlant))
	router.GET("/public/plant/:id/feedEntries", optionalAuth.Wrap(fetchPublicPlantFeedEntries))
	router.GET("/public/feedEntries/commented", optionalAuth.Wrap(fetchLatestCommentedFeedEntries))
	router.GET("/public/liked", optionalAuth.Wrap(fetchLatestLikedFeedEntries))
	router.GET("/public/feedEntry/:id", optionalAuth.Wrap(fetchPublicFeedEntry))
	router.GET("/public/feedEntry/:id/feedMedias", optionalAuth.Wrap(fetchPublicEntryFeedMedias))
	router.GET("/public/feedMedia/:id", optionalAuth.Wrap(fetchPublicFeedMedia))
}
