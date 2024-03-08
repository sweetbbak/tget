package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/jessevdk/go-flags"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

var opts struct {
	Magnet      string `short:"t" long:"torrent" description:"path to torrent or magnet link"`
	Output      string `short:"o" long:"output" description:"path to a directory to output the torrent"`
	Proxy       string `short:"p" long:"proxy" description:"proxy URL to use"`
	Blocklist   string `short:"b" long:"blocklist" description:"path or URL pointing to a plain-text IP blocklist"`
	DisableIPV6 bool   `short:"4" long:"ipv4" description:"dont use ipv6"`
	Quiet       bool   `short:"q" long:"quiet" description:"dont output text or progress bar"`
	NoCleanup   bool   `short:"n" long:"no-cleanup" description:"dont delete torrent database files on exit"`
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

func Header() {
	// Initialize a big text display with the letters "P" and "Term"
	// "P" is displayed in cyan and "Term" is displayed in light magenta
	pterm.DefaultBigText.WithLetters(
		// putils.LettersFromStringWithStyle("T", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithRGB("T", pterm.NewRGB(61, 238, 253)),
		putils.LettersFromStringWithRGB("get", pterm.NewRGB(249, 46, 254))).
		Render() // Render the big text to the terminal
}

func Progress(t *torrent.Torrent) {
	title := TruncateString(t.Name(), 100)
	fmt.Printf("name [%s]...", title)

	p, _ := pterm.DefaultProgressbar.WithTotal(100).Start()

	for !t.Complete.Bool() {
		pc := float64(t.BytesCompleted()) / float64(t.Length()) * 100
		numpeers := len(t.PeerConns())
		p.Increment().Current = int(pc)
		p.UpdateTitle(fmt.Sprintf("peers [%v]", numpeers))
		time.Sleep(time.Millisecond * 50)
	}
}

func CreateOutput(dir string) {
	_, err := os.Stat(opts.Output)
	if err == nil {
		return
	} else {
		os.MkdirAll(opts.Output, 0o755)
	}
}

func Download() error {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DisableIPv6 = opts.DisableIPV6

	CreateOutput(opts.Output)
	cfg.DefaultStorage = storage.NewFile(opts.Output)

	if opts.Blocklist != "" {
		f, err := openBlocklist(opts.Blocklist)
		if err != nil {
			return err
		}
		ipb, err := parseBlocklist(f)
		cfg.IPBlocklist = ipb
	}

	if opts.Proxy != "" {
		u, err := url.Parse(opts.Proxy)
		if err != nil {
			return err
		}
		cfg.HTTPProxy = http.ProxyURL(u)
	}

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}

	var t *torrent.Torrent
	if strings.Contains(opts.Magnet, "magnet") {
		t, err = client.AddMagnet(opts.Magnet)
		if err != nil {
			return err
		}
	} else if strings.Contains(opts.Magnet, "http") {
		success, _ := pterm.DefaultSpinner.Start("Downloading torrent from remote...")
		path, err := fileFromURL(opts.Magnet)
		if err != nil {
			success.Fail("Unable to download remote torrent")
			return fmt.Errorf("tget: unable to download torrent from URL [%s]: %v", opts.Magnet, err)
		}
		success.Success("Downloaded torrent file")
		t, err = client.AddTorrentFromFile(path)
		if err != nil {
			return err
		}
	} else {
		t, err = client.AddTorrentFromFile(opts.Magnet)
		if err != nil {
			return err
		}
	}

	success, _ := pterm.DefaultSpinner.Start("Getting torrent info...")
	<-t.GotInfo()
	success.Success("Got torrent info!")

	if !opts.Quiet {
		go func() {
			Progress(t)
		}()
	}

	t.DownloadAll()
	if client.WaitAll() {
		pterm.Success.Printf("Downloaded: %s\n", t.Name())
		return nil
	} else {
		return fmt.Errorf("Unable to completely download torrent: %s", t.Name())
	}
}

func Cleanup() error {
	files := []string{".torrent.db", ".torrent.db-shm", ".torrent.db-wal", ".torrent.bolt.db"}
	for _, f := range files {
		fp := path.Join(opts.Output, f)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			continue
		}
		err := os.Remove(fp)
		if err != nil {
			return fmt.Errorf("unable to remove torrent db: %v", err)
		}
	}
	return nil
}

func init() {
	if opts.Output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			opts.Output = "."
		} else {
			opts.Output = cwd
		}
	}
}

func HandleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println()
		fmt.Printf("\n\n\n%s\n", "Are you sure you want to exit? press [ctrl+c]")
		<-c
		fmt.Printf("\n\n\n")
		pterm.Info.Println("Exiting...")
		Cleanup()
		os.Exit(0)
	}()
}

func main() {
	args, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	if !opts.Quiet {
		Header()
	}

	if opts.Magnet == "" && len(args) < 1 {
		log.Fatal("tget: must provide torrent file")
	}

	if opts.Magnet == "" && len(args) > 0 {
		opts.Magnet = args[0]
	}

	HandleExit()

	if err := Download(); err != nil {
		log.Fatal(err)
	}

	if !opts.NoCleanup {
		if err := Cleanup(); err != nil {
			log.Fatal(err)
		}
	}
}
