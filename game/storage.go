package game

import "strings"

// Storage manages items and provides search/filter functionality
type Storage struct {
	items []Item
}

// NewStorage creates a new storage instance
func NewStorage(items []Item) *Storage {
	return &Storage{
		items: items,
	}
}

// GetItems returns all items in storage
func (s *Storage) GetItems() []Item {
	return s.items
}

// SetItems updates the storage items
func (s *Storage) SetItems(items []Item) {
	s.items = items
}

// FilterOptions contains criteria for filtering items
type FilterOptions struct {
	SearchTerm       string
	CategoryFilter   string
	IncludeEquipped  bool
	TagFilter        []string
}

// Filter returns items matching the given criteria
func (s *Storage) Filter(opts FilterOptions) []Item {
	filtered := []Item{}
	searchTerm := strings.ToLower(opts.SearchTerm)

	for _, item := range s.items {
		// Category filter
		if opts.CategoryFilter != "" && opts.CategoryFilter != "All Items" {
			if !matchesCategory(item, opts.CategoryFilter) {
				continue
			}
		}

		// Search filter (by name or tags)
		if searchTerm != "" {
			if !matchesSearch(item, searchTerm) {
				continue
			}
		}

		// Tag filter
		if len(opts.TagFilter) > 0 {
			if !hasAnyTag(item, opts.TagFilter) {
				continue
			}
		}

		// Equipped filter
		if !opts.IncludeEquipped && item.Equipped {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

// SearchByName returns items matching the search term in their name
func (s *Storage) SearchByName(searchTerm string) []Item {
	return s.Filter(FilterOptions{
		SearchTerm: searchTerm,
	})
}

// GetByCategory returns all items in a specific category
func (s *Storage) GetByCategory(category string) []Item {
	return s.Filter(FilterOptions{
		CategoryFilter: category,
	})
}

// GetByTag returns all items with the specified tag
func (s *Storage) GetByTag(tag string) []Item {
	return s.Filter(FilterOptions{
		TagFilter: []string{tag},
	})
}

// FindByID returns an item by its ID, or nil if not found
func (s *Storage) FindByID(id string) *Item {
	for i, item := range s.items {
		if item.ID == id {
			return &s.items[i]
		}
	}
	return nil
}

// CountByCategory returns the number of items in each category
func (s *Storage) CountByCategory() map[string]int {
	counts := make(map[string]int)
	for _, item := range s.items {
		counts[item.Category]++
	}
	return counts
}

// TotalValue calculates the total value of all items
func (s *Storage) TotalValue() int {
	total := 0
	for _, item := range s.items {
		total += item.Value * item.Quantity
	}
	return total
}

// Helper functions

// matchesCategory checks if an item matches the category filter
func matchesCategory(item Item, category string) bool {
	// Direct category match
	if item.Category == category {
		return true
	}

	// Check if any tags contain the category (case-insensitive)
	categoryLower := strings.ToLower(category)
	for _, tag := range item.Tags {
		if strings.Contains(strings.ToLower(tag), categoryLower) {
			return true
		}
	}

	return false
}

// matchesSearch checks if an item matches the search term
func matchesSearch(item Item, searchTerm string) bool {
	// Check name
	if strings.Contains(strings.ToLower(item.Name), searchTerm) {
		return true
	}

	// Check tags
	for _, tag := range item.Tags {
		if strings.Contains(strings.ToLower(tag), searchTerm) {
			return true
		}
	}

	return false
}

// hasAnyTag checks if an item has any of the specified tags
func hasAnyTag(item Item, tags []string) bool {
	for _, filterTag := range tags {
		for _, itemTag := range item.Tags {
			if strings.EqualFold(itemTag, filterTag) {
				return true
			}
		}
	}
	return false
}
