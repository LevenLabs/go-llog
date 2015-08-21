package llog

import (
	"bytes"
	"io/ioutil"
	"regexp"
	. "testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLLog(t *T) {
	// Unfortunately due to the nature of the package all testing involving Out
	// must be syncronously
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	Out = buf

	assertOut := func(expected string) {
		out, err := buf.ReadString('\n')
		require.Nil(t, err)
		assert.Equal(t, expected, out)
	}

	SetLevelFromString("INFO")
	Debug("foo")
	Info("bar")
	Warn("baz")
	Error("buz")
	time.Sleep(100 * time.Millisecond)
	assertOut("~ INFO -- bar\n")
	assertOut("~ WARN -- baz\n")
	assertOut("~ ERROR -- buz\n")

	SetLevelFromString("WARN")
	Debug("foo")
	Info("bar")
	Warn("baz")
	Error("buz", KV{"a": "b"})
	time.Sleep(100 * time.Millisecond)
	assertOut("~ WARN -- baz\n")
	assertOut("~ ERROR -- buz -- a=b\n")
}

func TestEntryPrintOut(t *T) {
	assertEntry := func(postfix string, e entry) {
		expectedRegex := regexp.MustCompile(`^~ ` + postfix + `\n$`)
		expectedRegexTS := regexp.MustCompile(`^~ \[[^\]]+\] ` + postfix + `\n$`)

		buf := bytes.NewBuffer(make([]byte, 0, 128))

		require.Nil(t, e.printOut(buf, false))
		require.Nil(t, e.printOut(buf, true))

		noTS, err := buf.ReadString('\n')
		require.Nil(t, err)
		assert.True(t, expectedRegex.MatchString(noTS), "regex: %q line: %q", expectedRegex.String(), noTS)

		withTS, err := buf.ReadString('\n')
		require.Nil(t, err)
		assert.True(t, expectedRegexTS.MatchString(withTS), "regex: %q line: %q", expectedRegexTS.String(), withTS)
	}

	e := entry{
		level: InfoLevel,
		msg:   "this is a test",
	}
	assertEntry("INFO -- this is a test", e)

	e.kv = KV{}
	assertEntry("INFO -- this is a test", e)

	e.kv = KV{
		"foo": "a",
	}
	assertEntry("INFO -- this is a test -- foo=a", e)

	e.kv = KV{
		"foo": "a",
		"bar": "a",
	}
	assertEntry("INFO -- this is a test -- (foo|bar)=a (foo|bar)=a", e)
}

func BenchmarkLLog(b *B) {
	Out = ioutil.Discard
	for n := 0; n < b.N; n++ {
		Info("This is a generic message", KV{"foo": "bar"})
	}
}
