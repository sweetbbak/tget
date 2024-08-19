package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const megabyte = 1024 * 1024

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

// Download retrieves the download information for a torrent and
// returns it as a string.
func Downloaded(t *torrent.Torrent, showPercent bool) string {
	var (
		done    = t.BytesCompleted()
		total   = t.Length()
		percent = float64(done) / float64(total) * 100

		tail string
	)

	if showPercent {
		tail = fmt.Sprintf(" (%d%%)", uint64(percent))
	}

	return fmt.Sprintf(
		"%s/%s%s ↓",
		humanize.Bytes(uint64(done)),
		humanize.Bytes(uint64(total)),
		tail,
	)
}

// calculateDownloadRate calculates the download rate in MB/s.
func calculateDownloadRate(bytesCompleted, startSize int64, elapsedTime time.Duration) float64 {
	return float64(bytesCompleted-startSize) / elapsedTime.Seconds() / megabyte
}

// Peers retrieves the peer information for a torrent and returns it as
// a string.
func Peers(t *torrent.Torrent) string {
	stats := t.Stats()

	return fmt.Sprintf(
		"%d/%d peers",
		stats.ActivePeers,
		stats.TotalPeers,
	)
}

// Upload retrieves the amount of data seeded for a torrent and returns
// it as a string.
func Upload(t *torrent.Torrent) string {
	var (
		stats  = t.Stats()
		upload = stats.BytesWritten.Int64()
	)

	return fmt.Sprintf(
		"%s ↑",
		humanize.Bytes(uint64(upload)),
	)
}

// Get largest file inside of a Torrent
func GetLargestFile(t *torrent.Torrent) *torrent.File {
	var target *torrent.File
	var maxSize int64

	for _, file := range t.Files() {
		if maxSize < file.Length() {
			maxSize = file.Length()
			target = file
		}
	}
	return target
}

// returns a seed ratio compared to the entire torrent
func TorrentSeedRatio(t *torrent.Torrent) float64 {
	stats := t.Stats()
	seedratio := float64(stats.BytesWrittenData.Int64()) / float64(stats.BytesReadData.Int64())
	return seedratio
}

func TruncateString(s string, length int) string {
	var l int
	var sb strings.Builder

	// early return if string is shorter then requested length
	if length >= len(s) {
		return s
	}

	for _, r := range s {
		if l <= length {
			sb.WriteRune(r)

		} else {
			break
		}
		l++
	}
	return sb.String()
}

func SeedProgress(t *torrent.Torrent) {
	title := TruncateString(t.Name(), 100)
	fmt.Printf("name [%s]...\n", title)
	tlen := float64(t.Length())

	for {
		fmt.Print("\x1b7")       // save the cursor position
		fmt.Print("\x1b[2k")     // erase the current line
		defer fmt.Print("\x1b8") // restore the cursor position

		stats := t.Stats()
		upload := stats.ConnStats.BytesWritten.Int64()
		pc := float64(float64(upload)/tlen) * 100

		fmt.Printf("upload [%d] [%s|%3.1f]\n", upload, Upload(t), pc)

		time.Sleep(time.Millisecond * 800)
		if pc >= 100.00 {
			break
		}
	}
}

func Progress(t *torrent.Torrent) {
	title := TruncateString(t.Name(), 100)
	fmt.Printf("name [%s]...", title)

	p, _ := pterm.DefaultProgressbar.WithTotal(100).Start()
	p.RemoveWhenDone = true

	for !t.Complete.Bool() {
		pc := float64(t.BytesCompleted()) / float64(t.Length()) * 100
		numpeers := len(t.PeerConns())
		p.Increment().Current = int(pc)
		p.UpdateTitle(fmt.Sprintf("%v peers [%v]", t.Name(), numpeers))
		time.Sleep(time.Millisecond * 50)
	}
	p.Stop()
}

func Header() {
	// Initialize a big text display with the letters "P" and "Term"
	// "P" is displayed in cyan and "Term" is displayed in light magenta
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithRGB("T", pterm.NewRGB(61, 238, 253)),
		putils.LettersFromStringWithRGB("get", pterm.NewRGB(249, 46, 254))).WithWriter(os.Stderr).Render()
}

func CreateOutput(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return err
	} else {
		return os.MkdirAll(dir, 0o755)
	}
}
