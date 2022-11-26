package ignorepattern

type Access int // want Access:"^Standard,User,Group$"

const Standard Access = User

const (
	User  Access = 1
	Group Access = 2
)

// The member User is ignored by the -ignore-enum-members flag.
// The member Standard, though it has the same constant value as User, must
// still be reported in the diagnostic.
func _c(a Access) {
	switch a { // want "^missing cases in switch of type ignorepattern.Access: ignorepattern.Standard, ignorepattern.Group$"
	}

	_ = map[Access]int{ // want "^missing keys in map of key type ignorepattern.Access: ignorepattern.Standard, ignorepattern.Group$"
		0: 0,
	}
}
