package testutil

import "time"

type Duration time.Duration

func (d Duration) StdDuration() time.Duration {
	return time.Duration(d)
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalText(txt []byte) error {
	v, err := time.ParseDuration(string(txt))
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}
