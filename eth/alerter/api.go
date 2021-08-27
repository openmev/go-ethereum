package alerter


// PublicAlerterAPI allows a party to register new alert means
// and to create new patterns for which it should be alerted
type PublicAlerterAPI struct {
	alertDestinations map[string]bool // this is used as a set
}


// NewPublicAlerterAPI create a new PublicAlerterAPI.
func NewPublicAlerterAPI() *PublicAlerterAPI {
	return &PublicAlerterAPI{
		alertDestinations: make(map[string]bool),
	}
}


// RegisterDestination registers a new destination to which to send
// the alert when one is triggered
// TODO: use persistent storage
func (api *PublicAlerterAPI) RegisterDestination(destination string) (bool, error) {
	if _, ok := api.alertDestinations[destination]; ok {
		return false, nil
	}
	api.alertDestinations[destination] = true
	return true, nil
}
