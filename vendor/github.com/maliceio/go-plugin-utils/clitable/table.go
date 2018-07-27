// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <http://unlicense.org>
//

// Package clitable implements methods for pretty command line table output.
package clitable

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Table - Table structure.
type Table struct {
	Fields     []string
	Footer     map[string]string
	Rows       []map[string]string
	HideHead   bool // when true doesn't print header
	Markdown   bool
	fieldSizes map[string]int
}

// New - Creates a new table.
func New(fields []string) *Table {
	return &Table{
		Fields:     fields,
		Rows:       make([]map[string]string, 0),
		fieldSizes: make(map[string]int),
	}
}

// PrintTable - Prints table.
func PrintTable(fields []string, rows []map[string]interface{}) {
	table := New(fields)
	for _, r := range rows {
		table.AddRow(r)
	}
	table.Print()
}

// PrintHorizontal - Prints horizontal table from a map.
func PrintHorizontal(m map[string]interface{}) {
	table := New([]string{"Key", "Value"})
	rows := mapToRows(m)
	for _, row := range rows {
		table.AddRow(row)
	}
	table.HideHead = true
	table.Print()
}

// PrintRow - Prints table with only one row.
func PrintRow(fields []string, row map[string]interface{}) {
	table := New(fields)
	table.AddRow(row)
	table.Print()
}

// AddRow - Adds row to the table.
func (t *Table) AddRow(row map[string]interface{}) {
	newRow := make(map[string]string)
	for _, k := range t.Fields {
		v := row[k]
		// If is not nil format
		// else value is empty string
		var val string
		if v == nil {
			val = ""
		} else {
			val = fmt.Sprintf("%v", v)
		}

		newRow[k] = val
	}

	t.calculateSizes(newRow)

	if len(newRow) > 0 {
		t.Rows = append(t.Rows, newRow)
	}
}

// AddFooter - Adds footer to the table.
func (t *Table) AddFooter(footer map[string]string) {
	t.Footer = footer
}

// Print - Prints table.
func (t *Table) Print() {
	if len(t.Rows) == 0 && t.Footer == nil {
		return
	}

	t.calculateSizes(t.Footer)

	if !t.Markdown {
		t.printDash()
	}

	if !t.HideHead {
		fmt.Println(t.getHead())
		t.printTableDash()
	}

	for _, r := range t.Rows {
		fmt.Println(t.rowString(r))
		if !t.Markdown {
			t.printDash()
		}
	}

	if t.Footer != nil {
		t.printTableDash()
		fmt.Println(t.rowString(t.Footer))
		if !t.Markdown {
			t.printTableDash()
		}
	}
}

// String - Ouput table as a string.
func (t *Table) String(title string) string {

	tableString := []string{}

	if len(t.Rows) == 0 && t.Footer == nil {
		return ""
	}

	tableString = append(tableString, "### "+title)

	t.calculateSizes(t.Footer)

	if !t.Markdown {
		// t.printDash()
		tableString = append(tableString, t.stringDash())
	}

	if !t.HideHead {
		tableString = append(tableString, t.getHead())
		tableString = append(tableString, t.stringTableDash())
		// fmt.Println(t.getHead())
		// t.printTableDash()
	}

	for _, r := range t.Rows {
		// fmt.Println(t.rowString(r))
		tableString = append(tableString, t.rowString(r))
		if !t.Markdown {
			tableString = append(tableString, t.stringDash())
		}
	}

	if t.Footer != nil {
		tableString = append(tableString, t.stringTableDash())
		tableString = append(tableString, t.rowString(t.Footer))
		if !t.Markdown {
			tableString = append(tableString, t.stringDash())
		}
	}

	return strings.Join(tableString[:], "\n")
}

// getHead - Returns table header containing fields names.
func (t *Table) getHead() string {
	s := "|"
	for _, name := range t.Fields {
		s += t.fieldString(name, strings.Title(name)) + "|"
	}
	return s
}

// rowString - Creates a string row.
func (t *Table) rowString(row map[string]string) string {
	s := "|"
	for _, name := range t.Fields {
		value := row[name]
		s += t.fieldString(name, value) + "|"
	}
	return s
}

// fieldString - Creates field value string.
func (t *Table) fieldString(name, value string) string {
	value = fmt.Sprintf(" %s ", value)
	spacesLeft := t.fieldSizes[name] - runewidth.StringWidth(value)
	if spacesLeft > 0 {
		for i := 0; i < spacesLeft; i++ {
			value += " "
		}
	}
	return value
}

// printTableDash - Prints table dash. Markdown or not depending on settings.
func (t *Table) printTableDash() {
	if t.Markdown {
		t.printMarkdownDash()
	} else {
		t.printDash()
	}
}

// stringTableDash - output table dash. Markdown or not depending on settings.
func (t *Table) stringTableDash() string {
	if t.Markdown {
		return t.stringMarkdownDash()
	}
	return t.stringDash()
}

// printDash - Prints dash (on top and header).
func (t *Table) printDash() {
	s := "|"
	for i := 0; i < t.lineLength()-2; i++ {
		s += "-"
	}
	s += "|"
	fmt.Println(s)
}

// stringDash - output dash (on top and header).
func (t *Table) stringDash() string {
	s := "|"
	for i := 0; i < t.lineLength()-2; i++ {
		s += "-"
	}
	s += "|"
	return s
}

// printMarkdownDash - Prints dash in middle of table.
func (t *Table) printMarkdownDash() {
	r := make(map[string]string)
	for _, name := range t.Fields {
		r[name] = strings.Repeat("-", t.fieldSizes[name]-2)
	}
	fmt.Println(t.rowString(r))
}

// stringMarkdownDash - output dash in middle of table.
func (t *Table) stringMarkdownDash() string {
	r := make(map[string]string)
	for _, name := range t.Fields {
		r[name] = strings.Repeat("-", t.fieldSizes[name]-2)
	}
	return t.rowString(r)
}

// lineLength - Counts size of table line length (with spaces etc.).
func (t *Table) lineLength() (sum int) {
	for _, l := range t.fieldSizes {
		sum += l + 1
	}
	return sum + 1
}

func (t *Table) calculateSizes(row map[string]string) {
	for _, k := range t.Fields {
		v, ok := row[k]
		if !ok {
			continue
		}

		vlen := runewidth.StringWidth(v)
		// align to field name length
		if klen := runewidth.StringWidth(k); vlen < klen {
			vlen = klen
		}
		vlen += 2 // + 2 spaces
		if t.fieldSizes[k] < vlen {
			t.fieldSizes[k] = vlen
		}
	}
}

func mapToRows(m map[string]interface{}) (rows []map[string]interface{}) {
	rows = []map[string]interface{}{}
	for key, value := range m {
		row := map[string]interface{}{}
		row["Key"] = strings.Title(key)
		row["Value"] = value
		rows = append(rows, row)
	}
	return
}
