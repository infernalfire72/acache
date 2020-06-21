package beatmaps

import (
	"database/sql"
	"sync"
	"time"

	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/log"
)

var bMutex sync.RWMutex
var Cache *BeatmapCache

type BeatmapCache struct {
	Beatmaps map[string]*Beatmap
}

func (c *BeatmapCache) Get(md5 string) *Beatmap {
	bMutex.RLock()
	bmp := c.Beatmaps[md5]
	bMutex.RUnlock()
	if bmp != nil {
		now := time.Now()
		if now.Sub(bmp.LastUpdate).Seconds() >= 60 {
			return c.UpdateCache(md5)
		}
		return bmp
	} else {
		return c.UpdateCache(md5)
	}
}

func (c *BeatmapCache) UpdateCache(md5 string) *Beatmap {
	b := &Beatmap{
		Md5: md5,
	}

	err := config.DB.QueryRow("SELECT beatmap_id, beatmapset_id, song_name, ranked, playcount, passcount FROM beatmaps WHERE beatmap_md5 = ?", md5).Scan(
		&b.ID, &b.SetID, &b.Name, &b.Status, &b.Playcount, &b.Passcount,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
	}
	b.LastUpdate = time.Now()
	bMutex.Lock()
	c.Beatmaps[md5] = b
	bMutex.Unlock()
	return b
}

func (c BeatmapCache) Clear() {
	bMutex.Lock()
	c.Beatmaps = make(map[string]*Beatmap)
	bMutex.Unlock()
}
