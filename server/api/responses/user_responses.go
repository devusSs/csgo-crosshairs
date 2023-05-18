// Responses for successful requests against any user route.
// Naming conventions should tell which route is in question.
package responses

// Login user.
type LoginUserResponse struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

// This struct serves any route except Login, Logout, GetUser.
type GeneralUserResponse struct {
	Message string `json:"message"`
}
