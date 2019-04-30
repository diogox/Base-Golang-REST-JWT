package authentication

const FREE_USER_ROLE = "freeUser"
const PREMIUM_USER_ROLE = "PremiumUser"

func ResolveUserRole(isPaidUser bool) string {
	if isPaidUser {
		return PREMIUM_USER_ROLE
	}

	return FREE_USER_ROLE
}