package util

func FindDifferences(cur, new []string) (add, remove []string) {
	// create map store "cur" element
	curMap := make(map[string]bool)
	for _, v := range cur {
		curMap[v] = true
	}

	// Finds elements in "new" but not in "cur"
	for _, v := range new {
		if _, ok := curMap[v]; !ok {
			add = append(add, v)
		}
	}

	// create map store "new" element
	newMap := make(map[string]bool)
	for _, v := range new {
		newMap[v] = true
	}

	// Finds elements in "cur" but not in "new"
	for _, v := range cur {
		if _, ok := newMap[v]; !ok {
			remove = append(remove, v)
		}
	}

	return add, remove
}

func Paginate[T any](data []T, page, pageSize int) []T {
	// Check page and pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Calc offset
	startIndex := (page - 1) * pageSize
	endIndex := min(startIndex+pageSize, len(data))

	// Return results
	return data[startIndex:endIndex]
}
