package domain
import "time"

// UserProfile defines model for UserProfile.
type UserProfile struct {
	UserID                 int64                      `json:"userId" db:"id"`
	Name                   string                   `json:"name" db:"name"`
	UserName               string                   `json:"userName" db:"username"`
	LastName               string                   `json:"lastName" db:"last_name"`
	Email                  string                   `json:"email" db:"email"`
	Identification         Identification           `json:"identification"`
	SocialNetwork          []SocialNetwork          `json:"socialNetwork,omitempty"`
	Gender                 Gender                   `json:"gender"`
	Address                Address                  `json:"address,omitempty"`
	Birthday               time.Time                `json:"birthday" time_format:"2006-01-02"`
	Nationality            Nationality              `json:"nationality"`
	ProfessionalExperience []ProfessionalExperience `json:"professionalExperience,omitempty"`
	AcademicFormation      []AcademicFormation      `json:"academicFormation"`
	UserContactPhoneNumber []Telephone              `json:"userContactPhoneNumber,omitempty"`
	Locale                 string                   `json:"locale"`
	Timezone               string                   `json:"timezone"`
	Picture                string                   `json:"picture,omitempty"`
	PublicProfile          PublicProfile            `json:"publicProfile,omitempty"`
	UserBlocked            bool                     `json:"userBlocked,omitempty"`
}

// Identification defines model for Identification.
type Identification struct {
	IdentificationID   string    `json:"identificationId"`
	IdentificationType string    `json:"identificationType"`
	IssuingCountry     string    `json:"issuingCountry,omitempty"`
	IssuingDate        time.Time `json:"issuingDate,omitempty" time_format:"2006-01-02"`
}

// SocialNetwork defines model for SocialNetwork.
type SocialNetwork struct {
	SocialNetworkName string `json:"socialNetworkName"`
	UserName          string `json:"userName"`
	URL               string `json:"url"`
}

// Gender defines model for Gender.
type Gender struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

// Address defines model for Address.
type Address struct {
	StreetName      string   `json:"streetName"`
	StreetNumber    string   `json:"streetNumber"`
	ZipCode         string   `json:"zipCode"`
	Department      string   `json:"department"`
	Additionals     string   `json:"additionals,omitempty"`
	Country         string   `json:"country"`
	StateOrProvince string   `json:"stateOrProvince"`
	City            string   `json:"city"`
	Locality        Locality `json:"locality"`
	Geo             Geo      `json:"geo,omitempty"`
}

// Locality defines model for Locality.
type Locality struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// Geo defines model for Geo.
type Geo struct {
	Coordinates Point  `json:"coordinates"`
	Type        string `json:"type"`
}

// Point defines model for Point.
type Point struct {
	Latitud  string `json:"latitud"`
	Longitud string `json:"longitud"`
}

// Nationality defines model for Nationality.
type Nationality struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// ProfessionalExperience defines model for ProfessionalExperience.
type ProfessionalExperience struct {
	AdditionalInformation string    `json:"additionalInformation"`
	Address               Address   `json:"address"`
	Company               string    `json:"company"`
	Description           string    `json:"description"`
	EndDate               time.Time `json:"endDate,omitempty" time_format:"2006-01-02"`
	Position              string    `json:"position"`
	StartingDate          time.Time `json:"startingDate" time_format:"2006-01-02"`
}

// AcademicFormation defines model for AcademicFormation.
type AcademicFormation struct {
	Additional   *string                  `json:"additional,omitempty"`
	Career       *string                  `json:"career,omitempty"`
	EndDate      *time.Time               `json:"endDate,omitempty"`
	Institution  *string                  `json:"institution,omitempty"`
	StartingDate *time.Time               `json:"startingDate,omitempty"`
	Status       *AcademicFormationStatus `json:"status,omitempty"`
	Type         *AcademicEducationType   `json:"type,omitempty"`
	University   *University              `json:"university,omitempty"`
}

// AcademicFormationStatus defines model for AcademicFormationStatus.
type AcademicFormationStatus struct {
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"`
}

// AcademicEducationType defines model for AcademicEducationType.
type AcademicEducationType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// University defines model for University.
type University struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Telephone defines model for Telephone.
type Telephone struct {
	AreaCode    int32  `json:"areaCode"`
	CountryCode int32  `json:"countryCode"`
	PhoneNumber int32  `json:"phoneNumber"`
	Type        string `json:"type"`
}

// PublicProfile defines model for PublicProfile.
type PublicProfile *struct {
	CaaUserProfile *string `json:"caaUserProfile,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
}