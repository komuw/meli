package meli

/*
The functionality/code in this file is taken from: https://github.com/subosito/gotenv
Which is released by the author; Alif Rachmawadi(subosito) under MIT license.
*/

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	// Pattern for detecting valid line format
	linePattern = `\A\s*(?:export\s+)?([\w\.]+)(?:\s*=\s*|:\s+?)('(?:\'|[^'])*'|"(?:\"|[^"])*"|[^#\n]+)?\s*(?:\s*\#.*)?\z`

	// Pattern for detecting valid variable within a value
	variablePattern = `(\\)?(\$)(\{?([A-Z0-9_]+)?\}?)`
)

// env holds key/value pair of valid environment variable
// TODO: replace env with a []string since ComposeService.Environment is a []string
type env map[string]string

// parsedotenv is a function to parse line by line any io.Reader supplied and returns the valid Env key/value pair of valid variables.
// It expands the value of a variable from environment variable, but does not set the value to the environment itself.
// This function is skipping any invalid lines and only processing the valid one.
func parsedotenv(r io.Reader) env {
	e, _ := strictParse(r)
	return e
}

// StrictParse is a function to parse line by line any io.Reader supplied and returns the valid Env key/value pair of valid variables.
// It expands the value of a variable from environment variable, but does not set the value to the environment itself.
// This function is returning an error if there is any invalid lines.
func strictParse(r io.Reader) (env, error) {
	e := make(env)
	scanner := bufio.NewScanner(r)

	i := 1
	bom := string([]byte{239, 187, 191})

	for scanner.Scan() {
		line := scanner.Text()

		if i == 1 {
			line = strings.TrimPrefix(line, bom)
		}

		i++

		err := parseLine(line, e)
		if err != nil {
			return e, err
		}
	}

	return e, nil
}

func parseLine(s string, env env) error {
	rl := regexp.MustCompile(linePattern)
	rm := rl.FindStringSubmatch(s)

	if len(rm) == 0 {
		st := strings.TrimSpace(s)

		if (st == "") || strings.HasPrefix(st, "#") {
			return nil
		}

		if strings.HasPrefix(st, "export") {
			vs := strings.SplitN(st, " ", 2)

			if len(vs) > 1 {
				if _, ok := env[vs[1]]; !ok {
					return fmt.Errorf("line `%s` has an unset variable", st)

				}
			}
		}
		return fmt.Errorf("line `%s` doesn't match format", s)
	}

	key := rm[1]
	val := rm[2]

	// determine if string has quote prefix
	hdq := strings.HasPrefix(val, `"`)

	// determine if string has single quote prefix
	hsq := strings.HasPrefix(val, `'`)

	// trim whitespace
	val = strings.Trim(val, " ")

	// remove quotes '' or ""
	rq := regexp.MustCompile(`\A(['"])(.*)(['"])\z`)
	val = rq.ReplaceAllString(val, "$2")

	if hdq {
		val = strings.Replace(val, `\n`, "\n", -1)
		val = strings.Replace(val, `\r`, "\r", -1)

		// Unescape all characters except $ so variables can be escaped properly
		re := regexp.MustCompile(`\\([^$])`)
		val = re.ReplaceAllString(val, "$1")
	}

	rv := regexp.MustCompile(variablePattern)
	fv := func(s string) string {
		if strings.HasPrefix(s, "\\") {
			return strings.TrimPrefix(s, "\\")
		}

		if hsq {
			return s
		}

		sn := `(\$)(\{?([A-Z0-9_]+)\}?)`
		rn := regexp.MustCompile(sn)
		mn := rn.FindStringSubmatch(s)

		if len(mn) == 0 {
			return s
		}

		v := mn[3]

		replace, ok := env[v]
		if !ok {
			replace = os.Getenv(v)
		}

		return replace
	}

	val = rv.ReplaceAllStringFunc(val, fv)

	if strings.Contains(val, "=") {
		if !(val == "\n" || val == "\r") {
			kv := strings.Split(val, "\n")

			if len(kv) == 1 {
				kv = strings.Split(val, "\r")
			}

			if len(kv) > 1 {
				val = kv[0]

				for i := 1; i < len(kv); i++ {
					parseLine(kv[i], env)
				}
			}
		}
	}

	env[key] = val
	return nil
}
