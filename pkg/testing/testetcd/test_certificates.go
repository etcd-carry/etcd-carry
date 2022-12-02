package testetcd

import (
	"io/ioutil"
	"path"
	"testing"
)

const CACertContent = `
-----BEGIN CERTIFICATE-----
MIIC0jCCAbqgAwIBAgIJALkNCS8h/g0KMA0GCSqGSIb3DQEBCwUAMBUxEzARBgNV
BAMMCmt1YmVybmV0ZXMwIBcNMjIwNzE0MDU0MjUyWhgPMjEyMjA2MjAwNTQyNTJa
MBUxEzARBgNVBAMMCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDiww1p3DABhjaflUVNh516KlpUNBe4cXUsfJPtVB7ZVdd39v2OGqW5
pKYiwIaYCVrHeMfQn8xDBehgi4yru2ijXuPsy1WhwffowfG+ojTmv+HK4uk8dEK9
2yEx4AGiS8CeBTbJ+nLknY//wFaePij/dRi/cP625zRpOEUCYDUWqWi7wUSiEpRq
YBDIzIoVrwbv934fxlP1wKX4g7WpYO52mg2OPDVwY9nMxRdUSGBtRXWg/+vrdkyv
IiNosKU1GA+GAZ5h1UEknc2SHfNvihlOlZ6VsEKya54v7hStEA16w3AQsZ/HiU4I
FAWAyrncTksgYGt6EuWlDPTmrB+SXJ63AgMBAAGjIzAhMA8GA1UdEwEB/wQFMAMB
Af8wDgYDVR0PAQH/BAQDAgKkMA0GCSqGSIb3DQEBCwUAA4IBAQAaOtGaD6e+FPi0
VKTHI6NTuAfiw7Uubq58/YJtph0h2nQzHSjfrHJxg6SkKWpLX/IGWnCpP+z8A+yE
7c8ny5xU1WqPvpBEQUf4S5W4tNR5euZMGkqWzU3KYVeqO+ji+nefy/cDB9TF7FwB
U9AZezsrZLh4cFFsU7sD7rgVWoZwuuCJqI6B0bwK0usVNZAL+FWbaMOIcTXLxiDM
8r59KLfs5E0WEnT8BUIA84SIzW+WSRMSHwhz/0bFHGXCyWluSVFRgdmb0+qFbqoi
nda0/6sDBk1QgG/xggDLtgEZUqo+B1hPddzcSkYQCzeUTD8Hbj4Bx/KyjMe0SocK
BAt5fltE
-----END CERTIFICATE-----
`

const CertContent = `
-----BEGIN CERTIFICATE-----
MIIDIjCCAgqgAwIBAgIJAOGbndYSAhgXMA0GCSqGSIb3DQEBCwUAMBUxEzARBgNV
BAMMCmt1YmVybmV0ZXMwIBcNMjIwNzE0MDU1MTI3WhgPMjEyMjA2MjAwNTUxMjda
MBExDzANBgNVBAMMBm1hc3RlcjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAMhu4yKy0k+HcBHI0vlsXBIRgtBTfeh2hHXihyu00iX8/G+8az95NDRHMMro
nYKzqTiUy1sXdXC+zsqH8/aD/XNhuZHVoBsxJaKIOC+XY2XyxTzMTpovgWh9Z5D6
Qm9SXoUyA4pBxJL8N95xQHqZRHu+CUGQnUzsWU8usREX/OYdV57w/UQfnkph68G8
rq1F1eUzQYWc4Ua8kKuH3GU4Lr2TKJ/xA1IDwktip+usCfoNIjvYcy1z8AncbBWN
209x3sy1JSmhFyZfPz0HgHeaEbKy8/kcPt50BPKqykQTzNDJrD13ctp5DWggqP0T
gmxoeyefqY1CghJpCDe6qk0QYR0CAwEAAaN3MHUwCQYDVR0TBAIwADAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMDkGA1UdEQQy
MDCCCWxvY2FsaG9zdIIFbm9kZTGHBH8AAAGHEAAAAAAAAAAAAAAAAAAAAAGHBMCo
MNwwDQYJKoZIhvcNAQELBQADggEBAKtzbxJlU3qXe0m3pkogzL2UYu3bF8M21Thw
VLjRUvVeimcKV7t9YCjQhbTzL9MFtQWhuJrrgWEFBp7oZBXAuQBmiNPgjGYYbAS2
rgTFEmjVuJs8bXXAbGZbLxgWptKVTIHJedzFosk1ZwPQJSyYHiXZlG405zCReblC
9VopsiDzDIhipoUJuEL+YtfyhytMASsOpmZx1skcC+pNVOrHTohH0xwo7+i1SG5r
L0VUy9WGAnWS0mxSamXr/hTQO4nyQSja3nNZmFDufjx26pLlS1xzLfrhGLknhviS
M+F0LE79mBQtf90wdA/Ifh14/YZ2Pn6+8CkgFridQyJ24uKqCVk=
-----END CERTIFICATE-----
`

const KeyContent = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAyG7jIrLST4dwEcjS+WxcEhGC0FN96HaEdeKHK7TSJfz8b7xr
P3k0NEcwyuidgrOpOJTLWxd1cL7Oyofz9oP9c2G5kdWgGzEloog4L5djZfLFPMxO
mi+BaH1nkPpCb1JehTIDikHEkvw33nFAeplEe74JQZCdTOxZTy6xERf85h1XnvD9
RB+eSmHrwbyurUXV5TNBhZzhRryQq4fcZTguvZMon/EDUgPCS2Kn66wJ+g0iO9hz
LXPwCdxsFY3bT3HezLUlKaEXJl8/PQeAd5oRsrLz+Rw+3nQE8qrKRBPM0MmsPXdy
2nkNaCCo/ROCbGh7J5+pjUKCEmkIN7qqTRBhHQIDAQABAoIBAQCjBU3IkxlrhgUo
8eZm6DOanpN/TelCgeMK3syCR/gE3deUKfQxFCvZuW1+G+TAwdAJLTyZQmNK6GbZ
Y+qRvpkOl5WPf+lRNJAfuhu42bEG4oZ6BNKJpcnjatwpluMiGCS4wQ0QDp4LzwjB
6+s9zBtfahmtMio+vp2FQbzS4mfX4OPBzZib/++dErdKTfJr7N46uWITUEL8iZNg
wSZ/bZKfX3LZEu2BxRilOKQMc1z2mAbpwcnoxhuISxFDsGySyX3nUxyQJ3gkJYdJ
48LQhkw/y41FqSGhbLINy/RDZFqslwkzHcer0/jh/TY7y9wdvwBhYKlAp05pHVIt
oTkZtRoBAoGBAPfcC94KiKi6owV7IBpQfAJwI9AQx7aU38VAJOoWpPBny6XJ8asL
D6u6DsoYsztIO1RedccnfuMzMO82fasWGkDzE22rsT00orLn5Br+i1nWmMf+/MAM
NlZydSn70QavPPqG3x10ud8o+zEtGtC5TRuzmEN9xYFmKULbkR9I59yRAoGBAM8E
FtPL7uNJk9jN59y9Juj5nXaLOJIGSkuUyGHDjEeISw3q/G6AHJKKotqwAvGcKhue
3vJcLJYf++8kSj3PRVMn0Fdij8wOkEVBmw48gB1M/AZ44cth3f8O1pnK0Bid30WV
+vxUXOpKGwQyx9KoQddm9fFehJN6XE1RdwGR3THNAoGBAJEcWjJQFnw6cLEHyd6+
GixPPRhWiqZBeNUB3drTERPSoO7aUUujeTRABOKbHWvquRmHCAtl+yrHULHsRBzD
HvHBnjFKIMVFqK93hhurxSf+tIn6pj1FqRZpgmDnFhSEyf2esseLDDszwgSjdJyY
sCU0u0NgQh6lEikbZVZcl6qRAoGBALvR9fQLHp5QbzdQ+YCojNjrQBYBkj3KPzX7
syIgDPIJki76eDS5PzMlXUQUVVdoXDvbFGPHhRxfwG/j+QfDOh6MDNZ7sgNtYy+y
qj9sXMA4zKACpLml/Ygfqky2Wb873QqBXMn6sKJQwdo5SFq0Faic3Z80Jgy4A26S
7uoZsRoxAoGBANmgEqzpDafeKiU1wqWCiA04+ZA+NfxaYyEUNqi9h5pcAvfd7woT
mNv4jj4lC/xMsYEQ3fdxQStHU2U/+QkysPZr/Lc6crB6IGRl3fasoZG/OIcCUlmt
d+1cOfShH7929IlDE9MkPhf8UlEIWiLDx6oJu2KDc64vHt3F2pku9KC5
-----END RSA PRIVATE KEY-----
`

func GetTLSCerts(t *testing.T, certsDir string) (caFile, certFile, keyFile string) {
	caFile = path.Join(certsDir, "ca.crt")
	if err := ioutil.WriteFile(caFile, []byte(CACertContent), 0644); err != nil {
		t.Fatal(err)
	}
	certFile = path.Join(certsDir, "server.crt")
	if err := ioutil.WriteFile(certFile, []byte(CertContent), 0644); err != nil {
		t.Fatal(err)
	}
	keyFile = path.Join(certsDir, "server.key")
	if err := ioutil.WriteFile(keyFile, []byte(KeyContent), 0644); err != nil {
		t.Fatal(err)
	}
	return caFile, certFile, keyFile
}
