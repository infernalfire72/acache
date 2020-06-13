package api

import (
	"strconv"
	"runtime"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/infernalfire72/acache/beatmaps"
	"github.com/infernalfire72/acache/leaderboards"
	"github.com/infernalfire72/acache/log"
	"github.com/infernalfire72/acache/tools"
)

func Start() {
	r := router.New()
	r.GET("/", MemHandler)
	r.GET("/beatmap/", BeatmapHandler)
	r.GET("/leaderboard/", LeaderboardHandler)
	log.Info("Starting API")
	log.Error(fasthttp.ListenAndServe(":5000", r.Handler))
}

func MemHandler(ctx *fasthttp.RequestCtx) {
	var m runtime.MemStats
    runtime.ReadMemStats(&m)

	ctx.WriteString(strconv.FormatUint(m.TotalAlloc, 10))
}

func BeatmapHandler(ctx *fasthttp.RequestCtx) {

}

func LeaderboardHandler(ctx *fasthttp.RequestCtx) {
	qs := ctx.QueryArgs()

	hash := string(qs.Peek("md5"))
	if len(hash) != 32 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString("error: invalid hash")
		return
	}

	mode := qs.GetUintOrZero("m")
	if mode > 3 {
		mode = 0
	}

	rx := qs.GetBool("rx")

	limit, err := qs.GetUint("limit")
	if err != nil {
		limit = 50
	}

	// We can ignore the error here because it will default to -1 when there is no arg!
	mods, _ := qs.GetUint("mods")

	u := qs.GetUintOrZero("user")

	fl := qs.GetBool("friends") && u != 0 // jg vs jne who will win
	var friendsFilter []int
	if fl {
		friendsFilter = tools.GetFriends(u)
	}

	sw := tools.Stopwatch{}
	sw.Start()
	lb := leaderboards.Cache.Get(leaderboards.Identifier{hash, byte(mode), rx})
	bmap := lb.Map()

	ctx.WriteString(bmap.String(len(lb.Scores)))

	if u != 0 {
		personalBest, position := lb.FindUserScore(int(u))
		if personalBest != nil {
			ctx.WriteString(personalBest.String(!lb.Relax || bmap.Status == beatmaps.Loved, position + 1))
		} else {
			ctx.WriteString("\n")
		}
	} else {
		ctx.WriteString("\n")
	}
	pos := 0
	for _, score := range lb.Scores {
		if pos >= int(limit) {
			break
		}

		// We have applied a mod filter
		if mods >= 0 && score.Mods != int(mods) {
			continue
		} else if fl && !tools.Has(friendsFilter, score.UserID) {
			continue
		}

		ctx.WriteString(score.String(!lb.Relax || bmap.Status == beatmaps.Loved, pos + 1))
		pos++
	}
	
	sw.Stop()
	log.Infof("Served Leaderboard for %s[%t, %d, %d] in %s", hash, rx, mode, limit, sw.ElapsedReadable())
}