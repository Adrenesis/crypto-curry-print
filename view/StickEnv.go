package view

import (
	"fmt"
	"github.com/tyler-sommer/stick"
	"os"
	"strings"
)

func GetEnv() *stick.Env {
	fsRoot, _ := os.Getwd() // Templates are loaded relative to this directory.
	fsRoot = fmt.Sprintf(fsRoot + "/templates")
	env := stick.New(stick.NewFilesystemLoader(fsRoot))
	env.Filters["sprintf"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		// Do some formatting.
		return fmt.Sprintf("%+v", val)
	}
	env.Filters["sprintfjson"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		// Do some formatting.
		v := fmt.Sprintf("%+v", val)
		v = strings.Replace(v, " ", ", \"", -1)
		v = strings.Replace(v, ":", "\":", -1)
		return fmt.Sprintf("%+v", v)
	}
	env.Filters["percent"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.1f", v) + "%"
	}
	env.Filters["f4d"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.4f", v)
	}
	env.Filters["number_format"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.10f", v)
	}
	env.Filters["number_format_vol"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.2f", v)
	}
	env.Filters["date"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceString(val)
		// Do some formatting.
		v = fmt.Sprintf("%s", v)
		if len(v) > 19 {
			v = v[:19]
		}

		return v
	}
	return env
}
