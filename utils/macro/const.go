package macro

type Timezone int8

const (
	MY_PST = Timezone(-8)
	MY_MST = Timezone(-7)
	MY_EST = Timezone(-5)
	MY_BST = Timezone(+1)
	//UTC/GMT
	MY_UTC = Timezone(+0)
	//(UTC+04:00) Asia/Dubai
	MY_GST = Timezone(+4)
	//(UTC+08:00) Asia/shanghai, Beijing(China)
	MY_CST = Timezone(+8)
	MY_JST = Timezone(+9)

	MY_PST_nano_sec = 3600 * int64(MY_PST) * 1e9
	MY_MST_nano_sec = 3600 * int64(MY_MST) * 1e9
	MY_EST_nano_sec = 3600 * int64(MY_EST) * 1e9
	MY_BST_nano_sec = 3600 * int64(MY_BST) * 1e9
	MY_UTC_nano_sec = 3600 * int64(MY_UTC) * 1e9
	MY_GST_nano_sec = 3600 * int64(MY_GST) * 1e9
	MY_CST_nano_sec = 3600 * int64(MY_CST) * 1e9
	MY_JST_nano_sec = 3600 * int64(MY_JST) * 1e9
)

type Level uint8

const (
	LVL_FATAL Level = 0
	LVL_ERROR Level = 1
	LVL_WARN  Level = 2
	LVL_INFO  Level = 3
	LVL_TRACE Level = 4
	LVL_DEBUG Level = 5
)

type Style uint16

const (
	F_SYNC              Style = 0x1000
	F_DETAIL            Style = 0x0001
	F_TMSTMP            Style = 0x0002
	F_FN                Style = 0x0004
	F_TMSTMP_FN         Style = 0x0008
	F_FL                Style = 0x0010
	F_TMSTMP_FL         Style = 0x0020
	F_FL_FN             Style = 0x0040
	F_TMSTMP_FL_FN      Style = 0x0080
	F_TEXT              Style = 0x0100
	F_PURE              Style = 0x0200
	F_DETAIL_SYNC             = F_DETAIL | F_SYNC
	F_TMSTMP_SYNC             = F_TMSTMP | F_SYNC
	F_FN_SYNC                 = F_FN | F_SYNC
	F_TMSTMP_FN_SYNC          = F_TMSTMP_FN | F_SYNC
	F_FL_SYNC                 = F_FL | F_SYNC
	F_TMSTMP_FL_SYNC          = F_TMSTMP_FL | F_SYNC
	F_FL_FN_SYNC              = F_FL_FN | F_SYNC
	F_TMSTMP_FL_FN_SYNC       = F_TMSTMP_FL_FN | F_SYNC
	F_TEXT_SYNC               = F_TEXT | F_SYNC
	F_PURE_SYNC               = F_PURE | F_SYNC
)

type Mode uint8

const (
	M_STDOUT_ONLY Mode = iota
	M_FILE_ONLY
	M_STDOUT_FILE
)
