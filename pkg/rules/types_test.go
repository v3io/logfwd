package rules

import (
	"testing"

	"encoding/json"
	"fmt"
	"strings"

	"github.com/stretchr/testify/assert"
)

var data = `
admin:
  namespaces:
    - xx
    - zz
    - yy
  output:
    http:
      method: POST
      endpoint: example.com:9090
      headers:
        EXAMPLE: -1
      authentication:
        header:
          AUTH: open

output:
  - namespace: endpoint1 
    http:
      method: POST
      endpoint: example1.com:9090
      headers:
        EXAMPLE: 1
      authentication:
        header:
          AUTH1: open
  - namespace: endpoint2
    http:
      method: POST
      endpoint: example2.com:9090
      headers:
        EXAMPLE: 2
      authentication:
        header:
          AUTH2: open
`

func TestRead(t *testing.T) {
	v, err := fromReader(strings.NewReader(data))
	assert.NoError(t, err)
	fmt.Printf("%+v\n", v)
	output, err := json.Marshal(v)
	assert.NoError(t, err)
	fmt.Printf("%s\n", string(output))
}
