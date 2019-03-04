package fittools /* import "ivo.qa/plfit-tablegen/fittools" */

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	// The naming pattern of all dat files. Files that don't adhere to it
	// shall be ignored.
	nameRegex = `^[\w\d\-]+_(?P<Frame>\d+)_(?P<Type>[\w\d\-]+)_` +
		`(?P<Period>\d+)_[\w\d\-]+_(?P<StructNo>\d+)\.dat$`
	noneFoundErr = "no files found, make sure you are using the correct " +
		"naming pattern: %q"
)

// Settings related to printing data as a table.
// See https://golang.org/pkg/text/tabwriter/#NewWriter
const (
	tableMinWidth = 2
	tabWidth      = 2
	padding       = 1
)

var (
	nameExpression = regexp.MustCompile(nameRegex)
)

// Data represents all known data about a file.
// It includes the information that can be obtained from the dat file
// (frame, type, period, and struct number), as well as alpha, xmin, and L,
// which are obtained using `plfit`.
type Data struct {
	Filename                string
	Type                    string
	Frame, Period, StructNo int
	Alpha, XMin, L          float64
}

func (d *Data) String() string {
	return fmt.Sprintf(`
File name: %s
Frame:     %d
Type:      %s
Period:    %d
Struct#:   %d
Alpha:     %f
XMin:      %f
L:         %f`,
		d.Filename, d.Frame, d.Type, d.Period,
		d.StructNo, d.Alpha, d.XMin, d.L)
}

// ExecutePlfit runs `plfit` and extracts data from its output, populating
// the alpha and xmin fields.
// In order for the call to succeed, make sure `plfit` is available via
// the system path.
// See http://tuvalu.santafe.edu/~aaronc/powerlaws/
func (d *Data) ExecutePlfit() error {
	out, err := exec.Command("plfit", d.Filename).Output()
	if err != nil {
		return err
	}
	output := *parsePlfitOutput(&out)
	d.Alpha, err = strconv.ParseFloat(output["alpha"], 32)
	if err != nil {
		return err
	}
	d.Alpha -= 1.0

	d.XMin, err = strconv.ParseFloat(output["xmin"], 32)
	if err != nil {
		return err
	}

	d.L, err = strconv.ParseFloat(output["L"], 32)
	if err != nil {
		return err
	}
	return nil
}

// parsePlfitOutput parses plfit's output into a map.
// Plfit's output consists of key-value pairs, separated by an equals sign,
// one per line.
// Lines that don't contain an equal sign have meta information and
// are ignored.
func parsePlfitOutput(output *[]byte) *map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(*output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return &result
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

// AsTable generates a table representation of the presented data and writes
// it to the specified file (usually os.Stdin).
func AsTable(f *os.File, data []*Data) {
	w := tabwriter.NewWriter(
		f, tableMinWidth, tabWidth, padding, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Frame\tH2/Ht\tPeriod\tStructure\tDP\tn\tL\t")
	for _, f := range data {
		fmt.Fprintf(
			w, "%d\t%s\t%d\t%d\t%f\t%f\t%f\t\n",
			f.Frame, f.Type, f.Period, f.StructNo, f.XMin, f.Alpha, f.L)
	}
	w.Flush()
}
