package test

import (
	"github.com/small-ek/ginp/crypto/rsa"
	"log"
	"testing"
)
/*var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgQChfYIOm7bOOlHF2VcEwqGaR1OxVT8fpJIjmvz6+k2EbiMo5e6b
L1pyRUXEVyXf+x06bJtXC1QCOKtwmnUKPAsSrWSP5jYmyLCjNaTJqwcCgyAXMA7k
AVoLRfjTxt4XYSm+P4HS+9w1yLZvsZ1XbRIhxgVHDtknkgANt7OREpBYowIDAQAB
An8BHMtcai9AZCJ9cQDcDnswQI0KnNAwxeCXIqrRD8yUsNBFwEnLtAnKucdBQqHb
cxIC9MaAZtoGmwIMZEl2BxmRYMhYZZyHdrKBNZJJgmV8Cz3Q1sRC12KJ1LMlSziV
8NO5tbKUOKzkY9GFItjgriDDfd1SAjNE+B7ydbN5cFKBAkEA/HTW3Xqt7Y8qxM51
4Za0RF1NQmVuyspKsaVcRCE/Igeq+T/pZ0jWvo/BFGDrFME4BTFbLJsP5URjZatH
wVzkIwJBAKPBy/xxzrpCwRVTMP2PkouqtUFYi8vYWK1j0DVPvNNkikYugOC+h6OQ
lbPiFspGgA1O2/rRqGzDT0fGdAjhwYECQQCwjeHKiMZkchB+DMmiF6xAd2PVwGxI
REsSi8vIFdw6J1Sp9cl8oxMTuCNW5iThofNUplzWCCeItlgxPST0lMszAkAhY9un
DsGbQw9BvOPJX+P+rIEm4NooZ2W1fRuwMyEKbX6wTr0illbr6AhOVHRXLEbh78l0
/Bj+jFh3ByUTxoyBAkAGZjTBaWJuQuNnD3Do+CuHAoJ3k9drAWEEAepQdzIIom0m
7+Vc8wQsdZSLK+ATqIM/KkC00dR1462axPXUR6f3
-----END RSA PRIVATE KEY-----
`)
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChfYIOm7bOOlHF2VcEwqGaR1Ox
VT8fpJIjmvz6+k2EbiMo5e6bL1pyRUXEVyXf+x06bJtXC1QCOKtwmnUKPAsSrWSP
5jYmyLCjNaTJqwcCgyAXMA7kAVoLRfjTxt4XYSm+P4HS+9w1yLZvsZ1XbRIhxgVH
DtknkgANt7OREpBYowIDAQAB
-----END PUBLIC KEY-----
`)*/
func TestRSA(t *testing.T)  {
	/*var rsa=rsa.Default()*/
	for i:=0;i<100;i++{
		var str1,_= rsa.Default().Encrypt("admin");
		log.Println(str1)
		log.Println(rsa.Default().Decrypt(str1))
	}
}