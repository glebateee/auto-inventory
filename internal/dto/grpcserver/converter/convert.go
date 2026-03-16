package converter

import (
	"math"
	"sort"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/dto/grpcserver"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func getBiggestThree(grid [][]int) []int {
	n, m := len(grid), len(grid[0])
	sums := make([]int, 0, m*n)
	for i := range n {
		for j := range m {
			maxS := min(i, j, n-i-1, m-j-1)
			for s := range maxS + 1 {
				total := 0
				if s == 0 {
					total = grid[i][j]
				} else {
					total += grid[i-s][j] + grid[i+s][j]
					for k := 1; k <= s+1; k++ {
						total += grid[i-s+k][j-(s-int(math.Abs(float64(s)-float64(k))))] +
							grid[i-s+k][j+(s-int(math.Abs(float64(s)-float64(k))))]
					}
				}
				sums = append(sums, total)
			}
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sums)))

	if len(sums) > 3 {
		sums = sums[:3]
	}
	return sums
}

func dtoToUpdateFields(dto *grpcserver.UpdateProductDTO, mask *fieldmaskpb.FieldMask) *models.UpdateProductFields {
	model := &models.UpdateProductFields{}
	for _, key := range mask.GetPaths() {
		switch key {
		case "name":
			model.Name = &dto.Name
		case "description":
			model.Description = &dto.Description
		case "category":
			model.Category = &dto.Category
		case "manufacturer":
			model.Manufacturer = &dto.Manufacturer
		case "weight":
			model.Weight = &dto.Weight
		case "price":
			model.Price = &dto.Price
		case "base_price":
			model.BasePrice = &dto.BasePrice
		case "issue_year":
			// IssueYear in DTO is int64, domain expects *int16
			val := int16(dto.IssueYear)
			model.IssueYear = &val
		}
	}
	return model
}
