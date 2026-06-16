package usenix

// Paper is a paper published at a USENIX conference.
type Paper struct {
	Rank       int    `json:"rank"`
	Conference string `json:"conference"`
	Title      string `json:"title"`
	URL        string `json:"url"`
}

// Conference is a known USENIX conference.
type Conference struct {
	Rank int    `json:"rank"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Year int     `json:"year"`
}

// KnownConferences is the hardcoded list of known USENIX conferences.
var KnownConferences = []Conference{
	{1, "usenixsecurity24", "USENIX Security '24", 2024},
	{2, "usenixsecurity23", "USENIX Security '23", 2023},
	{3, "usenixsecurity22", "USENIX Security '22", 2022},
	{4, "usenixsecurity21", "USENIX Security '21", 2021},
	{5, "usenixsecurity20", "USENIX Security '20", 2020},
	{6, "nsdi24", "NSDI '24", 2024},
	{7, "nsdi23", "NSDI '23", 2023},
	{8, "nsdi22", "NSDI '22", 2022},
	{9, "osdi24", "OSDI '24", 2024},
	{10, "osdi23", "OSDI '23", 2023},
	{11, "osdi22", "OSDI '22", 2022},
	{12, "fast24", "FAST '24", 2024},
	{13, "fast23", "FAST '23", 2023},
	{14, "fast22", "FAST '22", 2022},
	{15, "atc24", "ATC '24", 2024},
	{16, "atc23", "ATC '23", 2023},
	{17, "atc22", "ATC '22", 2022},
}

// Conferences returns the list of known USENIX conferences.
func Conferences() []Conference { return KnownConferences }
