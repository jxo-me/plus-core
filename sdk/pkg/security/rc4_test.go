package security

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func BenchmarkRc4Cipher_EncryptAndDecrypt(b *testing.B) {
	nonce := time.Now().Unix()
	var tests = []struct {
		key       string
		plaintext string
	}{
		{"159054a86e3bfb85b5f1991cdb07645e", fmt.Sprintf("这是测试rc4加密内容.#$5.*x,time:%d", nonce)},
	}
	for i := 0; i < b.N; i++ {
		rc4 := NewRc4Cipher(tests[0].key)
		encrypt, err := rc4.Encrypt(tests[0].plaintext)
		if err != nil {
			b.Error("RC4 Encrypt err", err)
		}
		decrypt, err := rc4.Decrypt(encrypt)
		if err != nil {
			b.Error("RC4 Decrypt err", err)
		}
		if string(decrypt) != tests[0].plaintext {
			b.Error(`RC4 Decrypt fail`)
		}
	}
}

func TestRc4EncryptDecrypt(t *testing.T) {
	nonce := time.Now().Unix()
	var tests = []struct {
		key       string
		plaintext string
	}{
		{"159054a86e3bfb85b5f1991cdb07645e", fmt.Sprintf("这是测试rc4加密内容.#$5.*x,time:%d", nonce)},
		{"159054a86e3bfb85b5f1991cdb07645e", fmt.Sprintf("这是测试rc4加密内容.#$5.*x,time:%d", nonce)},
		{"159054a86e3bfb85b5f1991cdb07645e", fmt.Sprintf("这是测试rc4加密内容.#$5.*x,time:%d", nonce)},
	}
	for _, test := range tests {
		rc4 := NewRc4Cipher(test.key)
		encrypt, err := rc4.Encrypt(test.plaintext)
		if err != nil {
			t.Error("RC4 Encrypt err", err)
		}
		decrypt, err := rc4.Decrypt(encrypt)
		if err != nil {
			t.Error("RC4 Decrypt err", err)
		}
		if string(decrypt) != test.plaintext {
			t.Error(`RC4 Decrypt fail`)
		}
	}
}
func TestRc4Decrypt(t *testing.T) {
	plaintext := `{"nickname":"","loginIp":469982874,"loginMacCode":"c4f174bd56e5c497254c390fd9223e0549d6ed5f","loginDevice":1,"loginSrc":4,"LoginAddress":"美国 加利福尼亚 洛杉矶 0","sex":"","timestamp":1653051576,"channelId":"800001","head":""}`
	keyStr := "sBymrOdodHQ4mLcnQsehtKWolnMaR0aLmtmFlCZNZsHwVHb2ZASCYW2kdiN7OLX1"
	encrypt := "9BB9F1C0E9267EA7E1B950F865262E67102D46F088F6678F801A9AFE00A346EE9FE7991ACA3D3AF1E60EF8CE5779AC1F39DE7E7286BB9F2174AB94182CB9705A4A64EB39FFE09CD81BC22791D63E0961DDA1AE4EDBA24F05FE3DFD8F4D66A2F7D593C513FD52DB24F62863365CC72EE3EDD98C7991A4FCD5895B9995F6BA0FBE5ED498DC355AA9196A12BFD1FBB16E1F0BD7E27E6454D81227CACC0B8405C13B3A0BD07FC88BBCFC4C3D810EF9233D069227D57844BC973A891EB6A8E7F5F24BF7E97DC1060977687BF18F774999C3F7CB4B1F314DDE9C0FBD19FF2CF922445FE26745360F5264968DA8D07956"
	c := NewRc4Cipher(keyStr)
	decrypt, err := c.decrypt(encrypt)
	if err != nil {
		t.Error("RC4 Decrypt err", err)
	}
	//err = ioutil.WriteFile("test.json", []byte(decrypt), 0644)
	//if err != nil {
	//	panic(err)
	//}
	if bytes.Compare(decrypt, []byte(plaintext)) != 0 {
		t.Error(`Decrypt fail`)
	}
	//fmt.Println(decrypt)
}

func TestRc4EncryptAndDecrypt(t *testing.T) {
	keyStr := "sBymrOdodHQ4mLcnQsehtKWolnMaR0aLmtmFlCZNZsHwVHb2ZASCYW2kdiN7OLX1"
	plaintext := "{\"UserName\":\"d1a3a79eb27920febcec43\",\"Password\":\"aa123456\",\"loginIp\":\"1275925876\",\"loginMacCode\":\"a598e745b7d9666ef5a550d8aef79cd2\",\"SMSCode\":\"12344\",\"loginDevice\":\"1\",\"RegDevice\":\"1\"}"
	c := NewRc4Cipher(keyStr)
	encrypt, err := c.encrypt(plaintext)
	if err != nil {
		t.Error("RC4 Encrypt err", err)
	}
	decrypt, err := c.decrypt(encrypt)
	if err != nil {
		t.Error("Rc4Decrypt error:", err)
	}
	if bytes.Compare(decrypt, []byte(plaintext)) != 0 {
		t.Error(`Decrypt fail`)
	}
}

func BenchmarkRc4EncryptAndDecrypt(b *testing.B) {
	keyStr := "sBymrOdodHQ4mLcnQsehtKWolnMaR0aLmtmFlCZNZsHwVHb2ZASCYW2kdiN7OLX1"
	plaintext := "{\"UserName\":\"d1a3a79eb27920febcec43\",\"Password\":\"aa123456\",\"loginIp\":\"1275925876\",\"loginMacCode\":\"a598e745b7d9666ef5a550d8aef79cd2\",\"SMSCode\":\"12344\",\"loginDevice\":\"1\",\"RegDevice\":\"1\"}"
	c := NewRc4Cipher(keyStr)
	for i := 0; i < b.N; i++ {
		encrypt, err := c.encrypt(plaintext)
		if err != nil {
			b.Error("RC4 Encrypt err", err)
		}
		decrypt, err := c.decrypt(encrypt)
		if err != nil {
			b.Error("RC4 Decrypt err", err)
		}
		if bytes.Compare(decrypt, []byte(plaintext)) != 0 {
			b.Error(`Decrypt fail`)
		}
	}
}

func TestRc4ClientEncryptAndDecrypt(t *testing.T) {
	plaintext := `{
    "version": "1.0.108",
    "msg": "",
    "ip": "127.0.0.1",
    "port": "13000",
    "sourceID": "CLMwZxSK",
    "sourceIDWeb": "sTyp4RgA",
    "shieldstatus": "false",
    "OpenAnWeishi":"true",
	"checkSimulator":"false",
    "clienturl": "http://tm9527.com/?channelCode=8001",
    "downloadurl": "http://taojin.kim",
    "apkurl": "http://dsndn.oss-cn-hongkong.aliyuncs.com/dsn.apk",
    "ipaurl": "itms-services://?action=download-manifest&url=https://tmylg.oss-cn-hangzhou.aliyuncs.com/tmassets/ipa/800001.plist",
    "recharge":"",
    "RegisterWebURL" : "http://127.0.0.1:13088/Register",
    "CheckCodeAgoWebURL" : "",
    "OnQueryAccWebURL" : "",
    "SetAccPasswordWebURL" : "http://127.0.0.1:13088/ResetPwd",
    "PayWebURL" : "",
    "BindAccount" : "http://127.0.0.1:13088/BindPhone",
    "QRCodeURL" : "",
    "SendMsgOnPcWebURL" : "http://127.0.0.1:13088/SendSms",
    "CheckMsgOnPcWebURL" : "",
    "HOST" : "https://jhgg.oss-accelerate.aliyuncs.com/",
    "TgHttpUrl" : "http://127.0.0.1:13088/",
    "kefuUrl" :  "https://talk.nanfengmao1.com/kf.html",
    "NOTICE" : "http://127.0.0.1:13088/GetServiceInfo?p=C5ACFD8CB87F35F2E8F944F362323160487104AF809A21999F1899E20BA951E5C9F6870A83616ABDBB72BC9B2633FB493ED7693481E98B2576ECC61B3FEF7D1C4D69F938FCF1C9D807C470D0D9391E618EB4AD4FC7F54D45AC3DBE9B5361F8AB97C99E72AB149774A62F6F3655D67BA9BC82C719D6E2E8DD98418EFCAFEE43E62995CF9D751BE81E62537D5D4771C1C60E040EFBB7EE54C3B8601FD81DD74E8428DF41D74F330729F6B8901BE83A3A46CF7AC47243AD883DCE41FEFEA1A4A01EA2F874D5030F616B7EE588725A8782B19118546B4D97E352BA15EE31FF70515BB6605F25420161D799A7D769197A9752FF7A23DFF260EAFF9CD3823BC2279C82F702EA2B71A99B45DDBB37DE9474C0D62E61E9B20DD924CB8A4903702D981CFD666A443526F5E0902B76EC84A4FB56DD31898C3FEB8733785AA0C5BB13F8A385B130072EB1B9BCA30B559EB21ADE59BCB37EB6DE1D0E05A855AF00971858A87C7C4364D5031FBF47B1A4E856F3B4F85249133ED97FFFAF3C54D19F54E96CB8858EE158DA57C7297C6D50E4F09207CB2983D78386656927A4E2F10F482B55ED1C9C1FBED0EE2CEB3B9150AA457CD50557CD7A16433793BB4327",
	"channel": "http://127.0.0.1:13087/GetChannelId",
	"getip":"http://api.ip138.com/ipv4/?token=22593ac90f78a4001d7350148bf1350d",
	"debugMacCodes":"2d8fa5bf61e18ee38e170cfc031c63c4,"
}`
	c := NewRc4Cipher("")
	encrypt, err := c.Rc4ClientEncrypt(plaintext)
	if err != nil {
		t.Error("RC4 Client Encrypt err", err)
	}
	decrypt, err := c.Rc4ClientDecrypt(c.CiphertextReplace(c.CiphertextFormat(encrypt)))
	if err != nil {
		t.Error("RC4 Client Decrypt err", err)
	}
	if decrypt != plaintext {
		t.Error(`Client Decrypt fail`)
	}
}

func TestRc4ClientDecrypt(t *testing.T) {
	plaintext := `{
    "version": "1.0.108",
    "msg": "",
    "ip": "127.0.0.1",
    "port": "13000",
    "sourceID": "CLMwZxSK",
    "sourceIDWeb": "sTyp4RgA",
    "shieldstatus": "false",
    "OpenAnWeishi":"true",
	"checkSimulator":"false",
    "clienturl": "http://tm9527.com/?channelCode=8001",
    "downloadurl": "http://taojin.kim",
    "apkurl": "http://dsndn.oss-cn-hongkong.aliyuncs.com/dsn.apk",
    "ipaurl": "itms-services://?action=download-manifest&url=https://tmylg.oss-cn-hangzhou.aliyuncs.com/tmassets/ipa/800001.plist",
    "recharge":"",
    "RegisterWebURL" : "http://127.0.0.1:13088/Register",
    "CheckCodeAgoWebURL" : "",
    "OnQueryAccWebURL" : "",
    "SetAccPasswordWebURL" : "http://127.0.0.1:13088/ResetPwd",
    "PayWebURL" : "",
    "BindAccount" : "http://127.0.0.1:13088/BindPhone",
    "QRCodeURL" : "",
    "SendMsgOnPcWebURL" : "http://127.0.0.1:13088/SendSms",
    "CheckMsgOnPcWebURL" : "",
    "HOST" : "https://jhgg.oss-accelerate.aliyuncs.com/",
    "TgHttpUrl" : "http://127.0.0.1:13088/",
    "kefuUrl" :  "https://talk.nanfengmao1.com/kf.html",
    "NOTICE" : "http://127.0.0.1:13088/GetServiceInfo?p=C5ACFD8CB87F35F2E8F944F362323160487104AF809A21999F1899E20BA951E5C9F6870A83616ABDBB72BC9B2633FB493ED7693481E98B2576ECC61B3FEF7D1C4D69F938FCF1C9D807C470D0D9391E618EB4AD4FC7F54D45AC3DBE9B5361F8AB97C99E72AB149774A62F6F3655D67BA9BC82C719D6E2E8DD98418EFCAFEE43E62995CF9D751BE81E62537D5D4771C1C60E040EFBB7EE54C3B8601FD81DD74E8428DF41D74F330729F6B8901BE83A3A46CF7AC47243AD883DCE41FEFEA1A4A01EA2F874D5030F616B7EE588725A8782B19118546B4D97E352BA15EE31FF70515BB6605F25420161D799A7D769197A9752FF7A23DFF260EAFF9CD3823BC2279C82F702EA2B71A99B45DDBB37DE9474C0D62E61E9B20DD924CB8A4903702D981CFD666A443526F5E0902B76EC84A4FB56DD31898C3FEB8733785AA0C5BB13F8A385B130072EB1B9BCA30B559EB21ADE59BCB37EB6DE1D0E05A855AF00971858A87C7C4364D5031FBF47B1A4E856F3B4F85249133ED97FFFAF3C54D19F54E96CB8858EE158DA57C7297C6D50E4F09207CB2983D78386656927A4E2F10F482B55ED1C9C1FBED0EE2CEB3B9150AA457CD50557CD7A16433793BB4327",
	"channel": "http://127.0.0.1:13087/GetChannelId",
	"getip":"http://api.ip138.com/ipv4/?token=22593ac90f78a4001d7350148bf1350d",
	"debugMacCodes":"2d8fa5bf61e18ee38e170cfc031c63c4,"
}`
	ciphertext := `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
7l5ogVMK8IywtL3Ks5tAs1LCOMVUseFnZOJ+jg==`
	c := NewRc4Cipher("")
	decrypt, err := c.Rc4ClientDecrypt(c.CiphertextReplace(ciphertext))
	if err != nil {
		t.Error("RC4 Client Decrypt err", err)
	}
	err = os.WriteFile("test.json", []byte(decrypt), 0644)
	if err != nil {
		panic(err)
	}
	//fmt.Println([]byte(plaintext))
	if decrypt != plaintext {
		t.Error(`Client Client Decrypt fail`)
	}
}
