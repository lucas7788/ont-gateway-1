// +build !test,!prod,!k8sdev

package config

import "time"

var config = Value{
	RestConfig: RestConfig{
		PublicAddr:      ":8080",
		IntraAddr:       ":8081",
		GracefulUpgrade: false,
		PIDFile:         "/tmp/ont-gateway.pid",
		ReadTimeout:     time.Second * 3,
	},
	LoggerConfig: LoggerConfig{
		LogLevel: "debug",
	},
	RedisCacheConfig: RedisConfig{
		Addr: "172.168.3.46:6379",
	},
	MongoConfig: MongoConfig{
		ConnectionString: "mongodb://172.168.3.47:27017/ont",
		Timeout:          time.Second * 3,
	},
	CICDConfig: CICDConfig{
		AddonDeployAPI: AkSkURL{
			Host: "a0d771952588111ea89590659513bb5d-1585432770.ap-southeast-1.elb.amazonaws.com:8000",
			URI:  "/api/v1/ss",
		},
	},
	EsignConfig: EsignConfig{
		DocuConfig: DocuConfig{
			IntegratorKey: "b68ab0de-9f7f-4567-a2cf-40161f2ac975",
			PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAgC2C+FDB0D5nNuZGeg/HBhExMPnVaYyvQivTaqh2LOn6jd4X
0cqNXcp4pV7fq5wOeqNXozNd++pT6xc/e530kqRTPTrHjTvzyeYyIyvFX02Oh5I/
hW66qtH52HFqt7ZnFGTdkwub2O53uKvFGZ8I+yhsDUW2l8dRNgkDxQBhxE0CUVwx
u9Ei0wkEEZM+3sVBkzys0mhnlKhqKPJ2AO0JLr53ULf1iHZvUobXbYHN3gznzO2G
JZvj1CtYEGUcIck1rvQ0nUXSBvJNumqvM2TxYLBZnuCA4v0bynLfcrKzidtWezIK
ZHJWhMSXJIDwn7XIuL8FUQxrXlv5ai4FmA+nzwIDAQABAoIBAA0em203Nt9Nw9rG
rygHPWPNlq9gowturvGi8rzUCWSHfnHO7bk3dkjHVJn4oAQ8sO60kV/O/iuuzHAf
rQvGGyZ3U13NCmfWXxmnSjJ3ZHhgw9n0ijPqJYkefOg+k6HCNcLMoDFQ2t7VacOp
MT9yG8U5WfXx0MKwUAur2IrtgxRId6iL31v1geA7hAgib+DqwWqE5RbP6acmfz/c
ETVoATNyUWbOF6gj1oaeBMHr6uQmZi3+ha94OWiaUBn0mZqFJYGjcL4HK9KHnVaz
Y1jiDYVDOxBQX2DyGrXH7rZ4T6JB0WCSLHAgWJ00xnaZyTToK2wKw9hOq6eFAISZ
iiOhSQECgYEA4GmbZl+Y7epkggr/K2CvCwBwuke+MtGzNcOHM9Bum/oy1HwSLw8Y
2umtW/EUdlwdveBBHLfK3TUOn5uuBFqzFGiZ7N0y1h4nuGnyhpsK80MKc0SUCdpj
6bqM75zA+1JoGjORAeAbAwoEmxIJ1lcEf/+SWmmXUnL0VrpAg4kCi2ECgYEAkjg4
LE059+8APhuz+vChCfMvfhQDy8WMFmLZlhe3uY2PkptE5iQNLqYK90pUu8xhQnyZ
V6UoKFOPLw63Qyo0z/TzwkZpnuDf9bP4xHGtP++Em+ALckBCp3UBs3IlY4xsOjRC
WXz7271oEZ6+jUgumZ8j0WiBncwV9uQOh4J6sS8CgYB0vduOpRqcYfwJPolB2pkU
4xTBg0Lpkvdkd6QlC0APOlgo+6ZF/teSQk/h7YcUj5UVSsz0kJQjAU/rLgSX1Usl
yciRVPz2MFe/crYs2gkXRX/xOPK+MXMaiuZ4XBZ0Z4kqYDsGO7wxl4uP1BF0BG6d
26kaCaYjyRNc7qVTB/pf4QKBgEgFqx2fOHN4ZP0ythdf2WLGR1lp0GjZuGP6csSs
kBG0uchz9J8LmPamUPZ3xX7vb+TI7Nsv/bTHW9rI+9n4eyUHud2ywynACHDFIj7Y
44Z/mykQVXMEVhCX4KucCPCc5V5SCXB80K3vAMjVEXUT3ehLa+AlAttQAG2o7cMv
sPF/AoGAFc5ib0ipFwmdNMjE6yUmhAMTTT9CivWxes+GASfS95dSopEIiBMtsPw9
7Ls6473xH7sWLlxQiekQYHJz7m5Ph3TIM9E502sckZhWVNgOZvqcQiNai6ONb54j
WO8HMzH4VrSOYudVNdIlbP3TZx73sXQ70K+YQfSW2ZREdk+iOmw=
-----END RSA PRIVATE KEY-----`,
			KeyPairID:    "a503663f-6d3f-4875-aa75-1ac6b5a8019c",
			APIAccountID: "ef99483c-2abc-4fe7-9941-5960ea835c70",
			APIUserID:    "ecc89097-2d5c-4135-bca4-4f09ed86ca8a",
		},
	},
}
