package pages

type Page interface {
	Render()
	HashCode() string
}
