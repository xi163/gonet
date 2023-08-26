package logs

import (
	"time"
)

func String(timezone Timezone) string {
	switch timezone {
	case MY_PST:
		return "PST"
	case MY_MST:
		return "MST"
	case MY_EST:
		return "EST"
	case MY_BST:
		return "BST"
	case MY_UTC:
		return "UTC"
	case MY_GST:
		return "GST"
	case MY_CST:
		return "CST"
	case MY_JST:
		return "JST"
	}
	return ""
}

func convertUTC(t *time.Time, tm *time.Time, timezone Timezone) bool {
	switch timezone {
	case MY_UTC:
		*tm = t.UTC() //UTC/GMT
	case MY_PST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_PST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	case MY_MST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_MST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	case MY_EST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_EST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	case MY_BST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_BST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	case MY_GST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_GST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	case MY_CST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_CST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	case MY_JST:
		tm_utc := t.UTC()
		t_zone_nanosec := tm_utc.UnixNano() + MY_JST_nano_sec
		*tm = time.Unix(0, t_zone_nanosec).UTC()
	default:
		return false
	}
	return true
}

func setting(tm *time.Time, timezone Timezone) {
	// loc, _ := time.LoadLocation("Asia/Shanghai")
	// tm_zone, _ := time.ParseInLocation("2006/01/02 15:04:05", "2018-07-11 15:07:51", loc)
	switch timezone {
	case MY_UTC:
		switch tm {
		case nil:
			FatalfF("error")
		}
	case MY_PST:
		switch tm {
		case nil:
			FatalfF("error")
		}
	case MY_MST:
		switch tm {
		case nil:
			FatalfF("error")
		}
	case MY_EST:
		switch tm {
		case nil:
			FatalfF("error")
		}
		InfofF("%v %v %v %v America/New_York %v", LevelString(), ModeString(), StyleString(), TimezoneString(), tm.Format("2006/01/02 15:04:05"))
	case MY_BST:
		switch tm {
		case nil:
			FatalfF("error")
		}
		InfofF("%v %v %v %v Europe/London %v", LevelString(), ModeString(), StyleString(), TimezoneString(), tm.Format("2006/01/02 15:04:05"))
	case MY_GST:
		switch tm {
		case nil:
			FatalfF("error")
		}
		InfofF("%v %v %v %v Asia/Dubai %v", LevelString(), ModeString(), StyleString(), TimezoneString(), tm.Format("2006/01/02 15:04:05"))
	case MY_CST:
		switch tm {
		case nil:
			FatalfF("error")
		}
		InfofF("%v %v %v %v Asia/Shanghai %v", LevelString(), ModeString(), StyleString(), TimezoneString(), tm.Format("2006/01/02 15:04:05"))
	case MY_JST:
		switch tm {
		case nil:
			FatalfF("error")
		}
		InfofF("%v %v %v %v Asia/Tokyo %v", LevelString(), ModeString(), StyleString(), TimezoneString(), tm.Format("2006/01/02 15:04:05"))
	default:
		ErrorfTLF("%v %v %v %v", LevelString(), ModeString(), StyleString(), TimezoneString())
	}
}
