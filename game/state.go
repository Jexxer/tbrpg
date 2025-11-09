package game

import "time"

// State holds all game-related state
type State struct {
	Storage          *Storage
	ActivityLog      *ActivityLog
	SelectedCategory string
	SavedSearches    []SavedSearch
	// Future: Player stats, quests, equipment, etc.
}

// SavedSearch represents a saved search query
type SavedSearch struct {
	Name  string
	Query string
}

// NewState creates a new game state with default values
func NewState() *State {
	// Initialize with sample items
	storage := NewStorage(GetSampleItems())

	// Initialize activity log with sample entries
	activityLog := NewActivityLog()
	activityLog.AddEntry("System", "Game started", "Welcome to TBRPG!")
	activityLog.AddEntry("Navigation", "Traveled to Starting Town", "")
	activityLog.AddEntry("Woodcutting", "+1 Oak Log", "(1,235 total) +12 XP")
	activityLog.AddEntry("Combat", "Goblin defeated", "+Goblin Ear (12) +15 XP")
	activityLog.AddEntry("Fishing", "+1 Raw Trout", "(58 total) +8 XP")

	return &State{
		Storage:          storage,
		ActivityLog:      activityLog,
		SelectedCategory: "All Items",
		SavedSearches:    []SavedSearch{},
	}
}

// GetFilteredItems returns items filtered by current category and search term
func (s *State) GetFilteredItems(searchTerm string) []Item {
	return s.Storage.Filter(FilterOptions{
		SearchTerm:     searchTerm,
		CategoryFilter: s.SelectedCategory,
	})
}

// SetCategory updates the selected category
func (s *State) SetCategory(category string) {
	s.SelectedCategory = category
}

// SaveSearch saves a search query with a name
func (s *State) SaveSearch(name, query string) {
	s.SavedSearches = append(s.SavedSearches, SavedSearch{
		Name:  name,
		Query: query,
	})
}

// LoadSearch finds a saved search by name
func (s *State) LoadSearch(name string) (string, bool) {
	for _, saved := range s.SavedSearches {
		if saved.Name == name {
			return saved.Query, true
		}
	}
	return "", false
}

// ActivityLog manages the game's activity log
type ActivityLog struct {
	entries []LogEntry
}

// LogEntry represents a single activity log entry
type LogEntry struct {
	Timestamp time.Time
	Category  string
	Action    string
	Details   string
}

// NewActivityLog creates a new activity log
func NewActivityLog() *ActivityLog {
	return &ActivityLog{
		entries: []LogEntry{},
	}
}

// AddEntry adds a new entry to the log
func (al *ActivityLog) AddEntry(category, action, details string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Category:  category,
		Action:    action,
		Details:   details,
	}
	al.entries = append(al.entries, entry)
}

// GetEntries returns all log entries
func (al *ActivityLog) GetEntries() []LogEntry {
	return al.entries
}

// GetRecentEntries returns the last N entries
func (al *ActivityLog) GetRecentEntries(n int) []LogEntry {
	if n >= len(al.entries) {
		return al.entries
	}
	return al.entries[len(al.entries)-n:]
}

// Clear removes all log entries
func (al *ActivityLog) Clear() {
	al.entries = []LogEntry{}
}
