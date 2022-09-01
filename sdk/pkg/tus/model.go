package tus

type GetByIdInput struct {
	Id int64 `v:"required" json:"id" description:"主键Id"`
}

type NullStructRes struct{}
