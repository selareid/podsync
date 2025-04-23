package builder

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mxpv/podsync/pkg/feed"
	"github.com/mxpv/podsync/pkg/model"
	"github.com/pkg/errors"
)

type fnGetDuration func(ctx context.Context, url string, extra_args ...string) (duration int64, err error)

type NebulaBuilder struct {
	token       string
	getDuration fnGetDuration
}

func splitThumbnail(desc string) string {
	a := strings.Split(desc, "<img src=")

	if len(a) < 2 {
		return ""
	}

	b := strings.Split(a[1], ">")
	c := strings.Split(b[0], "?format")

	return c[0]
}

func (neb *NebulaBuilder) Build(ctx context.Context, cfg *feed.Config) (*model.Feed, error) {
	info, err := ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}

	// add authorisation token to yt-dlp arguments
	cfg.YouTubeDLArgs = append(cfg.YouTubeDLArgs, "--add-headers", "Authorization: Token "+neb.token)

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
		newEpisode := &model.Episode{
			ID:          item.GUID,
			Title:       item.Title,
			Description: item.Description,
			VideoURL:    item.Link,
			PubDate:     *item.PublishedParsed,
			// Size: ,
			Status: model.EpisodeNew,
		}

		dur, err := neb.getDuration(ctx, item.Link, "--add-headers", "Authorization: Token "+neb.token)
		if err != nil {
			newEpisode.Duration = dur
		}

		thumbnailURL := splitThumbnail(item.Description)
		if thumbnailURL != "" {
			newEpisode.Thumbnail = thumbnailURL
			pattern := `<img src=` + regexp.QuoteMeta(thumbnailURL) + `\?[^>]+>`
			re := regexp.MustCompile(pattern)
			newEpisode.Description = re.ReplaceAllString(newEpisode.Description, "")
		}

		_feed.Episodes = append(_feed.Episodes, newEpisode)

		added++

		if added >= _feed.PageSize {
			break
		}
	}

	return _feed, nil
}

func newNebulaBuilder(key string, durationGetter fnGetDuration) (*NebulaBuilder, error) {
	if key == "" {
		return nil, errors.New("invalid key given for Nebula")
	}

	return &NebulaBuilder{token: key, getDuration: durationGetter}, nil
}
