package dtrack

// FetchAll is a convenience function to retrieve all items of a paginated API resource.
func FetchAll[T any](f func(po PageOptions) (Page[T], error)) (items []T, err error) {
	const pageSize = 50

	var (
		page       Page[T]
		pageNumber = 1
	)

	for {
		page, err = f(PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			break
		}

		items = append(items, page.Items...)
		if len(page.Items) == 0 || len(items) >= page.TotalCount {
			break
		}

		pageNumber++
	}

	return
}
