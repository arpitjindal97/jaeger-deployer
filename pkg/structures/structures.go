package structures

// TenantPayload is the struct accepted by TenantPut
type TenantPayload struct {
	Customers           []string
	Domain              string
	StorageType         string
	ESHost              string
	CassandraHost       string
	CassandraDataCenter string
	KafkaEnabled        bool
	KafkaBroker         string
	HotrodExample       bool
	AuthType            string
}

// Tenant struct hold info about a particular Tenant
type Tenant struct {
	Customer           string `json:"customer"`
	JaegerCollectorURL string `json:"jaegerCollectorURL,omitempty"`
	JaegerQueryURL     string `json:"jaegerQueryURL,omitempty"`
	Username           string `json:"username,omitempty"`
	Password           string `json:"password,omitempty"`
	CACrt              string `json:"caCrt,omitempty"`
	CAKey              string `json:"caKey,omitempty"`
	ClientCrt          string `json:"clientCrt,omitempty"`
	ClientKey          string `json:"clientKey,omitempty"`
	Error              string `json:"error,omitempty"`
}

// PerformancePayload for accepting JSON
type PerformancePayload struct {
	Customers []string
	Domain    string
}

// Response struct is used for responding in JSON
type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
