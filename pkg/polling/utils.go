package polling

type Link struct {
	Self string
	Git  string
	HTML string
}

type RepoContent struct {
	Name         string
	Path         string
	Sha          string
	URL          string
	HTML_url     string
	Git_url      string
	Download_url string
	Type         string
	Links        Link
}

type YAMLContent struct {
	Name         string
	Path         string
	Sha          string
	Size         int
	URL          string
	HTML_url     string
	Git_url      string
	Download_url string
	Type         string
	Links        Link
}
