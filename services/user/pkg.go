package user

import "github.com/samber/do/v2"

var Package = do.Package(
	do.Lazy(newUserService),
)
