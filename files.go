package main

import (
	"io/ioutil"
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
	Type                    string
	Frame, Period, StructNo int
	Alpha, XMin             float64
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

func collectData(source string) (data []Data, err error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return data, err
	}

	for _, f := range files {
		if nameExpression.MatchString(f.Name()) {
			data = append(data, *fromRegexp(f.Name()))
		}
	}
	return data, nil
}
