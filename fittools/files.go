package fittools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

const (
	// The naming pattern of all dat files. Files that don't adhere to it
	// shall be ignored.
	nameRegex = `^[\w\d\-]+_(?P<Frame>\d+)_(?P<Type>[\w\d\-]+)_` +
		`(?P<Period>\d+)_[\w\d\-]+_(?P<StructNo>\d+)\.dat$`
	noneFoundErr = "no files found, make sure you are using the correct " +
		"naming pattern: %q"
)

var (
	nameExpression = regexp.MustCompile(nameRegex)
)

// Data represents all known data about a file.
// It includes the information that can be obtained from the dat file
// (frame, type, period, and struct number), as well as alpha and xmin,
// which are obtained using `plfit`.
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

// ExecutePlfit runs `plfit` and extracts data from its output, populating
// the alpha and xmin fields.
// In order for the call to succeed, make sure `plfit` is available via
// the system path.
// See http://tuvalu.santafe.edu/~aaronc/powerlaws/
func (d *Data) ExecutePlfit() (output string, err error) {
	out, err := exec.Command("plfit", d.Filename).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

// fromRegexp extracts data (frame, type, period, struct #) from the name
// of a .dat file that matches the expression defined as `nameRegex`.
// The function assumes the name has already been validated against the
// expression and will panic in case there is no match.
// See collectData() as an example how to call.
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

// CollectData goes through the files at the location specified as source,
// looking for ones that match the expression defined as nameRegex.
// The function will not traverse through directories recursively.
func CollectData(source string) (data []*Data, err error) {
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

	if len(data) == 0 {
		err = errors.New(fmt.Sprintf(noneFoundErr, nameRegex))
		return data, err
	}

	return data, nil
}
