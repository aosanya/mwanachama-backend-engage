package engage

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

func zoneLoads(name string) error {
	if name == "" {
		return nil
	}
	if _, err := time.LoadLocation(name); err != nil {
		return fmt.Errorf("time_zone %q is not a time zone name such as Europe/Budapest", name)
	}
	return nil
}
