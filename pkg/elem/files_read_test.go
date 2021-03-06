package elem

import (
	"testing"

	"github.com/Bitspark/slang/pkg/core"
	"github.com/Bitspark/slang/tests/assertions"
	"github.com/stretchr/testify/require"
)

func Test_FileRead__IsRegistered(t *testing.T) {
	a := assertions.New(t)

	ocFileRead := getBuiltinCfg(filesReadId)
	a.NotNil(ocFileRead)
}

func Test_FileRead__InPorts(t *testing.T) {
	a := assertions.New(t)

	o, err := buildOperator(
		core.InstanceDef{
			Operator: filesReadId,
		},
	)
	require.NoError(t, err)

	a.NotNil(o.Main().In())
	a.Equal(core.TYPE_STRING, o.Main().In().Type())
}

func Test_FileRead__OutPorts(t *testing.T) {
	a := assertions.New(t)

	o, err := buildOperator(
		core.InstanceDef{
			Operator: filesReadId,
		},
	)
	require.NoError(t, err)

	a.NotNil(o.Main().Out())
	a.Equal(core.TYPE_MAP, o.Main().Out().Type())
	a.Equal(core.TYPE_BINARY, o.Main().Out().Map("content").Type())
	a.Equal(core.TYPE_STRING, o.Main().Out().Map("error").Type())
}

func Test_FileRead__Simple(t *testing.T) {
	a := assertions.New(t)

	o, err := buildOperator(
		core.InstanceDef{
			Operator: filesReadId,
		},
	)
	require.NoError(t, err)

	o.Main().Out().Bufferize()
	o.Start()

	o.Main().In().Push("../../tests/test_data/hello.txt")
	a.Equal(core.Binary("hello slang"), o.Main().Out().Map("content").Pull())
	a.Nil(o.Main().Out().Map("error").Pull())
}

func Test_FileRead__NotFound(t *testing.T) {
	a := assertions.New(t)

	o, err := buildOperator(
		core.InstanceDef{
			Operator: filesReadId,
		},
	)
	require.NoError(t, err)

	o.Main().Out().Bufferize()
	o.Start()

	o.Main().In().Push("./tests/test_data/nonexistentfile")
	a.Nil(o.Main().Out().Map("content").Pull())
	a.NotNil(o.Main().Out().Map("error").Pull())
}
