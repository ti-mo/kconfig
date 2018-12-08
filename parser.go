package kconfig

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parse(s *bufio.Scanner, p map[string]string) error {

	r, _ := regexp.Compile("^(?:# *)?CONFIG_(\\w*)(?:=| )(y|n|m|is not set|\\d+|0x.+|\".*\")$")

	for s.Scan() {

		t := s.Text()

		// Skip line if empty.
		if t == "" {
			continue
		}

		// 0 is the match of the entire expression,
		// 1 is the key, 2 is the value.
		m := r.FindStringSubmatch(t)
		if m == nil {
			continue
		}

		if len(m) != 3 {
			return fmt.Errorf("match is not 3 chars long: %v", m)
		}

		// Remove all leading and trailing double quotes from the value.
		if len(m[2]) > 1 {
			m[2] = strings.Trim(m[2], "\"")
		}

		// Insert entry into map.
		p[m[1]] = m[2]
	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

func dump(p map[string]string, w *bufio.Writer) error {

	var err error

	for k, v := range p {

		if v == "y" || v == "n" || v == "m" {
			// No quotes needed around tri-state.
			_, err = w.WriteString(fmt.Sprintf("CONFIG_%s=%s\n", k, v))
		} else if v == "is not set" {
			// Value 'is not set'.
			_, err = w.WriteString(fmt.Sprintf("# CONFIG_%s is not set\n", k))
		} else if strings.HasPrefix(v, "0x") {
			// Value is a hex number, no quotes needed.
			// Cheap test, no need to do ParseInt here.
			_, err = w.WriteString(fmt.Sprintf("CONFIG_%s=%s\n", k, v))
		} else if _, err := strconv.Atoi(v); err == nil {
			// Value is a decimal, no quotes needed.
			_, err = w.WriteString(fmt.Sprintf("CONFIG_%s=%s\n", k, v))
		} else {
			// Type of value unrecognized or empty, quote it.
			_, err = w.WriteString(fmt.Sprintf("CONFIG_%s=\"%s\"\n", k, v))
		}

		if err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}
