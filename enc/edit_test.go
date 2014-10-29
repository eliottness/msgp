package enc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLocate(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)
	en.WriteMapHeader(2)
	en.WriteString("thing_one")
	en.WriteString("value_one")
	en.WriteString("thing_two")
	en.WriteFloat64(2.0)

	field := Locate("thing_one", buf.Bytes())
	if len(field) == 0 {
		t.Fatal("field not found")
	}

	var zbuf bytes.Buffer
	NewEncoder(&zbuf).WriteString("value_one")

	if !bytes.Equal(zbuf.Bytes(), field) {
		t.Errorf("got %q; wanted %q", field, zbuf.Bytes())
	}

	zbuf.Reset()
	NewEncoder(&zbuf).WriteFloat64(2.0)
	field = Locate("thing_two", buf.Bytes())
	if len(field) == 0 {
		t.Fatal("field not found")
	}
	if !bytes.Equal(zbuf.Bytes(), field) {
		t.Errorf("got %q; wanted %q", field, zbuf.Bytes())
	}

	field = Locate("nope", buf.Bytes())
	if len(field) != 0 {
		t.Fatalf("wanted %q; got %q", ErrFieldNotFound, nil)
	}

}

func TestReplace(t *testing.T) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)
	en.WriteMapHeader(2)
	en.WriteString("thing_one")
	en.WriteString("value_one")
	en.WriteString("thing_two")
	en.WriteFloat64(2.0)

	var fbuf bytes.Buffer
	NewEncoder(&fbuf).WriteFloat64(4.0)

	// replace 2.0 with 4.0 in field two
	raw, err := Replace("thing_two", buf.Bytes(), fbuf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	m := make(map[string]interface{})
	m, _, err = ReadMapStrIntfBytes(raw, m)
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m["thing_two"], 4.0) {
		t.Errorf("wanted %v; got %v", 4.0, m["thing_two"])
	}

	// replace 2.0 with []byte("hi!")
	fbuf.Reset()
	NewEncoder(&fbuf).WriteBytes([]byte("hello there!"))
	raw, err = Replace("thing_two", raw, fbuf.Bytes())
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	m, _, err = ReadMapStrIntfBytes(raw, m)
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m["thing_two"], []byte("hello there!")) {
		t.Errorf("wanted %v; got %v", []byte("hello there!"), m["thing_two"])
	}
}

func BenchmarkLocate(b *testing.B) {
	var buf bytes.Buffer
	en := NewEncoder(&buf)
	en.WriteMapHeader(3)
	en.WriteString("thing_one")
	en.WriteString("value_one")
	en.WriteString("thing_two")
	en.WriteFloat64(2.0)
	en.WriteString("thing_three")
	en.WriteBytes([]byte("hello!"))

	raw := buf.Bytes()
	// bytes/s will be the number of bytes traversed per unit of time
	field := Locate("thing_three", raw)
	b.SetBytes(int64(len(raw) - len(field)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Locate("thing_three", raw)
	}
}
