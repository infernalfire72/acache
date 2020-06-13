package main

import (
	"net/http"
	"strconv"
	"runtime"
)

func Start() {
	http.HandleFunc("/", getHpHandler)
	http.HandleFunc("/beatmap/", getBeatmapHandler)
	http.HandleFunc("/leaderboard/", getLeaderboardHandler)
	log.Info("Staring API")
	http.ListenAndServe(":5000", nil)
}

func getHpHandler(w http.ResponseWriter, req *http.Request) {
	var m runtime.MemStats
    runtime.ReadMemStats(&m)

	w.Write([]byte(strconv.FormatUint(m.TotalAlloc, 10)))
}

func getBeatmapHandler(w http.ResponseWriter, req *http.Request) {

}

func getLeaderboardHandler(w http.ResponseWriter, req *http.Request) {
	qs := req.URL.Query()

	hash := qs.Get("md5")
	if len(hash) != 32 {
		w.WriteHeader(400)
		w.Write([]byte("error: invalid hash"))
		return
	}

	mode, err := strconv.ParseInt(qs.Get("m"), 10, 8)
	if err != nil || mode < 0 || mode > 3 {
		mode = 0
	}

	rx, err := strconv.ParseBool(qs.Get("rx"))
	if err != nil {
		rx = false
	}

	limit, err := strconv.ParseInt(qs.Get("limit"), 10, 32)
	if err != nil {
		limit = 50
	}

	mods, err := strconv.ParseInt(qs.Get("mods"), 10, 32)
	if err != nil {
		mods = -1
	}

	// We can ignore the user if no user is present
	u, _ := strconv.ParseInt(qs.Get("user"), 10, 32)

	sw := Stopwatch{}
	sw.Start()
	lb := lcache.Get(struct {
		string
		byte
		bool
	}{hash, byte(mode), rx})
	bmap := lb.Map()

	output := bmap.String(len(lb.Scores))

	if u > 0 {
		personalBest, position := lb.FindUserScore(int(u))
		if personalBest != nil {
			output += personalBest.String(!lb.Relax || bmap.Status == Loved, position + 1)
		} else {
			output += "\n"
		}
	} else {
		output += "\n"
	}
	pos := 0
	for _, score := range lb.Scores {
		if pos >= int(limit) {
			break
		}

		// We have applied a mod filter
		if mods >= 0 && score.Mods != int(mods) {
			continue
		}

		output += score.String(!lb.Relax || bmap.Status == Loved, pos + 1)
		pos++
	}
	
	w.Write([]byte(output))
	sw.Stop()
	log.Infof("Served Leaderboard for %s[%t, %d, %d] in %s", hash, rx, mode, limit, sw.ElapsedReadable())
}