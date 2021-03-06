package crontabparser

import (
	"fmt"
	"strings"
)

// Env cron entry
type Env struct {
	raw string

	ers errors
	key string
	val string
}

func newEnv(raw string) *Env {
	env := &Env{raw: raw}
	env.parse()
	return env
}

// Type returns TypeEnv
func (env *Env) Type() Type {
	return TypeEnv
}

// Err returns error if error
func (env *Env) Err() error {
	return env.ers.err()
}

// Raw returns raw contents of line
func (env *Env) Raw() string {
	return env.raw
}

// Key of env
func (env *Env) Key() string {
	return env.key
}

// Val of env
func (env *Env) Val() string {
	return env.val
}

func (env *Env) parse() {
	kv := strings.SplitN(env.raw, "=", 2)
	if len(kv) == 2 {
		k, err := dequote(kv[0])
		if err != nil {
			env.ers = append(env.ers, err)
		}
		v, err := dequote(kv[1])
		if err != nil {
			env.ers = append(env.ers, err)
		}
		if env.Err() == nil {
			env.key = k
			env.val = v
		}
		return
	}
	env.ers = append(env.ers, fmt.Errorf("invalid env entry: %q", env.raw))
}

func dequote(stuff string) (string, error) {
	ret := strings.TrimSpace(stuff)
	if strings.HasPrefix(ret, `"`) {
		if !strings.HasSuffix(ret, `"`) {
			return "", fmt.Errorf("invalid env element: %q", stuff)
		}
		ret = ret[1 : len(ret)-1]
		if strings.Contains(ret, `"`) {
			return "", fmt.Errorf("invalid env element: %q", stuff)
		}
	} else if strings.HasPrefix(ret, `'`) {
		if !strings.HasSuffix(ret, `'`) {
			return "", fmt.Errorf("invalid env element: %q", stuff)
		}
		ret = ret[1 : len(ret)-1]
		if strings.Contains(ret, `'`) {
			return "", fmt.Errorf("invalid env element: %q", stuff)
		}
	}
	return ret, nil
}
