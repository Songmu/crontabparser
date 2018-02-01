package crontab

import (
	"fmt"
	"strings"
)

type Env struct {
	raw  string
	line int

	key string
	val string
}

func (env *Env) Type() Type {
	return TypeEnv
}

func (env *Env) Err() error {
	return nil
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
		env.key, err = dequote(kv[0])
		if err != nil {
			// XXX set error?
			return err
		}
		env.val, err = dequote(kv[1])
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("invalid env entry: %q", env.raw)
}

func dequote(stuff string) (string, error) {
	stuff = strings.TrimSpace(stuff)
	if strings.HasPrefix(stuff, `"`) && strings.HasSuffix(stuff, `"`) {
		stuff = strings.Trim(stuff, `"`)
	} else if strings.HasPrefix(stuff, `'`) && strings.HasSuffix(stuff, `'`) {
		stuff = strings.Trim(stuff, `'`)
	}
	return stuff, nil
}
