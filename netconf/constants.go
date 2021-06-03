package netconf

// constants to represent netconf capabilities and patterns.
const (
	Version10             = "1.0"
	Version10Capability   = "urn:ietf:params:netconf:base:1.0"
	Version10Capabilities = "" +
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<hello xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">\n" +
		"     <capabilities>\n" +
		"         <capability>urn:ietf:params:netconf:base:1.0</capability>\n" +
		"     </capabilities>\n" +
		"</hello>]]>]]>"
	Version10DelimiterPattern = "]]>]]>"
	Version11                 = "1.1"
	Version11Capability       = "urn:ietf:params:netconf:base:1.1"
	Version11Capabilities     = "" +
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<hello xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">\n" +
		"     <capabilities>\n" +
		"         <capability>urn:ietf:params:netconf:base:1.1</capability>\n" +
		"     </capabilities>\n" +
		"</hello>]]>]]>"
	Version11DelimiterPattern = `(?m)^##$`
	Version11ChunkPattern     = `(?ms)(\d+)\n(.*?)#`
	XMLHeader                 = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	DefaultPort               = 830
)
