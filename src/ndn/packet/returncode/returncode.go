package returncode

/* reference : https://tools.ietf.org/html/draft-irtf-icnrg-ccnxmessages-00#section-3 */
type Base int

const (
	NoRoute Base = 1 + iota
	HopLimitExceeded
	NoResource
	PathError
	Prohibited
	Congested
	MTUTooLarge
)
