package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cyphar/paperback/pkg/shamir"
)

func main() {
	secret := []byte("this is a secret which is a bit larger")

	shards, err := shamir.Split(3, 6, secret)
	if err != nil {
		fmt.Printf("error creating shards: %v", err)
		return
	}

	fmt.Printf("our secret is: %v\n", string(secret))
	for _, shard := range shards {
		shardJson, _ := json.Marshal(shard)
		fmt.Printf("shard: %+v\n", string(shardJson))
	}

	fmt.Printf("combine1\n")
	if secret2, err := shamir.Combine(shards...); err != nil {
		fmt.Printf("ERROR COMBINE: %v\n", err)
	} else if !bytes.Equal(secret, secret2) {
		fmt.Printf("ERROR COMBINE MISMATCH: %q != %q\n", string(secret), string(secret2))
	}

	fmt.Printf("combine2\n")
	if secret2, err := shamir.Combine(shards[:3]...); err != nil {
		fmt.Printf("ERROR COMBINE: %v\n", err)
	} else if !bytes.Equal(secret, secret2) {
		fmt.Printf("ERROR COMBINE MISMATCH: %q != %q\n", string(secret), string(secret2))
	}

	fmt.Printf("combine3\n")
	if secret2, err := shamir.Combine(shards[2:]...); err != nil {
		fmt.Printf("ERROR COMBINE: %v\n", err)
	} else if !bytes.Equal(secret, secret2) {
		fmt.Printf("ERROR COMBINE MISMATCH: %q != %q\n", string(secret), string(secret2))
	}
}
