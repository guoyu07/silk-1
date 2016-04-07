package runner_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/silk/parse"
	"github.com/matryer/silk/runner"
	"github.com/matryer/silk/testutil"
)

func TestTInter(t *testing.T) {
	var tt runner.T
	tt = &testing.T{}
	_ = tt
}

func TestRunGroupSuccess(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	g, err := parse.ParseFile("../testfiles/success/echo.success.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.False(subT.Failed())
}

func TestRunFileSuccess(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/echo.success.silk.md")
	is.False(subT.Failed())
}

// https://github.com/matryer/silk/issues/31
func TestIssue31(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/issue-31.silk.md")
	is.False(subT.Failed())
}

// https://github.com/matryer/silk/issues/2
func TestCapturedVars(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoDataHandler())
	defer s.Close()
	os.Setenv("$AppNameFromEnv", "Silk")
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/captured-vars.silk.md")
	is.False(subT.Failed())
}

func TestData(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoDataHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/data.silk.md")
	is.False(subT.Failed())
}

func TestBodyField(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoDataHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/body-as-field.silk.md")
	is.False(subT.Failed())
}

func TestRunFileSuccessNoBody(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/echo.nobody.success.silk.md")
	is.False(subT.Failed())
}

func TestFailureWrongBody(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	var logs []string
	r.Log = func(s string) {
		logs = append(logs, s)
	}
	g, err := parse.ParseFile("../testfiles/failure/echo.failure.wrongbody.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.True(subT.Failed())
	logstr := strings.Join(logs, "\n")
	is.True(strings.Contains(logstr, "body expected:"))
	is.True(strings.Contains(logstr, "GET /echo"))
	is.True(strings.Contains(logstr, "Hello silky."))
	is.True(strings.Contains(logstr, "actual:"))
	is.True(strings.Contains(logstr, "GET /echo"))
	is.True(strings.Contains(logstr, "Hello silk."))
	is.True(strings.Contains(logstr, "--- FAIL: GET /echo"))
	is.True(strings.Contains(logstr, "../testfiles/failure/echo.failure.wrongbody.silk.md:14 - body doesn't match"))
}

func TestFailureWrongHeader(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	var logs []string
	r.Log = func(s string) {
		logs = append(logs, s)
	}
	g, err := parse.ParseFile("../testfiles/failure/echo.failure.wrongheader.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.True(subT.Failed())
	logstr := strings.Join(logs, "\n")

	is.True(strings.Contains(logstr, `Content-Type expected: "wrong/type"  actual: "text/plain; charset=utf-8"`))
	is.True(strings.Contains(logstr, "--- FAIL: GET /echo"))
	is.True(strings.Contains(logstr, "../testfiles/failure/echo.failure.wrongheader.silk.md:22 - Content-Type doesn't match"))
}

func TestGlob(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.Log = func(s string) {} // don't bother logging
	r.RunGlob(filepath.Glob("../testfiles/failure/echo.*.silk.md"))
	is.True(subT.Failed())
}

func TestCookies(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	r.RunFile("../testfiles/success/cookies.silk.md")
	is.False(subT.Failed())
}

func TestFailureFieldsSameType(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	var logs []string
	r.Log = func(s string) {
		logs = append(logs, s)
	}
	g, err := parse.ParseFile("../testfiles/failure/echo.failure.fieldssametype.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.True(subT.Failed())
	logstr := strings.Join(logs, "\n")

	is.True(strings.Contains(logstr, "Status expected: 400  actual: 200"))
}

func TestFailureFieldsDifferentTypes(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	var logs []string
	r.Log = func(s string) {
		logs = append(logs, s)
	}
	g, err := parse.ParseFile("../testfiles/failure/echo.failure.fieldsdifferenttypes.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.True(subT.Failed())
	logstr := strings.Join(logs, "\n")

	is.True(strings.Contains(logstr, `Status expected string: "400"  actual float64: 200`))
}

func TestRunJsonModesSuccess(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoRawHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	g, err := parse.ParseFile("../testfiles/success/echoraw.success.jsonmodes.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.False(subT.Failed())
}

func TestRunJsonModesFailure(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoRawHandler())
	defer s.Close()
	r := runner.New(subT, s.URL)
	g, err := parse.ParseFile("../testfiles/failure/echoraw.failure.jsonmodes.silk.md")
	is.NoErr(err)
	r.RunGroup(g...)
	is.True(subT.Failed())
}

type testT struct {
	log    []string
	failed bool
}

func (t *testT) FailNow() {
	t.failed = true
}

func (t *testT) Failed() bool {
	return t.failed
}

func (t *testT) LogString() string {
	return strings.Join(t.log, "\n")
}

func (t *testT) Log(args ...interface{}) {
	t.log = append(t.log, fmt.Sprint(args...))
}
