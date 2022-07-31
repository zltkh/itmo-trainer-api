package itmo_trainer_api

func minString(strings ...string) string {
	mn := strings[0]
	for _, v := range strings {
		if v < mn {
			mn = v
		}
	}
	return mn
}

func maxString(strings ...string) string {
	mx := strings[0]
	for _, v := range strings {
		if v > mx {
			mx = v
		}
	}
	return mx
}
