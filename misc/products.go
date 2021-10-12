package misc

type Product struct {
	ProductID   int
	ProductName string
}

func GetProducts() []Product {
	return []Product{
		{1111, "RTX 3090"},
		{2222, "GTX 1080"},
		{3333, "Intel i9 9900k"},
	}
}

type Project struct {
	ChargeNumber int
	ProjectName  string
}

func GetProjects() []Project {
	return []Project{
		{2011, "database"},
		{2022, "operating system"},
		{2023, "compiler"},
		{2024, "web application"},
	}
}
