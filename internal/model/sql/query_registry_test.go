package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryRegistry(t *testing.T) {
	var registry QueryRegistry = QueryRegistry{}
	require.NotNil(t, registry, "invalid query registry instance")
}
