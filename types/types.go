package types

type Credentials struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Role_code string `json:"role_code"`
}


type RoleCode struct {
	Role string `json:"role" db:"role"`
	Code string `json:"code" db:"code"`
}

type Secret struct {
	Password  string `json:"password" db:"password"`
	UserId string `json:"userid" db:"userid"`
}

type User struct {
	UserId    string `json:"userid" db:"userid"`
	Username  string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Role      string `json:"role" db:"role"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type UserRole struct {
	UserId string `json:"userid" db:"userid"`
	TeamId string `json:"team_id" db:"team_id"`
}

type Token struct {
	Token_id string `json:"token_id" db:"token_id"`
	UserId string `json:"userid" db:"userid"`
	Role string `json:"role" db:"role"`
	Type string `json:"token_type" db:"token__type"`
	Exp string `json:"expires_at" db:"expires_at"`
	Iat string `json:"issues_at" db:"issued_at"`
	Status bool `json:"status" db:"status"`
}