package feeds

import (
	"time"
)

func loadFeedMediaPublicURLs(fm FeedMedia) (FeedMedia, error) {
	expiry := time.Second * 60 * 60
	minioClient := createMinioClient()
	url1, err := minioClient.PresignedGetObject("feedmedias", fm.FilePath, expiry, nil)
	if err != nil {
		return fm, err
	}
	fm.FilePath = url1.RequestURI()

	url2, err := minioClient.PresignedGetObject("feedmedias", fm.ThumbnailPath, expiry, nil)
	if err != nil {
		return fm, err
	}
	fm.ThumbnailPath = url2.RequestURI()
	return fm, nil
}
