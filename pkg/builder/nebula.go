package builder

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mxpv/podsync/pkg/feed"
	"github.com/mxpv/podsync/pkg/model"
)

type NebulaBuilder struct {
}

func (neb *NebulaBuilder) Build(ctx context.Context, cfg *feed.Config) (*model.Feed, error) {
	info, err := ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	fp := gofeed.Parser{}
	rssFeed, err := fp.ParseURL("https://rss.nebula.app/video/channels/" + info.ItemID + ".rss")
	if err != nil {
		return nil, err
	}

	_feed := &model.Feed{
		ItemID:    info.ItemID,
		Provider:  info.Provider,
		LinkType:  info.LinkType,
		Format:    cfg.Format,
		Quality:   cfg.Quality,
		PageSize:  cfg.PageSize,
		UpdatedAt: time.Now().UTC(),

		Title:       rssFeed.Title,
		ItemURL:     rssFeed.Link,
		Description: rssFeed.Description,
		PubDate:     *rssFeed.UpdatedParsed,
	}

	// setup episodes and add to feed
	added := 0
	for _, item := range rssFeed.Items {
		_feed.Episodes = append(_feed.Episodes, &model.Episode{
			ID:          item.GUID,
			Title:       item.Title,
			Description: item.Description,
			VideoURL:    item.Link,
			PubDate:     *item.PublishedParsed,
			// Thumbnail: item.,
			// Duration: item.length,
			// Size: ,
			Status: model.EpisodeNew,
		})

		added++

		if added >= _feed.PageSize {
			break
		}
	}

	return _feed, nil
}

func newNebulaBuilder() (*NebulaBuilder, error) {
	return &NebulaBuilder{}, nil
}
