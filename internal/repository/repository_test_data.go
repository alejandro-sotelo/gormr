package repository

// Car model and test cases for repository tests
type Car struct {
	ID    uint
	Brand string
	Color string
	Year  int
	Model string
}

var carTestCases = []struct {
	name     string
	input    Car
	expected Car
}{
	{
		 name:     "create and fetch car",
		 input:    Car{Brand: "Toyota", Color: "Red", Year: 2020, Model: "Corolla"},
		 expected: Car{ID: 1, Brand: "Toyota", Color: "Red", Year: 2020, Model: "Corolla"},
	},
	{
		 name:     "create and fetch another car",
		 input:    Car{Brand: "Ford", Color: "Blue", Year: 2018, Model: "Focus"},
		 expected: Car{ID: 2, Brand: "Ford", Color: "Blue", Year: 2018, Model: "Focus"},
	},
}

// Data for GetAll
var carGetAllTestData = []Car{
	{Brand: "Toyota", Color: "Red", Year: 2020, Model: "Corolla"},
	{Brand: "Ford", Color: "Blue", Year: 2018, Model: "Focus"},
	{Brand: "Peugeot", Color: "White", Year: 2019, Model: "208"},
}

// Data for GetPaginated
var carPaginatedTestData = []Car{
	{Brand: "Toyota", Color: "Red", Year: 2020, Model: "Corolla"},
	{Brand: "Ford", Color: "Blue", Year: 2018, Model: "Focus"},
	{Brand: "Peugeot", Color: "White", Year: 2019, Model: "208"},
	{Brand: "Nissan", Color: "Green", Year: 2017, Model: "Sentra"},
	{Brand: "Mazda", Color: "Gray", Year: 2019, Model: "3"},
}

// Data for GetByField
var carByFieldTestData = []Car{
	{Brand: "Toyota", Color: "Red", Year: 2020, Model: "Corolla"},
	{Brand: "Peugeot", Color: "White", Year: 2019, Model: "208"},
	{Brand: "Ford", Color: "Blue", Year: 2018, Model: "Focus"},
}
