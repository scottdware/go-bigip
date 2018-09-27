package bigip

const (
	uriClientSSL       = "client-ssl"
	uriDatagroup       = "data-group"
	uriHttp            = "http"
	uriHttpCompression = "http-compression"
	uriIRule           = "rule"
	uriInternal        = "internal"
	uriLtm             = "ltm"
	uriMonitor         = "monitor"
	uriNode            = "node"
	uriOneConnect      = "one-connect"
	uriPolicy          = "policy"
	uriPool            = "pool"
	uriPoolMember      = "members"
	uriProfile         = "profile"
	uriServerSSL       = "server-ssl"
	uriSnatPool        = "snatpool"
	uriTcp             = "tcp"
	uriUdp             = "udp"
	uriVirtual         = "virtual"
	uriVirtualAddress  = "virtual-address"
	uriGtm             = "gtm"
	uriWideIp          = "wideip"
	uriARecord         = "a"
	uriAAAARecord      = "aaaa"
	uriCNameRecord     = "cname"
	uriMXRecord        = "mx"
	uriNaptrRecord     = "naptr"
	uriSrvRecord       = "srv"
	ENABLED            = "enable"
	DISABLED           = "disable"
	CONTEXT_SERVER     = "serverside"
	CONTEXT_CLIENT     = "clientside"
	CONTEXT_ALL        = "all"

	// Newer policy APIs have a draft-publish workflow that this library does not support.
	policyVersionSuffix = "?ver=11.5.1"
)

var cidr = map[string]string{
	"0":  "0.0.0.0",
	"1":  "128.0.0.0",
	"2":  "192.0.0.0",
	"3":  "224.0.0.0",
	"4":  "240.0.0.0",
	"5":  "248.0.0.0",
	"6":  "252.0.0.0",
	"7":  "254.0.0.0",
	"8":  "255.0.0.0",
	"9":  "255.128.0.0",
	"10": "255.192.0.0",
	"11": "255.224.0.0",
	"12": "255.240.0.0",
	"13": "255.248.0.0",
	"14": "255.252.0.0",
	"15": "255.254.0.0",
	"16": "255.255.0.0",
	"17": "255.255.128.0",
	"18": "255.255.192.0",
	"19": "255.255.224.0",
	"20": "255.255.240.0",
	"21": "255.255.248.0",
	"22": "255.255.252.0",
	"23": "255.255.254.0",
	"24": "255.255.255.0",
	"25": "255.255.255.128",
	"26": "255.255.255.192",
	"27": "255.255.255.224",
	"28": "255.255.255.240",
	"29": "255.255.255.248",
	"30": "255.255.255.252",
	"31": "255.255.255.254",
	"32": "255.255.255.255",
}

type gtmType string

const (
	aRecord     gtmType = uriARecord
	aaaaRecord  gtmType = uriAAAARecord
	cnameRecord gtmType = uriCNameRecord
	mxRecord    gtmType = uriMXRecord
	naptrRecord gtmType = uriNaptrRecord
	srvRecord   gtmType = uriSrvRecord
)
