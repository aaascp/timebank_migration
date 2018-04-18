package collection

type Collection interface {
	ToString() string
	ToDbFormat() map[string]interface{}
}
