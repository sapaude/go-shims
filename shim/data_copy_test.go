package shim

import (
    "testing"
)

func TestDeepCopyByJSON(t *testing.T) {

    type address struct {
        Address  string
        Position string
    }

    type userAddr struct {
        Name    string `json:"name,omitempty"`
        Sex     string
        Address *address `json:"address,omitempty"`
        age     int
    }

    ua := &userAddr{
        Name: "a",
        Sex:  "6",
        Address: &address{
            Address:  "HK",
            Position: "66.66",
        },
        age: 18,
    }
    var ub *userAddr

    // copy
    err := DeepCopyByJSON(ua, &ub)
    if err != nil {
        t.Fatal(err)
    }

    // update
    ua.Address.Position = "77.77" // ua: (HK, 77.77)
    ub.Address.Address = "AM"     // ub: (AM, 66.66)

    t.Logf("ua: %s, age: %d", ToJsonString(ua, false), ua.age)
    t.Logf("ub: %s, age: %d", ToJsonString(ub, false), ub.age)
}
