package types

type Status string

const (
	Created    Status = "Created"
	Processing Status = "Processing"
	Updated    Status = "Updated"
	Done       Status = "Done"
)
