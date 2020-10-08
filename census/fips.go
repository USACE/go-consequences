package census

func StateToCountyFipsMap() map[string][]string {
	m := make(map[string][]string)
	m["01"] = []string{"01001", "01003"}
	m["02"] = []string{"02001", "02003"}
	return m
}
