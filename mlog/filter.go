package mlog

type Filter func(*Record) bool

func LevelFilter(level ...string) Filter {
	var mask Level

	if len(level) == 0 {
		mask = LevelAll
	} else {
		for _, v := range level {
			mask |= NewLevel(v)
		}
	}

	return func(r *Record) (e bool) {
		e = r.level&mask > 0
		return
	}
}
