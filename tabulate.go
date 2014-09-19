package main

import "fmt"
import "bytes"

// Basic Structure of TableFormat
type TableFormat struct {
	LineTop         Line
	LineBelowHeader Line
	LineBetweenRows Line
	LineBottom      Line
	HeaderRow       Row
	DataRow         Row
	Padding         int
	HeaderHide      bool
	FitScreen       bool
}

// Represents a Line
type Line struct {
	begin string
	hline string
	sep   string
	end   string
}

// Represents a Row
type Row struct {
	begin string
	sep   string
	end   string
}

// Table Formats that are available to the user
// The user can define his own format, just by addind an entry to this map
// and calling it with Render function e.g t.Render("customFormat")
var TableFormats = map[string]TableFormat{
	"simple": TableFormat{
		LineTop:         Line{"", "-", "  ", ""},
		LineBelowHeader: Line{"", "-", "  ", ""},
		LineBottom:      Line{"", "-", "  ", ""},
		HeaderRow:       Row{"", "  ", ""},
		DataRow:         Row{"", "  ", ""},
		Padding:         1,
	},
	"plain": TableFormat{
		HeaderRow: Row{"", "  ", ""},
		DataRow:   Row{"", "  ", ""},
		Padding:   1,
	},
	"grid": TableFormat{
		LineTop:         Line{"+", "-", "+", "+"},
		LineBelowHeader: Line{"+", "=", "+", "+"},
		LineBetweenRows: Line{"+", "-", "+", "+"},
		LineBottom:      Line{"+", "-", "+", "+"},
		HeaderRow:       Row{"|", "|", "|"},
		DataRow:         Row{"|", "|", "|"},
		Padding:         1,
	},
}

// Main Tabulate structure
type Tabulate struct {
	Data        []*TabulateRow
	Headers     []string
	FormatFloat byte
	TableFormat TableFormat
	Align       string
	EmptyVar    string
	HideLines   []string
}

// Represents normalized tabulate Row
type TabulateRow struct {
	Elements []string
}

func (t *Tabulate) normalize() {

}

func (t *Tabulate) drawLine() {

}

func (t *Tabulate) padLine() {

}

func (t *Tabulate) format() {

}

func (t *Tabulate) alignCol(elements []string, align string, minwidth int) {

}

func (t *Tabulate) padRow(arr []string, padding int) []string {
	if len(arr) < 1 {
		return arr
	}
	padded := make([]string, len(arr))
	for index, el := range arr {
		var buffer bytes.Buffer
		// Pad left
		for i := 0; i < padding; i++ {
			buffer.WriteString(" ")
		}

		buffer.WriteString(el)

		// Pad Right
		for i := 0; i < padding; i++ {
			buffer.WriteString(" ")
		}

		padded[index] = buffer.String()
	}
	return padded
}

func (t *Tabulate) padLeft(width int, str string) string {
	var buffer bytes.Buffer
	// Pad left
	padding := width - len(str)
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}

	buffer.WriteString(str)
	return buffer.String()
}

func (t *Tabulate) padRight(width int, str string) string {
	var buffer bytes.Buffer
	padding := width - len(str)

	buffer.WriteString(str)

	// Add Padding right
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}
	return buffer.String()
}

func (t *Tabulate) padCenter(width int, str string) string {
	var buffer bytes.Buffer
	padding := int((width - len(str)) / 2)

	// Add padding left
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}
	// Write string
	buffer.WriteString(str)

	// Add padding right
	for i := 0; i < padding; i++ {
		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (t *Tabulate) buildLine(padded_widths []int, padding []int, l Line) string {
	cells := make([]string, len(padded_widths))

	for i, _ := range cells {
		var buffer bytes.Buffer
		for j := -1; j <= padding[i]; j++ {
			buffer.WriteString(l.hline)
		}
		cells[i] = buffer.String()
	}
	var buffer bytes.Buffer

	// Print begin
	buffer.WriteString(l.begin)

	// Print contents
	for i := 0; i < len(cells); i++ {
		if i != len(cells)-1 {
			buffer.WriteString(cells[i] + l.sep)
		} else {
			buffer.WriteString(cells[i])
		}
	}

	// Print end
	buffer.WriteString(l.end)
	return buffer.String()
}

func (t *Tabulate) buildRow(elements []string, padded_widths []int, paddings []int, d Row) string {

	var buffer bytes.Buffer
	buffer.WriteString(d.begin)
	padFunc := t.getAlignFunc()
	// Print contents
	for i := 0; i < len(padded_widths); i++ {
		output := padFunc(padded_widths[i], t.EmptyVar)
		if len(elements) > i {
			output = padFunc(padded_widths[i], elements[i])
		}
		buffer.WriteString(output)
		if i != len(padded_widths)-1 {
			buffer.WriteString(d.sep)
		}
	}
	// Print end
	buffer.WriteString(d.end)

	return buffer.String()
}

func (t *Tabulate) Render(format ...interface{}) string {
	var lines []string

	// If headers are set use them, otherwise use the first row
	if len(t.Headers) < 1 {
		t.Headers = t.Data[0].Elements
	}

	// Use the format that was passed as parameter, otherwise
	// use the format defined in the struct
	if len(format) > 0 {
		t.TableFormat = TableFormats[format[0].(string)]
	}

	// Check if Data is present
	if len(t.Data) < 1 {
		panic("No Data specified")
	}

	// Get Min widths for columns, based on headers
	min_widths := make([]int, len(t.Headers))
	for index, item := range t.Headers {
		min_widths[index] = len(item)
	}

	// Get Column widths for all columns
	cols := t.getWidths(t.Headers, t.Data)

	padded_widths := make([]int, len(cols))
	for i, _ := range padded_widths {
		padded_widths[i] = cols[i] + 2*t.TableFormat.Padding
	}

	// Start appending lines

	// Append top line if not hidden
	if !inSlice("top", t.HideLines) {
		lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineTop))
	}

	// Add Header
	lines = append(lines, t.buildRow(t.padRow(t.Headers, t.TableFormat.Padding), padded_widths, cols, t.TableFormat.HeaderRow))

	// Add Line Below Header if not hidden
	if !inSlice("belowheader", t.HideLines) {
		lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineBelowHeader))
	}

	// Add Data Rows
	for index, element := range t.Data {
		lines = append(lines, t.buildRow(t.padRow(element.Elements, t.TableFormat.Padding), padded_widths, cols, t.TableFormat.DataRow))
		if index < len(t.Data)-1 {
			lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineBetweenRows))
		}
	}

	if !inSlice("bottomLine", t.HideLines) {
		lines = append(lines, t.buildLine(padded_widths, cols, t.TableFormat.LineBottom))
	}

	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.WriteString(line + "\n")
	}

	return buffer.String()
}

func (t *Tabulate) getWidths(headers []string, data []*TabulateRow) []int {
	widths := make([]int, len(headers))
	current_max := 0
	for i := 0; i < len(headers); i++ {
		current_max = len(headers[i])
		for _, item := range data {
			if len(item.Elements) > i && len(widths) > i {
				element := item.Elements[i]
				if len(element) > current_max {
					widths[i] = len(element)
					current_max = len(element)
				} else {
					widths[i] = current_max
				}
			}
		}
	}

	return widths
}

func (t *Tabulate) SetHeaders(headers []string) *Tabulate {
	t.Headers = headers
	return t
}

func (t *Tabulate) SetColWidth(width int) {

}

func (t *Tabulate) SetTableFormat(format string) {

}

func (t *Tabulate) SetFormatting(format string) *Tabulate {
	return t
}

func (t *Tabulate) SetAlign(align string) {
	t.Align = align
}

func (t *Tabulate) getAlignFunc() func(int, string) string {
	if len(t.Align) < 1 || t.Align == "right" {
		return t.padLeft
	} else if t.Align == "left" {
		return t.padRight
	} else {
		return t.padCenter
	}
}

func (t *Tabulate) SetEmptyString(empty string) {
	t.EmptyVar = empty
}

func Create(data interface{}) *Tabulate {
	t := &Tabulate{FormatFloat: 'f'}

	switch v := data.(type) {
	case [][]string:
		t.Data = createFromString(data.([][]string))
	case [][]int32:
		t.Data = createFromInt32(data.([][]int32))
	case [][]int64:
		t.Data = createFromInt64(data.([][]int64))
	case [][]int:
		t.Data = createFromInt(data.([][]int))
	case [][]bool:
		t.Data = createFromBool(data.([][]bool))
	case [][]float64:
		t.Data = createFromFloat64(data.([][]float64), t.FormatFloat)
	case [][]interface{}:
		t.Data = createFromMixed(data.([][]interface{}), t.FormatFloat)
	case []string:
		t.Data = createFromString([][]string{data.([]string)})
	case []interface{}:
		t.Data = createFromMixed([][]interface{}{data.([]interface{})}, t.FormatFloat)
	case map[string][]interface{}:
		t.Headers, t.Data = createFromMapMixed(data.(map[string][]interface{}), t.FormatFloat)
	case map[string][]string:
		t.Headers, t.Data = createFromMapString(data.(map[string][]string))
	default:
		fmt.Println(v)
	}

	return t
}

func main() {
	row1 := []interface{}{"test_row1_1", "test_row1_2", "test_row1_3", 1, 2, 3, 11.10, '1'}
	row2 := []interface{}{"test_row2_1", "testssss_row2_222222", "test_row2_3", 1, 2, 3}
	row3 := []interface{}{"test_row3_2", "test_row3_3", 1}
	t := Create([][]interface{}{row1, row2, row3})
	t.SetHeaders([]string{"Test", "Test Hea222der 2", "Test Header 3", "T 4", "Test Header 5"})
	t.SetEmptyString("None")
	t.SetAlign("right")
	fmt.Println(t.Render("grid"))

	// Test Map
	maptest := map[string][]interface{}{"test": row1, "test222": row2, "test header22": row3}
	maptest_tabulate := Create(maptest)
	maptest_tabulate.SetEmptyString("None")
	maptest_tabulate.SetAlign("right")
	fmt.Println(maptest_tabulate.Render("grid"))

}
