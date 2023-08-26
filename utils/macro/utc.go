package macro

import "time"

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
