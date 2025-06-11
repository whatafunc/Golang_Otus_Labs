package main

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestReadDir(t *testing.T) {
	dir := "testdata/env"
	// Test case: Happy path (valid directory)
	t.Run("InvalidDir", func(t *testing.T) {
		res, err := ReadDir(dir + "no-folder")
		assert.Assert(t, cmp.Nil(res)) // Check res is nil

		// Check error is os.ErrNotExist
		assert.Assert(t, cmp.ErrorIs(err, os.ErrNotExist))

		// Double check: Check error contains a substring
		assert.ErrorContains(t, err, "no such file or directory")
	})
	t.Run("ValidDir", func(t *testing.T) {
		res, err := ReadDir(dir)
		assert.Assert(t, cmp.ErrorIs(err, nil)) // Check error is nil
		assert.NilError(t, err)
		assert.Equal(t, len(res), 5)
		expectedValues := map[string]EnvValue{
			"BAR": {
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": {
				Value:      "",
				NeedRemove: true, // Assuming empty file or empty first line
			},
			// TO-DO: add other Golden files...
		}
		assert.DeepEqual(t, res["BAR"], EnvValue{Value: "bar", NeedRemove: false})

		for key, expected := range expectedValues {
			actual, exists := res[key]
			assert.Assert(t, exists, "key %s not found", key)
			assert.DeepEqual(t, actual, expected)
		}
		// log.Println("------ dir is OK")
		// log.Println("-struct = ", res)
	})
	t.Run("ValidEnvVars", func(t *testing.T) {
		_, err := ReadDir(dir)
		// Double check all
		assert.Assert(t, cmp.ErrorIs(err, nil)) // Check error is nil
		assert.NilError(t, err)
		// Assert environment changes
		assert.Equal(t, "", os.Getenv("EMPTY")) // Should be unset
		assert.Equal(t, "bar", os.Getenv("BAR"))
		// assert.Equal(t, "a\nb", os.Getenv("NULLTEST"))
		// assert.Equal(t, "value", os.Getenv("SPACES"))
		// log.Println("------ dir is OK")
		// log.Println("-struct = ", res)
	})
}
