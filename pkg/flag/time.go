package flag

import "time"

func ParseTime(flag string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", flag)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
