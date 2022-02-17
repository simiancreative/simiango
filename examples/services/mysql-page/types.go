package mysqlpage

type Service struct {
	Pattern  string `param:"pattern"`
	Pattern2 string `param:"pattern2"`
	Page     int    `param:"page"`
	Size     int    `param:"size"`
}

type Product struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
