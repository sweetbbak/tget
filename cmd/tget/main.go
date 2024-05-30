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

	"github.com/anacrolix/torrent"
	"github.com/carlmjohnson/versioninfo"
	"github.com/jessevdk/go-flags"
	"github.com/pterm/pterm"
)

const (
	Version = "v0.1"
)

var opts struct {
	Magnet      string `short:"t" long:"torrent" description:"path to torrent or magnet link"`
	FileIndex   string `short:"f" long:"files" description:"select files to download by index ex: (0,2,3-5)"`
	Info        bool   `short:"i" long:"info" description:"display torrent info and exit"`
	FileInfo    bool   `short:"F" long:"list-files" description:"list files and their index of a given torrent"`
	Output      string `short:"o" long:"output" description:"path to a directory to output the torrent"`
	Proxy       string `short:"p" long:"proxy" description:"proxy URL to use"`
	Blocklist   string `short:"b" long:"blocklist" description:"path or URL pointing to a plain-text IP blocklist"`
	DisableIPV6 bool   `short:"4" long:"ipv4" description:"dont use ipv6"`
	Quiet       bool   `short:"q" long:"quiet" description:"dont output text or progress bar"`
	Version     bool   `short:"V" long:"version" description:"display the version and exit"`
	NoCleanup   bool   `short:"n" long:"no-cleanup" description:"dont delete torrent database files on exit"`
}

func Download(client *torrent.Client, tor string) error {
	t, err := AddTorrent(client, tor)
	if err != nil {
		return err
	}

	if opts.FileIndex != "" {
		idxs, err := ParseIndices(opts.FileIndex)
		if err != nil {
			return err
		}

		if err := SelectFileByIndex(t, client, idxs); err != nil {
			return err
		}
	} else {
		t.DownloadAll()
		if !opts.Quiet {
			go func() {
				Progress(t)
			}()
		}
	}

	if t.Complete.Bool() {
		t.VerifyData()
	}

	if client.WaitAll() {
		pterm.Success.Printf("Downloaded: %s\n", t.Name())
		// client.WriteStatus(os.Stdout)
		SeedProgress(t)
		return nil
	} else {
		return fmt.Errorf("Unable to completely download torrent: %s", t.Name())
	}
}

func ConfigClient() (*torrent.Client, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DisableIPv6 = opts.DisableIPV6

	CreateOutput(opts.Output)
	tmpdir := os.TempDir()
	// cfg.DefaultStorage = storage.NewFile(opts.Output)

	// set default output directory
	stor, err := getMetadataDir(tmpdir, opts.Output)
	if err != nil {
		return nil, err
	}

	cfg.DefaultStorage = stor

	if opts.Blocklist != "" {
		f, err := openBlocklist(opts.Blocklist)
		if err != nil {
			return nil, err
		}
		ipb, err := parseBlocklist(f)
		cfg.IPBlocklist = ipb
	}

	if opts.Proxy != "" {
		u, err := url.Parse(opts.Proxy)
		if err != nil {
			return nil, err
		}
		cfg.HTTPProxy = http.ProxyURL(u)
	}

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func AddTorrent(client *torrent.Client, tor string) (*torrent.Torrent, error) {
	var t *torrent.Torrent
	var err error
	if strings.Contains(tor, "magnet") {
		t, err = client.AddMagnet(tor)
		if err != nil {
			return nil, err
		}
	} else if strings.Contains(tor, "http") {
		success, _ := pterm.DefaultSpinner.Start("Downloading torrent from remote...")
		path, err := fileFromURL(tor)
		if err != nil {
			success.Fail("Unable to download remote torrent")
			return nil, fmt.Errorf("tget: unable to download torrent from URL [%s]: %v", opts.Magnet, err)
		}
		success.Success("Downloaded torrent file")
		t, err = client.AddTorrentFromFile(path)
		if err != nil {
			return nil, err
		}
	} else {
		t, err = client.AddTorrentFromFile(tor)
		if err != nil {
			return nil, err
		}
	}

	success, _ := pterm.DefaultSpinner.Start("Getting torrent info...")
	<-t.GotInfo()
	success.Success("Got torrent info!")
	return t, nil
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
		// Cleanup()
		os.Exit(0)
	}()
}

func showVersion() {
	fmt.Printf("tget %s (%s)\n", Version, versioninfo.Revision)
	os.Exit(0)
}

func main() {
	args, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	if opts.Version {
		showVersion()
	}

	if opts.Magnet == "" && len(args) < 1 {
		log.Fatal("tget: must provide torrent file")
	}

	if opts.Magnet == "" && len(args) > 0 {
		opts.Magnet = args[0]
	}

	client, err := ConfigClient()
	if err != nil {
		log.Fatal(err)
	}

	if opts.FileInfo {
		if err := Info(opts.Magnet, client); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if opts.Info {
		if err := Info(opts.Magnet, client); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	HandleExit()

	if !opts.Quiet {
		Header()
	}

	if err := Download(client, opts.Magnet); err != nil {
		log.Fatal(err)
	}
}
