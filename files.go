package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

const (
	nameRegex = `^[\w\d\-]+_(?P<Frame>\d+)_(?P<Type>[\w\d\-]+)_` +
		`(?P<Period>\d+)_[\w\d\-]+_(?P<StructNo>\d+)\.dat$`
)

var (
	nameExpression = regexp.MustCompile(nameRegex)
)

type Data struct {
	Filename                string
	Type                    string
	Frame, Period, StructNo int
	Alpha, XMin             float64
}

func (d *Data) String() string {
	return fmt.Sprintf(`
File name: %s
Frame:     %d
Type:      %s
Period:    %d
Struct#:   %d`, d.Filename, d.Frame, d.Type, d.Period, d.StructNo)
}

func (d *Data) ExecutePlfit() (output string, err error) {
	out, err := exec.Command("plfit", d.Filename).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func fromRegexp(name string) *Data {
	datum := new(Data)
	result := nameExpression.FindStringSubmatch(name)

	// Ignoring errors in Atoi as matching the regex guarantees we've
	// found numbers.
	datum.Frame, _ = strconv.Atoi(result[1])
	datum.Type = result[2]
	datum.Period, _ = strconv.Atoi(result[3])
	datum.StructNo, _ = strconv.Atoi(result[4])

	return datum
}

func collectData(source string) (data []*Data, err error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return data, err
	}

	for _, f := range files {
		if nameExpression.MatchString(f.Name()) {
			file := fromRegexp(f.Name())
			file.Filename = path.Join(source, f.Name())
			data = append(data, file)
		}
	}
	return data, nil
}
