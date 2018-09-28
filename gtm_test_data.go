package bigip

func wideIPsSample() []byte {
	return []byte(`{
		"kind": "tm:gtm:wideip:a:acollectionstate",
		"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a?ver=12.1.1",
		"items": [
			{
				"kind": "tm:gtm:wideip:a:astate",
				"name": "baseapp.domain.com",
				"partition": "Common",
				"fullPath": "/Common/baseapp.domain.com",
				"generation": 2,
				"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a/~Common~baseapp.domain.com?ver=12.1.1",
				"enabled": true,
				"failureRcode": "noerror",
				"failureRcodeResponse": "disabled",
				"failureRcodeTtl": 0,
				"lastResortPool": "",
				"minimalResponse": "enabled",
				"persistCidrIpv4": 32,
				"persistCidrIpv6": 128,
				"persistence": "disabled",
				"poolLbMode": "topology",
				"ttlPersistence": 3600,
				"pools": [
						{
								"name": "baseapp.domain.com_pool",
								"partition": "Common",
								"order": 0,
								"ratio": 1,
								"nameReference": {
										"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool_int_pool?ver=12.1.1"
								}
						}
				]
			},
			{
				"kind": "tm:gtm:wideip:a:astate",
				"name": "myapp.domain.com",
				"partition": "test",
				"fullPath": "/test/myapp.domain.com",
				"generation": 35,
				"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a/~test~myapp.domain.com?ver=12.1.1",
				"enabled": true,
				"failureRcode": "noerror",
				"failureRcodeResponse": "disabled",
				"failureRcodeTtl": 0,
				"lastResortPool": "",
				"minimalResponse": "enabled",
				"persistCidrIpv4": 32,
				"persistCidrIpv6": 128,
				"persistence": "disabled",
				"poolLbMode": "round-robin",
				"ttlPersistence": 3600,
				"pools": [
						{
								"name": "myapp.domain.com.com_pool",
								"partition": "test",
								"order": 0,
								"ratio": 1,
								"nameReference": {
										"link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1"
								}
						}
				]
			}
		]
	}`)
}

func wideIPSample(usePartition bool) []byte {
	if usePartition {
		return []byte(`{
			"kind": "tm:gtm:wideip:a:astate",
			"name": "myapp.domain.com",
			"partition": "test",
			"fullPath": "/test/myapp.domain.com",
			"generation": 35,
			"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a/~test~myapp.domain.com?ver=12.1.1",
			"enabled": true,
			"failureRcode": "noerror",
			"failureRcodeResponse": "disabled",
			"failureRcodeTtl": 0,
			"lastResortPool": "",
			"minimalResponse": "enabled",
			"persistCidrIpv4": 32,
			"persistCidrIpv6": 128,
			"persistence": "disabled",
			"poolLbMode": "round-robin",
			"ttlPersistence": 3600,
			"pools": [
					{
							"name": "myapp.domain.com.com_pool",
							"partition": "test",
							"order": 0,
							"ratio": 1,
							"nameReference": {
									"link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1"
							}
					}
			]
		}`)
	}

	return []byte(`{
		"kind": "tm:gtm:wideip:a:astate",
		"name": "baseapp.domain.com",
		"partition": "Common",
		"fullPath": "/Common/baseapp.domain.com",
		"generation": 2,
		"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a/~Common~baseapp.domain.com?ver=12.1.1",
		"enabled": true,
		"failureRcode": "noerror",
		"failureRcodeResponse": "disabled",
		"failureRcodeTtl": 0,
		"lastResortPool": "",
		"minimalResponse": "enabled",
		"persistCidrIpv4": 32,
		"persistCidrIpv6": 128,
		"persistence": "disabled",
		"poolLbMode": "topology",
		"ttlPersistence": 3600,
		"pools": [
				{
						"name": "baseapp.domain.com_pool",
						"partition": "Common",
						"order": 0,
						"ratio": 1,
						"nameReference": {
								"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool_int_pool?ver=12.1.1"
						}
				}
		]
	}`)
}

func wideIPReturn(usePartiion bool) string {
	if usePartiion {
		return `{"name":"myapp.domain.com","partition":"test","fullPath":"/test/baseapp.domain.com","generation":2,"enabled":true,"failureRcode":"noerror","failureRcodeResponse":"disabled","minimalResponse":"enabled","persistCidrIpv4":32,"persistCidrIpv6":128,"persistence":"disabled","poolLbMode":"topology","ttlPersistence":3600,"pools":[{"name":"myapp.domain.com_pool","partition":"test","ratio":1,"nameReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool_int_pool?ver=12.1.1"}}]}`
	}

	return `{"name":"baseapp.domain.com","partition":"Common","fullPath":"/Common/baseapp.domain.com","generation":2,"enabled":true,"failureRcode":"noerror","failureRcodeResponse":"disabled","minimalResponse":"enabled","persistCidrIpv4":32,"persistCidrIpv6":128,"persistence":"disabled","poolLbMode":"topology","ttlPersistence":3600,"pools":[{"name":"baseapp.domain.com_pool","partition":"Common","ratio":1,"nameReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool_int_pool?ver=12.1.1"}}]}`
}
