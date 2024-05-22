package data

type Transformer struct{}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) Transform(records []map[string]interface{}) []map[string]interface{} {
	// Implement your transformation logic here
	return records
}
