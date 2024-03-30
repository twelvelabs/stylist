package stylist

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandType(t *testing.T) {
	names := CommandTypeNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseCommandType(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestInputType(t *testing.T) {
	names := InputTypeNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseInputType(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestOutputType(t *testing.T) {
	names := OutputTypeNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseOutputType(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestOutputFormat(t *testing.T) {
	names := OutputFormatNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseOutputFormat(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestResultLevel(t *testing.T) {
	names := ResultLevelNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseResultLevel(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)

	_ = ResultLevel(99).String()
}

func TestCoerceResultLevel(t *testing.T) {
	tests := []struct {
		desc     string
		expected ResultLevel
		err      string
	}{
		{
			desc:     "",
			expected: ResultLevelNone,
			err:      "",
		},
		{
			desc:     "<no value>",
			expected: ResultLevelNone,
			err:      "",
		},
		{
			desc:     "info",
			expected: ResultLevelInfo,
			err:      "",
		},
		{
			desc:     "note",
			expected: ResultLevelInfo,
			err:      "",
		},
		{
			desc:     "warn",
			expected: ResultLevelWarning,
			err:      "",
		},
		{
			desc:     "warning",
			expected: ResultLevelWarning,
			err:      "",
		},
		{
			desc:     "err",
			expected: ResultLevelError,
			err:      "",
		},
		{
			desc:     "error",
			expected: ResultLevelError,
			err:      "",
		},
		{
			desc:     "unknown",
			expected: ResultLevelError,
			err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := CoerceResultLevel(tt.desc)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}

			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestResultFormat(t *testing.T) {
	names := ResultFormatNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseResultFormat(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestResultPath(t *testing.T) {
	names := ResultPathNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseResultPath(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestResultSort(t *testing.T) {
	names := ResultSortNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseResultSort(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}

func TestLogLevel(t *testing.T) {
	names := LogLevelNames()
	require.True(t, len(names) > 0)

	name := names[0]
	enum, _ := ParseLogLevel(name)
	require.True(t, enum.IsValid())
	require.Equal(t, enum, enum.Get())
	require.NoError(t, enum.Set(enum.String()))

	enumType := strings.TrimPrefix(fmt.Sprintf("%T", enum), "stylist.")
	require.Equal(t, enumType, enum.Type())

	marshalled, err := enum.MarshalText()
	require.NoError(t, err)
	err = enum.UnmarshalText(marshalled)
	require.NoError(t, err)
	err = enum.UnmarshalText([]byte{})
	require.Error(t, err)
}
