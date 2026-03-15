package tests

import (
	"math/rand/v2"
	"testing"

	"github.com/glebateee/auto-inventory/tests/suite"
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductPageSize(t *testing.T) {
	ctx, st := suite.New(t)
	var page int64 = rand.Int64N(5) + 1
	var size int64 = rand.Int64N(5) + 1
	//var total int64 = 20
	resp, err := st.AuthClient.ProductPageSize(ctx, &aiv1.ProductPageSizeRequest{
		Page: page,
		Size: size,
	})
	products := resp.GetProducts()
	//available := resp.GetAvailable()
	require.NoError(t, err)
	//assert.Equal(t, total, available)
	assert.Equal(t, (page-1)*size+1, products[0].Id)
	assert.Equal(t, (page-1)*size+size, products[len(products)-1].Id)
}
