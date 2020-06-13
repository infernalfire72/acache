package beatmaps

import "fmt"

type BeatmapStatus int

const (
	Unknown BeatmapStatus = iota - 2
	NotSubmitted
	Pending
	NeedsUpdate
	Ranked
	Approved
	Qualified
	Loved
)

type Beatmap struct {
	Md5			string
	ID			int
	SetID		int
	Name		string
	Status		BeatmapStatus
	Playcount	int
	Passcount	int
}

func (b *Beatmap) String(scoresCount int) string {
	return fmt.Sprintf("%d|false|%d|%d|%d\n0\n%s\n10.0\n", b.Status, b.ID, b.SetID, scoresCount, b.Name)
}