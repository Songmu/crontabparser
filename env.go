package checron

import (
	"fmt"
	"strings"
)

type Env struct {
	raw  string
	line int

	ers errors
	key string
	val string
}

func (env *Env) Type() Type {
	return TypeEnv
}

func (env *Env) Err() error {
	return env.ers.err()
}

func (env *Env) Raw() string {
	return env.raw
}

func (env *Env) Line() int {
	return env.line
}

func (env *Env) Key() string {
	if env.key == "" {
		env.parse()
	}
	return env.key
}

func (env *Env) Val() string {
	if env.val == "" {
		env.parse()
	}
	return env.val
}

func (env *Env) parse() (err error) {
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
		return env.Err()
	}
	env.ers = append(env.ers, fmt.Errorf("invalid env entry: %q", env.raw))
	return env.Err()
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
