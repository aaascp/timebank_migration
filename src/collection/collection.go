package collection

type Item interface {
	ToString() string
	ToDbFormat() map[string]interface{}
}
