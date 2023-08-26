package macro

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
