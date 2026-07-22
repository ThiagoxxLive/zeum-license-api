package cryptox

import "testing"

// Vetores reais extraídos do banco de dev do zeum-admin-api e decriptados via
// Defuse\Crypto\Crypto::decrypt() em PHP, para garantir compatibilidade byte a byte.
func TestDecryptMatchesPHPDefuse(t *testing.T) {

	secret := "def00000b32f61e77f1cd261cdedcec968633b98a09c3ac534eae77d3ed1b2a062112ed2bc99801bfb97f6a68c46368486b4047988caa98bcf5794623cf5d5799e1af572"
	ciphertext := "def5020059ed6d8942274b7274ece650a3cd31c860488223208916e3079c8588f02316d40928dc0edfd36a302f3bf409c78a7757c7ad184ae6ab63ad34d344423b5ee0bb6982afe6cc68e19011edf698be385dc78e4696c30e6ef04432b01aca16913dd1533ebc64353dfa8ba12c007757e6bb4f908a38cf3f064b"
	want := "zu-cd196fa7-0bb8-4278-a440-6d8f6ceb93bd"

	rawKey, err := ParseKey(secret)

	if err != nil {
		t.Fatalf("ParseKey retornou erro: %v", err)
	}

	got, err := Decrypt(ciphertext, rawKey)

	if err != nil {
		t.Fatalf("Decrypt retornou erro: %v", err)
	}

	if got != want {
		t.Fatalf("Decrypt = %q, esperado %q", got, want)
	}
}

func TestDecryptFailsWithWrongKey(t *testing.T) {

	secret := "def00000b32f61e77f1cd261cdedcec968633b98a09c3ac534eae77d3ed1b2a062112ed2bc99801bfb97f6a68c46368486b4047988caa98bcf5794623cf5d5799e1af572"
	ciphertext := "def5020059ed6d8942274b7274ece650a3cd31c860488223208916e3079c8588f02316d40928dc0edfd36a302f3bf409c78a7757c7ad184ae6ab63ad34d344423b5ee0bb6982afe6cc68e19011edf698be385dc78e4696c30e6ef04432b01aca16913dd1533ebc64353dfa8ba12c007757e6bb4f908a38cf3f064b"

	rawKey, err := ParseKey(secret)

	if err != nil {
		t.Fatalf("ParseKey retornou erro: %v", err)
	}

	rawKey[0] ^= 0xFF // corrompe a chave para simular uma API Key de outra aplicação

	if _, err := Decrypt(ciphertext, rawKey); err == nil {
		t.Fatal("esperava erro de integridade com a chave errada, mas decriptou com sucesso")
	}
}
