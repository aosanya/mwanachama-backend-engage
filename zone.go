package engage

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"
)

//go:embed zone_names.txt
var zoneNamesFile string

var zoneNames = namesIn(zoneNamesFile)

func namesIn(file string) map[string]bool {
	names := map[string]bool{}
	for _, name := range strings.Fields(file) {
		names[name] = true
	}
	return names
}

func zoneLoads(name string) error {
	if name == "" || name == "Local" {
		return nil
	}
	if _, err := time.LoadLocation(name); err != nil || !zoneNames[name] {
		return fmt.Errorf("time_zone %q is not a time zone name such as Europe/Budapest", name)
	}
	return nil
}
