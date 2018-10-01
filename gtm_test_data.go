package bigip

func wideIPSamples() []byte {
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
										"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool?ver=12.1.1"
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

func poolASamples() []byte {
	return []byte(
		`{
			"kind": "tm:gtm:pool:a:acollectionstate",
			"selfLink": "https://localhost/mgmt/tm/gtm/pool/a?ver=12.1.1",
			"items": [
					{
							"kind": "tm:gtm:pool:a:astate",
							"name": "baseapp.domain.com_pool",
							"partition": "Common",
							"fullPath": "/Common/baseapp.domain.com_pool",
							"generation": 2,
							"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool?ver=12.1.1",
							"alternateMode": "round-robin",
							"dynamicRatio": "disabled",
							"enabled": true,
							"fallbackIp": "any",
							"fallbackMode": "return-to-dns",
							"limitMaxBps": 0,
							"limitMaxBpsStatus": "disabled",
							"limitMaxConnections": 0,
							"limitMaxConnectionsStatus": "disabled",
							"limitMaxPps": 0,
							"limitMaxPpsStatus": "disabled",
							"loadBalancingMode": "round-robin",
							"manualResume": "disabled",
							"maxAnswersReturned": 1,
							"monitor": "default",
							"qosHitRatio": 5,
							"qosHops": 0,
							"qosKilobytesSecond": 3,
							"qosLcs": 30,
							"qosPacketRate": 1,
							"qosRtt": 50,
							"qosTopology": 0,
							"qosVsCapacity": 0,
							"qosVsScore": 0,
							"ttl": 30,
							"verifyMemberAvailability": "enabled",
							"membersReference": {
									"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1",
									"isSubcollection": true
							}
					},
					{
            "kind": "tm:gtm:pool:a:astate",
            "name": "myapp.domain.com_pool",
            "partition": "test",
            "fullPath": "/test/myapp.domain.com_pool",
            "generation": 182,
            "selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1",
            "alternateMode": "round-robin",
            "dynamicRatio": "disabled",
            "enabled": true,
            "fallbackIp": "any",
            "fallbackMode": "return-to-dns",
            "limitMaxBps": 0,
            "limitMaxBpsStatus": "disabled",
            "limitMaxConnections": 0,
            "limitMaxConnectionsStatus": "disabled",
            "limitMaxPps": 0,
            "limitMaxPpsStatus": "disabled",
            "loadBalancingMode": "round-robin",
            "manualResume": "disabled",
            "maxAnswersReturned": 1,
            "monitor": "default",
            "qosHitRatio": 5,
            "qosHops": 0,
            "qosKilobytesSecond": 3,
            "qosLcs": 30,
            "qosPacketRate": 1,
            "qosRtt": 50,
            "qosTopology": 0,
            "qosVsCapacity": 0,
            "qosVsScore": 0,
            "ttl": 30,
            "verifyMemberAvailability": "enabled",
            "membersReference": {
                "link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members?ver=12.1.1",
                "isSubcollection": true
            }
          }
      ]
    }`)
}

func poolASample(usePartition bool) []byte {
	if usePartition {
		return []byte(`{
			"kind": "tm:gtm:pool:a:astate",
			"name": "myapp.domain.com_pool",
			"partition": "test",
			"fullPath": "/test/myapp.domain.com_pool",
			"generation": 182,
			"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1",
			"alternateMode": "round-robin",
			"dynamicRatio": "disabled",
			"enabled": true,
			"fallbackIp": "any",
			"fallbackMode": "return-to-dns",
			"limitMaxBps": 0,
			"limitMaxBpsStatus": "disabled",
			"limitMaxConnections": 0,
			"limitMaxConnectionsStatus": "disabled",
			"limitMaxPps": 0,
			"limitMaxPpsStatus": "disabled",
			"loadBalancingMode": "round-robin",
			"manualResume": "disabled",
			"maxAnswersReturned": 1,
			"monitor": "default",
			"qosHitRatio": 5,
			"qosHops": 0,
			"qosKilobytesSecond": 3,
			"qosLcs": 30,
			"qosPacketRate": 1,
			"qosRtt": 50,
			"qosTopology": 0,
			"qosVsCapacity": 0,
			"qosVsScore": 0,
			"ttl": 30,
			"verifyMemberAvailability": "enabled",
			"membersReference": {
					"link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members?ver=12.1.1",
					"isSubcollection": true
			}
		}`)
	}

	return []byte(`{
		"kind": "tm:gtm:pool:a:astate",
		"name": "baseapp.domain.com_pool",
		"partition": "Common",
		"fullPath": "/Common/baseapp.domain.com_pool",
		"generation": 2,
		"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool?ver=12.1.1",
		"alternateMode": "round-robin",
		"dynamicRatio": "disabled",
		"enabled": true,
		"fallbackIp": "any",
		"fallbackMode": "return-to-dns",
		"limitMaxBps": 0,
		"limitMaxBpsStatus": "disabled",
		"limitMaxConnections": 0,
		"limitMaxConnectionsStatus": "disabled",
		"limitMaxPps": 0,
		"limitMaxPpsStatus": "disabled",
		"loadBalancingMode": "round-robin",
		"manualResume": "disabled",
		"maxAnswersReturned": 1,
		"monitor": "default",
		"qosHitRatio": 5,
		"qosHops": 0,
		"qosKilobytesSecond": 3,
		"qosLcs": 30,
		"qosPacketRate": 1,
		"qosRtt": 50,
		"qosTopology": 0,
		"qosVsCapacity": 0,
		"qosVsScore": 0,
		"ttl": 30,
		"verifyMemberAvailability": "enabled",
		"membersReference": {
				"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1",
				"isSubcollection": true
		}
	}`)
}

func poolAReturn(usePartiion bool) string {
	if usePartiion {
		return `{"name":"myapp.domain.com_pool","partition":"test","fullPath":"/test/myapp.domain.com_pool","generation":182,"dynamicRatio":"disabled","enabled":true,"fallbackIp":"any","fallbackMode":"return-to-dns","limitMaxBpsStatus":"disabled","limitMaxConnectionsStatus":"disabled","limitMaxPpsStatus":"disabled","loadBalancingMode":"round-robin","manualResume":"disabled","maxAnswersReturned":1,"monitor":"default","qosHitRatio":5,"qosKilobytesSecond":3,"qosLcs":30,"qosPacketRate":1,"qosRtt":50,"ttl":30,"verifyMemberAvailability":"enabled","MembersReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members?ver=12.1.1","isSubcollection":true}}`
	}

	return `{"name":"baseapp.domain.com_pool","partition":"Common","fullPath":"/Common/baseapp.domain.com_pool","generation":2,"dynamicRatio":"disabled","enabled":true,"fallbackIp":"any","fallbackMode":"return-to-dns","limitMaxBpsStatus":"disabled","limitMaxConnectionsStatus":"disabled","limitMaxPpsStatus":"disabled","loadBalancingMode":"round-robin","manualResume":"disabled","maxAnswersReturned":1,"monitor":"default","qosHitRatio":5,"qosKilobytesSecond":3,"qosLcs":30,"qosPacketRate":1,"qosRtt":50,"ttl":30,"verifyMemberAvailability":"enabled","MembersReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1","isSubcollection":true}}`
}

func poolAMemberSamples() []byte {
	return []byte(
		`{
			"kind": "tm:gtm:pool:a:members:memberscollectionstate",
			"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1",
			"items": [
				{
						"kind": "tm:gtm:pool:a:members:membersstate",
						"name": "baseapp_80_vs",
						"partition": "Common",
						"subPath": "someltm:/Common",
						"fullPath": "/Common/someltm:/Common/baseapp_80_vs",
						"generation": 197,
						"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members/~Common~someltm:~Common~baseapp_80_vs?ver=12.1.1",
						"enabled": true,
						"limitMaxBps": 0,
						"limitMaxBpsStatus": "disabled",
						"limitMaxConnections": 0,
						"limitMaxConnectionsStatus": "disabled",
						"limitMaxPps": 0,
						"limitMaxPpsStatus": "disabled",
						"memberOrder": 0,
						"monitor": "default",
						"ratio": 1
				},
				{
					"kind": "tm:gtm:pool:a:members:membersstate",
					"name": "baseapp_443_vs",
					"partition": "Common",
					"subPath": "someltm:/Common",
					"fullPath": "/Common/someltm:/Common/baseapp_443_vs",
					"generation": 197,
					"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members/~Common~someltm:~Common~baseapp_443_vs?ver=12.1.1",
					"enabled": true,
					"limitMaxBps": 0,
					"limitMaxBpsStatus": "disabled",
					"limitMaxConnections": 0,
					"limitMaxConnectionsStatus": "disabled",
					"limitMaxPps": 0,
					"limitMaxPpsStatus": "disabled",
					"memberOrder": 0,
					"monitor": "default",
					"ratio": 1
				},
				{
					"kind": "tm:gtm:pool:a:members:membersstate",
					"name": "myapp_80_vs",
					"partition": "Common",
					"subPath": "someltm:/test",
					"fullPath": "/Common/someltm:/test/myapp_80_vs",
					"generation": 197,
					"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members/~Common~someltm:~test~myapp_80_vs?ver=12.1.1",
					"enabled": true,
					"limitMaxBps": 0,
					"limitMaxBpsStatus": "disabled",
					"limitMaxConnections": 0,
					"limitMaxConnectionsStatus": "disabled",
					"limitMaxPps": 0,
					"limitMaxPpsStatus": "disabled",
					"memberOrder": 0,
					"monitor": "default",
					"ratio": 1
				},
				{
					"kind": "tm:gtm:pool:a:members:membersstate",
					"name": "myapp_443_vs",
					"partition": "Common",
					"subPath": "someltm:/test",
					"fullPath": "/Common/someltm:/test/myapp_443_vs",
					"generation": 197,
					"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members/~Common~someltm:~test~myapp_443_vs?ver=12.1.1",
					"enabled": true,
					"limitMaxBps": 0,
					"limitMaxBpsStatus": "disabled",
					"limitMaxConnections": 0,
					"limitMaxConnectionsStatus": "disabled",
					"limitMaxPps": 0,
					"limitMaxPpsStatus": "disabled",
					"memberOrder": 0,
					"monitor": "default",
					"ratio": 1
				}
			]
		}`)
}
