package graphpaper

import ()

type Configuration struct {
  Resolutions []ResolutionConfig
}

type ResolutionConfig struct {
  Name       string
  Resolution int64
  Size       int64
  DateFmt    string
}

var Config = Configuration{[]ResolutionConfig{{"5m", 5 * 60 * 1000000000, 24 * 60 * 60 * 1000000000, "2006-01-02-15-04"}, {"1h", 60 * 60 * 1000000000, 24 * 60 * 60 * 1000000000, "2006-01-02"}}}
