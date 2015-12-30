package main

import (
  "flag"
  "fmt"
  "github.com/jzelinskie/geddit"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
)

type Bot struct {
  Username string
  Password string
}

func main() {

  // Using gist 'aggrolite/160d5be23adb9e597553', 'go-reddit-bot' as resources
  var bot Bot
  var config string
  var outfile string

  flag.StringVar(&config, "config", "bot.yml", "Path to YAML config file.")
  flag.StringVar(&outfile, "output", "links.txt", "Path to an output text file.")
  flag.Parse()

  // Read config file.
  config_data, err := ioutil.ReadFile(config)
  if err != nil {
    log.Fatal(err)
  }
  err = yaml.Unmarshal(config_data, &bot)
  if err != nil {
    log.Fatal(err)
  }

  // Use own reddit API wrapper once done.
  session, err := geddit.NewLoginSession(
    bot.Username,
    bot.Password,
    "Link Saver written in Go by u/Alexius-CA",
  )
  if err != nil {
    log.Fatal(err)
  }

  // Query options if needed (blank right now).
  //opts := geddit.ListingOptions{}

  ticker := time.NewTicker(30 * time.Second)
  quit := make(chan os.Signal)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

  // Acquire saved links and write to file every 1 hour.
  go func() {
    for {
      select {
      case <-ticker.C:
        saved, err := session.MySaved(geddit.NewSubmissions, "")
        write := ""
        for _, s := range saved {
          write += fmt.Sprintf("#%s\n%s\n", s.Subreddit, s.Permalink)
        }
        err = ioutil.WriteFile(outfile, []byte(write), 0644)
        if err != nil {
          log.Print(err)
        }
      case <-quit:
        ticker.Stop()
        return
      }
    }
  }()
  <-quit

}
