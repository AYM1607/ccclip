package main

import (
	"encoding/base64"
)

func b64(i []byte) string {
	return base64.StdEncoding.EncodeToString(i)
}

var (
	priv1 = "~/dev/ccclip/keys1/private.key"
	pub1  = "~/dev/ccclip/keys1/public.key"

	priv2 = "~/dev/ccclip/keys2/private.key"
	pub2  = "~/dev/ccclip/keys2/public.key"
)

func main() {
	rootCmd.Execute()
	// key1 := crypto.LoadPrivateKey("../keys1/private.key")
	// key2 := crypto.LoadPrivateKey("../keys2/private.key")

	// secretMsg := "new-some-secret-messageeee"

	// encrypted := crypto.Encrypt(
	// 	crypto.NewSharedKey(key1, key2.PublicKey(), crypto.SendDirection),
	// 	[]byte(secretMsg),
	// )

	// fmt.Printf("Message %q was encrypted to %q\n", secretMsg, b64(encrypted))

	// decrypted := crypto.Decrypt(
	// 	crypto.NewSharedKey(key2, key1.PublicKey(), crypto.ReceiveDirection),
	// 	encrypted,
	// )

	// fmt.Printf("Message was decrypted as %q\n", string(decrypted))
}
