// Package game contains core game types and logic
package game

type Item struct {
	ID          string
	Name        string
	Quantity    int
	Value       int
	Tags        []string
	Category    string
	Description string
	Equipped    bool
}

type ItemCategory struct {
	Name     string
	Children []string // Sub-categories
	Filter   []string // Tags to filter by
}

func GetSampleItems() []Item {
	return []Item{
		{ID: "wood_oak", Name: "Oak Wood", Quantity: 150, Value: 5, Tags: []string{"resource", "wood"}, Category: "Resources"},
		{ID: "ore_iron", Name: "Iron Ore", Quantity: 45, Value: 10, Tags: []string{"resource", "ore"}, Category: "Resources"},
		{ID: "ore_coal", Name: "Coal", Quantity: 23, Value: 8, Tags: []string{"resource", "ore"}, Category: "Resources"},
		{ID: "fish_trout", Name: "Raw Trout", Quantity: 12, Value: 15, Tags: []string{"resource", "fish"}, Category: "Resources"},
		{ID: "stone", Name: "Stone", Quantity: 89, Value: 2, Tags: []string{"resource", "stone"}, Category: "Resources"},

		{ID: "sword_steel", Name: "Steel Sword", Quantity: 1, Value: 150, Tags: []string{"equipment", "weapon", "sword"}, Category: "Equipment", Description: "+25 ATK"},
		{ID: "sword_iron", Name: "Iron Sword", Quantity: 2, Value: 75, Tags: []string{"equipment", "weapon", "sword"}, Category: "Equipment", Description: "+15 ATK", Equipped: true},
		{ID: "dagger_iron", Name: "Iron Dagger", Quantity: 3, Value: 50, Tags: []string{"equipment", "weapon", "dagger"}, Category: "Equipment", Description: "+15 ATK"},
		{ID: "axe_bronze", Name: "Bronze Axe", Quantity: 1, Value: 30, Tags: []string{"equipment", "tool", "axe"}, Category: "Equipment", Description: "+10 Woodcutting"},
		{ID: "axe_steel", Name: "Steel Axe", Quantity: 1, Value: 100, Tags: []string{"equipment", "tool", "axe"}, Category: "Equipment", Description: "+15 Woodcutting"},
		{ID: "pickaxe_iron", Name: "Iron Pickaxe", Quantity: 1, Value: 60, Tags: []string{"equipment", "tool", "pickaxe"}, Category: "Equipment", Description: "+8 Mining", Equipped: true},

		{ID: "food_bread", Name: "Bread", Quantity: 15, Value: 5, Tags: []string{"consumable", "food"}, Category: "Consumables"},
		{ID: "potion_hp", Name: "Health Potion", Quantity: 8, Value: 25, Tags: []string{"consumable", "potion"}, Category: "Consumables"},
	}
}

func GetCategories() []ItemCategory {
	return []ItemCategory{
		{Name: "All Items", Children: nil, Filter: []string{}},
		{Name: "Resources", Children: []string{"Ores", "Woods", "Fish"}, Filter: []string{"resource"}},
		{Name: "Equipment", Children: []string{"Weapons", "Tools", "Armor"}, Filter: []string{"equipment"}},
		{Name: "Consumables", Children: []string{"Food", "Potions"}, Filter: []string{"consumable"}},
	}
}
