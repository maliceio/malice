package jsonpointer

import (
	"reflect"
	"testing"
)

type address struct {
	Street string `json:"street"`
	Zip    string
}

type person struct {
	Name               string `json:"name,omitempty"`
	Twitter            string
	Aliases            []string   `json:"aliases"`
	Addresses          []*address `json:"addresses"`
	NameTildeContained string     `json:"name~contained"`
	NameSlashContained string     `json:"name/contained"`
	AnActualArray      [4]int
	NestedMap          map[string]float64 `json:"nestedmap"`
	MapIntKey          map[int]string
	MapUintKey         map[uint]string
	MapFloatKey        map[float64]string
}

var input = &person{
	Name:    "marty",
	Twitter: "mschoch",
	Aliases: []string{
		"jabroni",
		"beer",
	},
	Addresses: []*address{
		&address{
			Street: "123 Sesame St.",
			Zip:    "99099",
		},
	},
	NameTildeContained: "yessir",
	NameSlashContained: "nosir",
	AnActualArray:      [4]int{0, 1, 2, 3},
	NestedMap: map[string]float64{
		"pi":          3.14,
		"back/saidhe": 2.71,
		"till~duh":    1.41,
	},
	MapIntKey: map[int]string{
		1: "one",
		2: "two",
	},
	MapUintKey: map[uint]string{
		3: "three",
		4: "four",
	},
	MapFloatKey: map[float64]string{
		3.14: "pi",
		4.15: "notpi",
	},
}

func benchReflect(b *testing.B, path string) {
	for i := 0; i < b.N; i++ {
		if Reflect(input, path) == nil {
			b.FailNow()
		}
	}
}

func BenchmarkReflectRoot(b *testing.B) {
	benchReflect(b, "")
}

func BenchmarkReflectToplevelExact(b *testing.B) {
	benchReflect(b, "/Twitter")
}

func BenchmarkReflectToplevelTagged(b *testing.B) {
	benchReflect(b, "/Name")
}

func BenchmarkReflectToplevelTaggedLower(b *testing.B) {
	benchReflect(b, "/name")
}

func BenchmarkReflectDeep(b *testing.B) {
	benchReflect(b, "/addresses/0/Zip")
}

func BenchmarkReflectSlash(b *testing.B) {
	benchReflect(b, "/name~1contained")
}

func BenchmarkReflectTilde(b *testing.B) {
	benchReflect(b, "/name~0contained")
}

func compareStringArrayIgnoringOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	tmp := make(map[string]bool, len(a))
	for _, av := range a {
		tmp[av] = true
	}
	for _, bv := range b {
		if tmp[bv] != true {
			return false
		}
	}
	return true
}

func TestReflectListPointers(t *testing.T) {
	pointers, err := ReflectListPointers(input)
	if err != nil {
		t.Fatal(err)
	}
	expect := []string{"", "/name", "/Twitter", "/aliases",
		"/aliases/0", "/aliases/1", "/addresses", "/addresses/0",
		"/addresses/0/street", "/addresses/0/Zip",
		"/name~0contained", "/name~1contained", "/AnActualArray",
		"/AnActualArray/0", "/AnActualArray/1", "/AnActualArray/2",
		"/AnActualArray/3", "/nestedmap", "/nestedmap/pi",
		"/nestedmap/back~1saidhe", "/nestedmap/till~0duh",
		"/MapIntKey", "/MapIntKey/1", "/MapIntKey/2", "/MapUintKey",
		"/MapUintKey/3", "/MapUintKey/4", "/MapFloatKey",
		"/MapFloatKey/3.14", "/MapFloatKey/4.15"}
	if !compareStringArrayIgnoringOrder(expect, pointers) {
		t.Fatalf("expected %#v, got %#v", expect, pointers)
	}
}

func TestReflectNonObjectOrSlice(t *testing.T) {
	got := Reflect(36, "/test")
	if got != nil {
		t.Errorf("expected nil, got %#v", got)
	}
}

type structThatCanBeUsedAsKey struct {
	name   string
	domain string
}

func TestReflectMapThatWontWork(t *testing.T) {

	amapthatwontwork := map[structThatCanBeUsedAsKey]string{}
	akey := structThatCanBeUsedAsKey{name: "marty", domain: "couchbase"}
	amapthatwontwork[akey] = "verycontrived"

	got := Reflect(amapthatwontwork, "/anykey")
	if got != nil {
		t.Errorf("expected nil, got %#v", got)
	}
}

func TestReflect(t *testing.T) {

	tests := []struct {
		path string
		exp  interface{}
	}{
		{
			path: "",
			exp:  input,
		},
		{
			path: "/", exp: nil,
		},
		{
			path: "/name",
			exp:  "marty",
		},
		{
			path: "/Name",
			exp:  "marty",
		},
		{
			path: "/Twitter",
			exp:  "mschoch",
		},
		{
			path: "/aliases/0",
			exp:  "jabroni",
		},
		{
			path: "/Aliases/0",
			exp:  "jabroni",
		},
		{
			path: "/addresses/0/street",
			exp:  "123 Sesame St.",
		},
		{
			path: "/addresses/4/street",
			exp:  nil,
		},
		{
			path: "/doesntexist",
			exp:  nil,
		},
		{
			path: "/does/not/exit",
			exp:  nil,
		},
		{
			path: "/doesntexist/7",
			exp:  nil,
		},
		{
			path: "/name~0contained",
			exp:  "yessir",
		},
		{
			path: "/name~1contained",
			exp:  "nosir",
		},
		{
			path: "/AnActualArray/2",
			exp:  2,
		},
		{
			path: "/AnActualArray/5",
			exp:  nil,
		},
		{
			path: "/nestedmap/pi",
			exp:  3.14,
		},
		{
			path: "/nestedmap/back~1saidhe",
			exp:  2.71,
		},
		{
			path: "/nestedmap/till~0duh",
			exp:  1.41,
		},
		{
			path: "/MapIntKey/1",
			exp:  "one",
		},
		{
			path: "/MapUintKey/3",
			exp:  "three",
		},
		{
			path: "/MapFloatKey/3.14",
			exp:  "pi",
		},
		{
			path: "/MapFloatKey/4.0",
			exp:  nil,
		},
	}

	for _, test := range tests {
		output := Reflect(input, test.path)
		if !reflect.DeepEqual(output, test.exp) {
			t.Errorf("Expected %#v for %q, got %#v", test.exp, test.path, output)
		}
	}
}
