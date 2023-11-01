package pkg

import (
	"testing"
)

func TestOpenIdWithKeyEncode(t *testing.T) {
	var tests = []struct {
		Key    string
		userId int
		result string
	}{
		{"8a6f2805b4515ac12058e79e66539be9", 1, "2fb9GnkjEMA4qlb7BRxz3nuRG50_2BgJ_2BbP91T6V4T"},
		{"8a6f2805b4515ac12058e79e66539be9", 2, "2fb9GnkjEMA4qlb7BUop2X7AG508icmdP9xd7FAQ"},
		{"8a6f2805b4515ac12058e79e66539be9", 3, "2fb9GnkjEMA4qlb7BRB6jiyRScFqhJ3HbN1Q7lER"},
		{"8a6f2805b4515ac12058e79e66539be9", 4, "2fb9GnkjEMA4qlb7BRx9iCmWFMpuhcrKbNJd6wIW"},
		{"8a6f2805b4515ac12058e79e66539be9", 5, "2fb9GnkjEMA4qlb7BU8tji2QT5xr15vOOthT7QMX"},
	}
	for _, test := range tests {
		open := NewOpenIdWithKey(test.Key)
		openId, err := open.Encode(test.userId)
		if err != nil {
			t.Error("open.Encode error:", err)
		}
		t.Log(openId)
		if openId != test.result {
			t.Error(`open.Encode fail`, test.Key, openId)
		}
	}
}

func TestOpenIdWithKeyDecode(t *testing.T) {
	var tests = []struct {
		Key    string
		userId string
		result string
	}{
		{"8a6f2805b4515ac12058e79e66539be9", "1", "2fb9GnkjEMA4qlb7BRxz3nuRG50_2BgJ_2BbP91T6V4T"},
		{"8a6f2805b4515ac12058e79e66539be9", "2", "2fb9GnkjEMA4qlb7BUop2X7AG508icmdP9xd7FAQ"},
		{"8a6f2805b4515ac12058e79e66539be9", "3", "2fb9GnkjEMA4qlb7BRB6jiyRScFqhJ3HbN1Q7lER"},
		{"8a6f2805b4515ac12058e79e66539be9", "4", "2fb9GnkjEMA4qlb7BRx9iCmWFMpuhcrKbNJd6wIW"},
		{"8a6f2805b4515ac12058e79e66539be9", "5", "2fb9GnkjEMA4qlb7BU8tji2QT5xr15vOOthT7QMX"},
	}
	for _, test := range tests {
		open := NewOpenIdWithKey(test.Key)
		userId, err := open.Decode(test.result)
		if err != nil {
			t.Error("open.Encode error:", err)
		}
		t.Log(userId)
		if userId != test.userId {
			t.Error(`open.Encode fail`, test.Key, userId)
		}
	}
}

func TestOpenIdEncode(t *testing.T) {
	var tests = []struct {
		KeyA   string
		KeyB   string
		Salt   string
		userId int
		result string
	}{
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", 1, "2727fS7bOhnpJXQR9zWiSrNox3unuZJyae_2B5NwuS"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", 2, "2727fS7bOhnpJXQR9zT0Trhiw3ml6MByaOLlNwCR"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", 3, "2727fS7bOhnpJXQR9zagT7Nrk3nz75skP7_2B4NluQ"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", 4, "2727fS7bOhnpJXQR92D0TbJslSig5sUjarrtMQCX"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", 5, "2727fS7bOhnpJXQR9zb2HOU4lCT37pZ8beq5YAiW"},
	}
	for _, test := range tests {
		open := NewOpenId(test.KeyA, test.KeyB, test.Salt)
		openId, err := open.Encode(test.userId)
		if err != nil {
			t.Error("open.Encode error:", err)
		}
		t.Log(openId)
		if openId != test.result {
			t.Error(`open.Encode fail`, test.KeyA, test.KeyB, openId)
		}
	}
}

func TestOpenIdDecode(t *testing.T) {
	var tests = []struct {
		KeyA   string
		KeyB   string
		Salt   string
		userId string
		result string
	}{
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", "1", "2727fS7bOhnpJXQR9zWiSrNox3unuZJyae_2B5NwuS"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", "2", "2727fS7bOhnpJXQR9zT0Trhiw3ml6MByaOLlNwCR"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", "3", "2727fS7bOhnpJXQR9zagT7Nrk3nz75skP7_2B4NluQ"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", "4", "2727fS7bOhnpJXQR92D0TbJslSig5sUjarrtMQCX"},
		{"8a6f2805b4515ac12058e79e66539be9", "1929f008c8edbcf5756f95bf876466e7", "50f0b71156b779e16825b4e8df502727", "5", "2727fS7bOhnpJXQR9zb2HOU4lCT37pZ8beq5YAiW"},
	}
	for _, test := range tests {
		open := NewOpenId(test.KeyA, test.KeyB, test.Salt)
		userId, err := open.Decode(test.result)
		if err != nil {
			t.Error("open.Encode error:", err)
		}
		t.Log(userId)
		if userId != test.userId {
			t.Error(`open.Encode fail`, test.KeyA, test.KeyB, userId)
		}
	}
}
