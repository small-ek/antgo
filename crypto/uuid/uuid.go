package uuid

import (
	"github.com/google/uuid"
	"log"
	"os"
)

// UUID representation compliant with specification
// described in RFC 4122.
type UUID = uuid.UUID

// UUID DCE domains.
const (
	DomainPerson = uuid.Domain(0)
	DomainGroup  = uuid.Domain(1)
	DomainOrg    = uuid.Domain(2)
)

// New creates a new random UUID or panics.
// It returns a Random (Version 4) UUID.
func New() UUID {
	return uuid.New()
}

// Create NewUUID returns a Version 1 UUID based on the current NodeID and clock
// sequence, and the current time.  If the NodeID has not been set by SetNodeID
// or SetNodeInterface then it will be set automatically.  If the NodeID cannot
// be set NewUUID returns nil.  If clock sequence has not been set by
// SetClockSequence then it will be set automatically.  If GetTime fails to
// return the current NewUUID returns nil and an error.
//
// In most cases, New should be used.
func Create() UUID {
	var result, err = uuid.NewUUID()
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

// NewDCEGroup returns a DCE Security (Version 2) UUID in the group
// domain with the id returned by os.Getgid.
//
//  NewDCESecurity(Group, uint32(os.Getgid()))
func NewDCEGroup() UUID {
	var result, err = uuid.NewDCESecurity(DomainGroup, uint32(os.Getgid()))
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

// NewDCEPerson returns a DCE Security (Version 2) UUID in the person
// domain with the id returned by os.Getuid.
//
//  NewDCESecurity(Person, uint32(os.Getuid()))
func NewDCEPerson() UUID {
	var result, err = uuid.NewDCESecurity(DomainPerson, uint32(os.Getuid()))
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

// NewMD5 returns a new MD5 (Version 3) UUID based on the
// supplied name space and data.  It is the same as calling:
//
//  NewHash(md5.New(), space, data, 3)
func NewMD5(space UUID, data []byte) UUID {
	return uuid.NewMD5(space, data)
}

// NewRandom returns a Random (Version 4) UUID.
//
// The strength of the UUIDs is based on the strength of the crypto/rand
// package.
//
// A note about uniqueness derived from the UUID Wikipedia entry:
//
//  Randomly generated UUIDs have 122 random bits.  One's annual risk of being
//  hit by a meteorite is estimated to be one chance in 17 billion, that
//  means the probability is about 0.00000000006 (6 × 10−11),
//  equivalent to the odds of creating a few tens of trillions of UUIDs in a
//  year and having one duplicate.
func NewRandom() UUID {
	var result, err = uuid.NewRandom()
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

// NewSHA1 returns a new SHA1 (Version 5) UUID based on the
// supplied name space and data.  It is the same as calling:
//
//  NewHash(sha1.New(), space, data, 5)
func NewSHA1(space UUID, data []byte) UUID {
	return uuid.NewSHA1(space, data)
}
