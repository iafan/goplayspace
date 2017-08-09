package draw

import (
	"regexp"
	"strconv"
	"strings"
)

const cmdStartDrawMode = "draw mode"

var cmdForwardR = regexp.MustCompile(`^forward$`)
var cmdForwardNR = regexp.MustCompile(`^forward (\d+(\.\d+)?)$`)
var cmdLeftR = regexp.MustCompile(`^left$`)
var cmdLeftNR = regexp.MustCompile(`^left (\d+(\.\d+)?)$`)
var cmdRightR = regexp.MustCompile(`^right$`)
var cmdRightNR = regexp.MustCompile(`^right (\d+(\.\d+)?)$`)
var cmdColorOffR = regexp.MustCompile(`^(?:color|colour) off$`)
var cmdColorSR = regexp.MustCompile(`^(?:color|colour) (.+)$`)
var cmdWidthNR = regexp.MustCompile(`^width (\d+(\.\d+)?)$`)
var cmdSaySR = regexp.MustCompile(`^say (.+)$`)

const (
	Step = iota
	Left
	Right
	Color
	Width
	Say
)

type Action struct {
	Cmd  string
	Kind int
	FVal float64
	SVal string
}

type ActionList []*Action

func ParseLines(lines []string) ActionList {
	a := make(ActionList, 0)

	isDrawMode := false

	for _, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))

		if !isDrawMode && line != cmdStartDrawMode {
			continue
		}
		isDrawMode = true

		if matches := cmdForwardR.FindAllStringSubmatch(line, -1); matches != nil {
			a = append(a, &Action{line, Step, 1, ""})
			continue
		}

		if matches := cmdForwardNR.FindAllStringSubmatch(line, -1); matches != nil {
			n, _ := strconv.ParseFloat(matches[0][1], 64)
			a = append(a, &Action{line, Step, n, ""})
			continue
		}

		if matches := cmdLeftR.FindAllStringSubmatch(line, -1); matches != nil {
			a = append(a, &Action{line, Left, 90, ""})
			continue
		}

		if matches := cmdLeftNR.FindAllStringSubmatch(line, -1); matches != nil {
			n, _ := strconv.ParseFloat(matches[0][1], 64)
			a = append(a, &Action{line, Left, n, ""})
			continue
		}

		if matches := cmdRightR.FindAllStringSubmatch(line, -1); matches != nil {
			a = append(a, &Action{line, Right, 90, ""})
			continue
		}

		if matches := cmdRightNR.FindAllStringSubmatch(line, -1); matches != nil {
			n, _ := strconv.ParseFloat(matches[0][1], 64)
			a = append(a, &Action{line, Right, n, ""})
			continue
		}

		if matches := cmdColorOffR.FindAllStringSubmatch(line, -1); matches != nil {
			a = append(a, &Action{line, Color, 0, ""})
			continue
		}

		if matches := cmdColorSR.FindAllStringSubmatch(line, -1); matches != nil {
			a = append(a, &Action{line, Color, 0, matches[0][1]})
			continue
		}

		if matches := cmdWidthNR.FindAllStringSubmatch(line, -1); matches != nil {
			n, _ := strconv.ParseFloat(matches[0][1], 64)
			a = append(a, &Action{line, Width, n, ""})
			continue
		}

		if matches := cmdSaySR.FindAllStringSubmatch(line, -1); matches != nil {
			a = append(a, &Action{line, Say, 0, matches[0][1]})
			continue
		}

	}

	return a
}

func ParseString(s string) ActionList {
	return ParseLines(strings.Split(s, "\n"))
}
