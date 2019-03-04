package wall

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidType(t *testing.T) {
	re := regexp.MustCompile(`.*`)
	input := "data"

	invalidTests := []struct {
		output  interface{}
		message string
	}{
		{
			output:  "somestring",
			message: "output must be a pointer to a slice, map, or struct",
		},
		{
			output:  []int{},
			message: "output must be a pointer to a slice, map, or struct",
		},
		{
			output:  map[int]string{},
			message: "output must be a pointer to a slice, map, or struct",
		},
		{
			output:  new(int),
			message: "output must be a pointer to a slice, map, or struct",
		},
		{
			output:  new([]int),
			message: "output slice must be []string",
		},
		{
			output:  new(map[string]int),
			message: "output map must be map[string]string",
		},
	}

	for _, it := range invalidTests {
		err := Parse(re, input, it.output)
		require.NotNil(t, err)
		require.Equal(t, err.Error(), it.message)
	}
}

func TestParseSlice(t *testing.T) {
	sliceTests := []struct {
		re     *regexp.Regexp
		input  string
		output []string
	}{
		{
			re:     regexp.MustCompile(`(\d{3})-(\d{3})-(\d{4})`),
			input:  "123-456-7890",
			output: []string{"123", "456", "7890"},
		},
		{
			re:     regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`),
			input:  "1.8.7",
			output: []string{"1", "8", "7"},
		},
		{
			re:     regexp.MustCompile("data"),
			input:  "data",
			output: []string{},
		},
	}

	for _, st := range sliceTests {
		var s []string
		err := Parse(st.re, st.input, &s)
		require.Nil(t, err)
		require.Equal(t, s, st.output)
	}
}

func TestParseMap(t *testing.T) {
	mapTests := []struct {
		re     *regexp.Regexp
		input  string
		output map[string]string
	}{
		{
			re:     regexp.MustCompile(`(\d{3})-(\d{3})-(\d{4})`),
			input:  "123-456-7890",
			output: map[string]string{"1": "123", "2": "456", "3": "7890"},
		},
		{
			re:     regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`),
			input:  "1.8.7",
			output: map[string]string{"major": "1", "minor": "8", "patch": "7"},
		},
		{
			re:     regexp.MustCompile("data"),
			input:  "data",
			output: map[string]string{},
		},
	}

	for _, mt := range mapTests {
		var m map[string]string
		err := Parse(mt.re, mt.input, &m)
		require.Nil(t, err)
		require.Equal(t, m, mt.output)
	}
}

func TestParseStructLabels(t *testing.T) {
	re := regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`)
	input := "1.8.7"
	var s struct {
		Major string `wall:"major"`
		Minor string `wall:"minor"`
		Patch string `wall:"patch"`
	}

	err := Parse(re, input, &s)

	require.Nil(t, err)
	require.Equal(t, s.Major, "1")
	require.Equal(t, s.Minor, "8")
	require.Equal(t, s.Patch, "7")
}

func TestParseStructDefaultLabels(t *testing.T) {
	re := regexp.MustCompile(`(\d{3})-(\d{3})-(\d{4})`)
	input := "123-456-7890"
	var s struct {
		Area   string `wall:"1"`
		Prefix string `wall:"2"`
		Number string `wall:"3"`
	}

	err := Parse(re, input, &s)

	require.Nil(t, err)
	require.Equal(t, s.Area, "123")
	require.Equal(t, s.Prefix, "456")
	require.Equal(t, s.Number, "7890")
}

func TestParseStructExtraData(t *testing.T) {
	re := regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`)
	input := "1.8.7"
	var s struct {
		Major string `wall:"major"`
		Minor string `wall:"minor"`
		Extra int
	}

	s.Extra = 1

	err := Parse(re, input, &s)

	require.Nil(t, err)
	require.Equal(t, s.Major, "1")
	require.Equal(t, s.Minor, "8")
	require.Equal(t, s.Extra, 1)
}

func TestParseStructBadTag(t *testing.T) {
	re := regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`)
	input := "1.8.7"
	var s struct {
		Major string `wall:"major"`
		Minor string `wall:"minor"`
		Patch int    `wall:"patch"`
	}

	err := Parse(re, input, &s)

	require.NotNil(t, err)
	require.Equal(t, err.Error(), `field Patch tagged with "wall" must be a string`)
}

func TestParseStructUnsettableField(t *testing.T) {
	re := regexp.MustCompile(`(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`)
	input := "1.8.7"
	var s struct {
		Major string `wall:"major"`
		Minor string `wall:"minor"`
		patch string `wall:"patch"`
	}

	err := Parse(re, input, &s)

	require.NotNil(t, err)
	require.Equal(t, err.Error(), `field patch is not settable`)
}
