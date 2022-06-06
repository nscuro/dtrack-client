package dtrack

// FetchAll is a convenience function to retrieve all items of a paginated API resource.
func FetchAll[T any](f func(po PageOptions) (Page[T], error)) (items []T, err error) {
	const pageSize = 50
	pageNumber := 1

	for {
		page, fErr := f(PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if fErr != nil {
			err = fErr
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
