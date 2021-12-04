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
	env.Filters["percent"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		v := stick.CoerceNumber(val)
		// Do some formatting.
		return fmt.Sprintf("%.1f", v) + "%"
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
		v = strings.TrimSuffix(v, ".000Z")
		return v
	}
	return env
}
