package activity

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestData(t *testing.T) {
	files, err := filepath.Glob("testdata/*.plist")
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range files {
		// Open the test file.
		f, err := os.Open(filename)
		if err != nil {
			t.Error(err)
			continue
		}
		// Read from the test file.
		data, err := ReadDotsData(f)
		if err != nil {
			t.Error(err)
			continue
		}
		// Close the test file.
		if err = f.Close(); err != nil {
			t.Error(err)
		}

		// Make sure Powerups() works, even if the
		// powerup keys aren't present.
		_, _, _ = data.Powerups()
		// Make sure getting and setting of powerups works.
		data.SetPowerups(100, 100, 100)
		a, b, c := data.Powerups()
		if a != 100 && b != 100 && c != 100 {
			t.Error("SetPowerups()")
		}
		// Make sure maximizing the powerups works.
		data.MaximizePowerups()
		a, b, c = data.Powerups()
		if a+1 >= 0 || b+1 >= 0 || c+1 >= 0 {
			t.Error("MaxPowerups()")
		}

		// Test reading and writing on some temp files.
		for i := 0; i < 10; i++ {
			// Create a copy of data with random powerups.
			data := data
			a := rand.Int31()
			b := rand.Int31()
			c := rand.Int31()
			data.SetPowerups(a, b, c)
			// Open a temporary file.
			tmp, err := ioutil.TempFile("", "dotbot")
			if err != nil {
				t.Error(err)
				continue
			}
			// Write the the temp file.
			if err = data.WriteTo(tmp); err != nil {
				t.Error(err)
				continue
			}
			// Close the temp file.
			if err := tmp.Close(); err != nil {
				t.Error(err)
			}
			// Reopen the temp file.
			tmp, err = os.Open(tmp.Name())
			if err != nil {
				t.Error(err)
				continue
			}
			// Read the temp file.
			data, err = ReadDotsData(tmp)
			if err != nil {
				fmt.Println(tmp, tmp.Name())
				t.Error(err)
				continue
			}
			// Make sure the powerups didn't change.
			a1, b1, c1 := data.Powerups()
			if a != a1 || b != b1 || c != c1 {
				t.Error("Powerups were not presevered")
			}
			// Close the temp file again.
			if err := tmp.Close(); err != nil {
				t.Error(err)
			}
			// Delete the temp file.
			if err := os.Remove(tmp.Name()); err != nil {
				t.Error(err)
			}
		}
	}
}
