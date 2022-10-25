package mongodbr

type IEntityBeforeCreate interface {
	BeforeCreate()
}

type IEntityBeforeUpdate interface {
	BeforeUpdate()
}
