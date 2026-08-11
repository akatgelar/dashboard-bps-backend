package models
  
type BaseFilter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"` // eq, neq, gt, gte, lt, lte, like
    Value    interface{} `json:"value"`
}