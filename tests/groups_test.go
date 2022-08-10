// SPDX-License-Group: MIT
//
// Copyright (C) 2021 Daniel Bourdrez. All Rights Reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree or at
// https://spdx.org/licenses/MIT.html

package group_test

import (
	"encoding/hex"
	"testing"

	"github.com/bytemare/crypto"
	"github.com/bytemare/crypto/internal/nist"
	"github.com/bytemare/crypto/internal/ristretto"
)

type testGroup struct {
	name          string
	h2c           string
	e2c           string
	id            crypto.Group
	elementLength uint
	scalarLength  uint
}

func testGroups() []*testGroup {
	return []*testGroup{
		{"Ristretto255", ristretto.H2C, ristretto.H2C, crypto.Ristretto255Sha512, 32, 32},
		{"P256", nist.H2CP256, nist.E2CP256, crypto.P256Sha256, 33, 32},
		{"P384", nist.H2CP384, nist.E2CP384, crypto.P384Sha384, 49, 48},
		{"P521", nist.H2CP521, nist.E2CP521, crypto.P521Sha512, 67, 66},
	}
}

func testAllGroups(t *testing.T, f func(*testing.T, *testGroup)) {
	for _, test := range testGroups() {
		t.Run(test.name, func(t *testing.T) {
			f(t, test)
		})
	}
}

func TestAvailability(t *testing.T) {
	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		if !group.id.Available() {
			t.Errorf("'%s' is not available, but should be", group.id.String())
		}
	})
}

func TestNonAvailability(t *testing.T) {
	oob := crypto.Group(0)
	if oob.Available() {
		t.Errorf("%v is considered available when it must not", oob)
	}

	d := crypto.Group(2) // decaf448
	if d.Available() {
		t.Errorf("%v is considered available when it must not", d)
	}

	oob = crypto.P521Sha512 + 1
	if oob.Available() {
		t.Errorf("%v is considered available when it must not", oob)
	}
}

func TestGroup_Base(t *testing.T) {
	base := map[crypto.Group]string{
		crypto.Ristretto255Sha512: "e2f2ae0a6abc4e71a884a961c500515f58e30b6aa582dd8db6a65945e08d2d76",
		crypto.P256Sha256:         "036b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296",
		crypto.P384Sha384:         "03aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7",
		crypto.P521Sha512:         "0200c6858e06b70404e9cd9e3ecb662395b4429c648139053fb521f828af606b4d3dbaa14b5e77efe75928fe1dc127a2ffa8de3348b3c1856a429bf97e7e31c2e5bd66",
	}

	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		if hex.EncodeToString(group.id.Base().Encode()) != base[group.id] {
			t.Fatalf("Got wrong base element %s", hex.EncodeToString(group.id.Base().Encode()))
		}
	})
}

func TestDST(t *testing.T) {
	app := "app"
	version := uint8(1)
	tests := map[crypto.Group]string{
		crypto.Ristretto255Sha512: "app-V01-CS01-ristretto255_XMD:SHA-512_R255MAP_RO_",
		crypto.P256Sha256:         "app-V01-CS03-P256_XMD:SHA-256_SSWU_RO_",
		crypto.P384Sha384:         "app-V01-CS04-P384_XMD:SHA-384_SSWU_RO_",
		crypto.P521Sha512:         "app-V01-CS05-P521_XMD:SHA-512_SSWU_RO_",
	}

	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		res := string(group.id.MakeDST(app, version))
		test := tests[group.id]
		if res != test {
			t.Errorf("Wrong DST. want %q, got %q", res, test)
		}
	})
}

func TestGroup_String(t *testing.T) {
	tests := map[crypto.Group]string{
		crypto.Ristretto255Sha512: "ristretto255_XMD:SHA-512_R255MAP_RO_",
		crypto.P256Sha256:         "P256_XMD:SHA-256_SSWU_RO_",
		crypto.P384Sha384:         "P384_XMD:SHA-384_SSWU_RO_",
		crypto.P521Sha512:         "P521_XMD:SHA-512_SSWU_RO_",
	}

	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		res := group.id.String()
		test := tests[group.id]
		if res != test {
			t.Errorf("Wrong DST. want %q, got %q", res, test)
		}
	})
}

func TestGroup_NewScalar(t *testing.T) {
	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		s := group.id.NewScalar().Encode()
		for _, b := range s {
			if b != 0 {
				t.Fatalf("expected zero scalar, but got %v", hex.EncodeToString(s))
			}
		}
	})
}

func TestGroup_NewElement(t *testing.T) {
	identity := map[crypto.Group]string{
		crypto.Ristretto255Sha512: "0000000000000000000000000000000000000000000000000000000000000000",
		crypto.P256Sha256:         "00",
		crypto.P384Sha384:         "00",
		crypto.P521Sha512:         "00",
	}

	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		e := hex.EncodeToString(group.id.NewElement().Encode())
		ref := identity[group.id]

		if e != ref {
			t.Fatalf("expected identity element %v, but got %v", ref, e)
		}
	})
}

func TestGroup_ScalarLength(t *testing.T) {
	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		if group.id.ScalarLength() != group.scalarLength {
			t.Fatalf("expected encoded scalar length %d, but got %d", group.scalarLength, group.id.ScalarLength())
		}
	})
}

func TestGroup_ElementLength(t *testing.T) {
	testAllGroups(t, func(t2 *testing.T, group *testGroup) {
		if group.id.ElementLength() != group.elementLength {
			t.Fatalf("expected encoded element length %d, but got %d", group.elementLength, group.id.ElementLength())
		}
	})
}
