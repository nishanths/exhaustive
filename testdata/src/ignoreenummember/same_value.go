package ignoreenummember

type Access int // want Access:"^Standard,User,Group$"

const Standard Access = User

const (
	User  Access = 1
	Group Access = 2
)

// Though the member User is ignored by the -ignore-enum-members flag,
// Standard ,which has the same value as User, must still be reported
// in the diagnostic.
func _c(a Access) {
	switch a { // want "^missing cases in switch of type Access: Group, Standard$"
	}
}
