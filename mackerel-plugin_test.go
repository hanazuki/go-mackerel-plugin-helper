package mackerelplugin

import (
	"bytes"
	"math"
	"testing"
	"time"
)

func TestCalcDiff(t *testing.T) {
	var mp MackerelPlugin

	val1 := 10.0
	val2 := 0.0
	now := time.Now()
	last := time.Unix(now.Unix()-10, 0)

	diff, err := mp.calcDiff(val1, now, val2, last)
	if diff != 60 {
		t.Errorf("calcDiff: %f should be %f", diff, 60.0)
	}
	if err != nil {
		t.Error("calcDiff causes an error")
	}
}

func TestCalcDiffWithUInt32WithReset(t *testing.T) {
	var mp MackerelPlugin

	val := uint32(10)
	now := time.Now()
	lastval := uint32(12345)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := mp.calcDiffUint32(val, now, lastval, last, 10)
	if err != nil {
	} else {
		t.Errorf("calcDiffUint32 with counter reset should cause an error: %f", diff)
	}
}

func TestCalcDiffWithUInt32Overflow(t *testing.T) {
	var mp MackerelPlugin

	val := uint32(10)
	now := time.Now()
	lastval := math.MaxUint32 - uint32(10)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := mp.calcDiffUint32(val, now, lastval, last, 10)
	if diff != 21.0 {
		t.Errorf("calcDiff: last: %d, now: %d, %f should be %f", val, lastval, diff, 21.0)
	}
	if err != nil {
		t.Error("calcDiff causes an error")
	}
}

func TestCalcDiffWithUInt64WithReset(t *testing.T) {
	var mp MackerelPlugin

	val := uint64(10)
	now := time.Now()
	lastval := uint64(12345)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := mp.calcDiffUint64(val, now, lastval, last, 10)
	if err != nil {
	} else {
		t.Errorf("calcDiffUint64 with counter reset should cause an error: %f", diff)
	}
}

func TestCalcDiffWithUInt64Overflow(t *testing.T) {
	var mp MackerelPlugin

	val := uint64(10)
	now := time.Now()
	lastval := math.MaxUint64 - uint64(10)
	last := time.Unix(now.Unix()-60, 0)

	diff, err := mp.calcDiffUint64(val, now, lastval, last, 10)
	if diff != 21.0 {
		t.Errorf("calcDiff: last: %d, now: %d, %f should be %f", val, lastval, diff, 21.0)
	}
	if err != nil {
		t.Error("calcDiff causes an error")
	}
}

func TestPrintValueUint32(t *testing.T) {
	var mp MackerelPlugin
	s := new(bytes.Buffer)
	var now = time.Unix(1437227240, 0)
	mp.printValue(s, "test", uint32(10), now)

	expected := []byte("test\t10\t1437227240\n")

	if bytes.Compare(expected, s.Bytes()) != 0 {
		t.Fatalf("not matched, expected: %s, got: %s", expected, s)
	}
}

func TestPrintValueUint64(t *testing.T) {
	var mp MackerelPlugin
	s := new(bytes.Buffer)
	var now = time.Unix(1437227240, 0)
	mp.printValue(s, "test", uint64(10), now)

	expected := []byte("test\t10\t1437227240\n")

	if bytes.Compare(expected, s.Bytes()) != 0 {
		t.Fatalf("not matched, expected: %s, got: %s", expected, s)
	}
}

func TestPrintValueFloat64(t *testing.T) {
	var mp MackerelPlugin
	s := new(bytes.Buffer)
	var now = time.Unix(1437227240, 0)
	mp.printValue(s, "test", float64(10.0), now)

	expected := []byte("test\t10.000000\t1437227240\n")

	if bytes.Compare(expected, s.Bytes()) != 0 {
		t.Fatalf("not matched, expected: %s, got: %s", expected, s)
	}
}

func ExampleFormatValues() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(1000)}
	lastStat := map[string]interface{}{"cmd_get": uint64(500), ".last_diff.cmd_get": 300.0}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.cmd_get	500.000000	1437227240
}

func ExampleFormatValuesWithCounterReset() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(10)}
	lastStat := map[string]interface{}{"cmd_get": uint64(500), ".last_diff.cmd_get": 300.0}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
}

func ExampleFormatValuesWithOverflow() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"cmd_get": uint64(500)}
	lastStat := map[string]interface{}{"cmd_get": uint64(math.MaxUint64 - 100), ".last_diff.cmd_get": float64(1.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValues(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
}

func ExampleFormatValuesWithWildcard() {
	var mp MackerelPlugin
	prefix := "foo.#"
	metric := Metrics{Name: "bar", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"foo.1.bar": uint64(1000), "foo.2.bar": uint64(2000)}
	lastStat := map[string]interface{}{"foo.1.bar": uint64(500), ".last_diff.foo.1.bar": float64(2.0), "foo.2.bar": uint64(1000), ".last_diff.foo.2.bar": float64(1.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValuesWithWildcard(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.1.bar	500.000000	1437227240
	// foo.2.bar	1000.000000	1437227240
}

func ExampleFormatValuesWithWildcardAndNoDiff() {
	var mp MackerelPlugin
	prefix := "foo.#"
	metric := Metrics{Name: "bar", Label: "Get", Diff: false}
	stat := map[string]interface{}{"foo.1.bar": float64(1000), "foo.2.bar": float64(2000)}
	lastStat := map[string]interface{}{"foo.1.bar": float64(500), ".last_diff.foo.1.bar": float64(2.0), "foo.2.bar": float64(1000), ".last_diff.foo.2.bar": float64(1.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValuesWithWildcard(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.1.bar	1000.000000	1437227240
	// foo.2.bar	2000.000000	1437227240
}

func ExampleFormatValuesWithWildcardAstarisk() {
	var mp MackerelPlugin
	prefix := "foo"
	metric := Metrics{Name: "*", Label: "Get", Diff: true, Type: "uint64"}
	stat := map[string]interface{}{"foo.1": uint64(1000), "foo.2": uint64(2000)}
	lastStat := map[string]interface{}{"foo.1": uint64(500), ".last_diff.foo.1": float64(2.0), "foo.2": uint64(1000), ".last_diff.foo.2": float64(1.0)}
	now := time.Unix(1437227240, 0)
	lastTime := now.Add(-time.Duration(60) * time.Second)
	mp.formatValuesWithWildcard(prefix, metric, &stat, &lastStat, now, lastTime)

	// Output:
	// foo.1	500.000000	1437227240
	// foo.2	1000.000000	1437227240
}

// an example implementation
type MemcachedPlugin struct {
}

var graphdef map[string](Graphs) = map[string](Graphs){
	"memcached.cmd": Graphs{
		Label: "Memcached Command",
		Unit:  "integer",
		Metrics: [](Metrics){
			Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"},
		},
	},
}

func (m MemcachedPlugin) GraphDefinition() map[string](Graphs) {
	return graphdef
}

func (m MemcachedPlugin) FetchMetrics() (map[string]interface{}, error) {
	var stat map[string]interface{}
	return stat, nil
}

func ExampleOutputDefinitions() {
	var mp MemcachedPlugin
	helper := NewMackerelPlugin(mp)
	helper.OutputDefinitions()

	// Output:
	// # mackerel-agent-plugin
	// {"graphs":{"memcached.cmd":{"label":"Memcached Command","unit":"integer","metrics":[{"name":"cmd_get","label":"Get","diff":true,"type":"uint64","stacked":false,"scale":0}]}}}
}
