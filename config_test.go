package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {

	c := &Config{
		LogDir:          "./logs",
		RestartWaitTime: (10 * time.Hour).String(),
		Processes:       nil,
	}

	j, err := json.Marshal(c)
	require.NoError(t, err)
	fmt.Println(string(j))

	d, err := time.ParseDuration("1m")
	require.NoError(t, err)
	t.Log(d)

}
