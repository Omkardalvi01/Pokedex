package main

type Pokemon struct {
	Name      string    `json:"name"`
	Height    int       `json:"height"`
	Weight    int       `json:"weight"`
	Abilities []Ability `json:"abilities"`
	MoveSet   []Moves   `json:"moves"`
	Types     []Type    `json:"types"`
	Stats     []Stat    `json:"stats"`
	Imgpaths  ImgPath   `json:"sprites"`
}

type Ability struct {
	Details struct {
		Name string `json:"name"`
	} `json:"ability"`
}

type Type struct {
	Details struct {
		Name string `json:"name"`
	} `json:"type"`
}

type Moves struct {
	Move struct {
		Name string `json:"name"`
	} `json:"move"`
}

type Stat struct {
	Base int `json:"base_stat"`
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Physical struct {
	Height int `json:"height"`
	Weight int `json:"weight"`
}

type ImgPath struct {
	Front string `json:"front_default"`
	Back  string `json:"back_default"`
}
