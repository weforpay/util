package util

import (
	"testing"
)

func TestByteArray(t *testing.T) {
	var err error
	ba := NewByteArray()
	ba.WriteByte(0xaa)
	t.Logf("len:%d", ba.Length())
	ba.SetPosition(ba.Length() - 1)
	v, err := ba.ReadByte()
	if err != nil {
		t.Fatalf("ba.ReadByte err:%v", err)
		t.FailNow()
	}
	if v != 0xaa {
		t.Fatalf("v(0x%x) != 0xaa", v)
		t.FailNow()
	}
	t.Logf("len:%d", ba.Length())
	ba.WriteBoolean(true)
	ba.SetPosition(ba.Length() - 1)
	vb, err := ba.ReadBoolean()
	if err != nil {
		t.Fatalf("ba.ReadBoolean err:%v", err)
		t.FailNow()
	}
	if !vb {
		t.Fatalf("v(%v) != true", vb)
	}
	t.Logf("len:%d", ba.Length())
	ba.WriteDouble(0.5)
	t.Logf("len:%d", ba.Length())
	ba.SetPosition(ba.Length() - 8)
	vdb, err := ba.ReadDouble()
	if err != nil {
		t.Fatalf("ba.ReadDouble err:%v", err)
		t.FailNow()
	}
	if vdb != 0.5 {
		t.Fatalf("v(%f) != 0.5", vdb)
	}
	ba.WriteFloat(0.6)
	t.Logf("len:%d", ba.Length())
	ba.SetPosition(ba.Length() - 4)
	vf, err := ba.ReadFloat()
	if err != nil {
		t.Fatalf("ba.ReadFloat err:%v", err)
		t.FailNow()
	}
	if vf != 0.6 {
		t.Fatalf("v(%f) != 0.6", vf)
	}
	ba.WriteInt(0xaaaaaaa)
	t.Logf("len:%d", ba.Length())
	ba.SetPosition(ba.Length() - 4)
	vn, err := ba.ReadInt()
	if err != nil {
		t.Fatalf("ba.ReadInt err:%v", err)
		t.FailNow()
	}
	if vn != 0xaaaaaaa {
		t.Fatalf("v(0x%x) != 0xaaaaaaaa", vn)
	}
	ba.WriteShort(0xaaa)
	t.Logf("len:%d", ba.Length())
	ba.SetPosition(ba.Length() - 2)
	vs, err := ba.ReadShort()
	if err != nil {
		t.Fatalf("ba.ReadShort err:%v", err)
		t.FailNow()
	}
	if vs != 0xaaa {
		t.Fatalf("v(0x%x) != 0xaaa", vs)
	}

	ba.WriteUnsignedInt(0xaaaaaaa)
	t.Logf("len:%d", ba.Length())
	ba.SetPosition(ba.Length() - 4)
	vun, err := ba.ReadUnsignedInt()
	if err != nil {
		t.Fatalf("ba.ReadUnsignedInt err:%v", err)
		t.FailNow()
	}
	if vun != 0xaaaaaaa {
		t.Fatalf("v(0x%x) != 0xaaaaaaa", vun)
	}
	ban := NewByteArray()
	ban.WriteByteArray(ba, 0, 0)

	t.Logf("ban.len:%d", ban.Length())
	ban.SetPosition(ban.Length() - 4)
	vun, err = ban.ReadUnsignedInt()
	if err != nil {
		t.Fatalf("ban.ReadUnsignedInt err:%v", err)
		t.FailNow()
	}
	if vun != 0xaaaaaaa {
		t.Fatalf("v(0x%x) != 0xaaaaaaa", vun)
	}

	ba.ReadByteArray(ban, 0, 0)

	t.Logf("ban.len:%d", ban.Length())
	ban.SetPosition(ban.Length() - 4)
	vun, err = ban.ReadUnsignedInt()
	if err != nil {
		t.Fatalf("ban.ReadUnsignedInt err:%v", err)
		t.FailNow()
	}
	if vun != 0xaaaaaaa {
		t.Fatalf("v(0x%x) != 0xaaaaaaa", vun)
	}

	ba.WriteBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd}, 0, 0)

	t.Logf("len:%d", ba.Length())

	buf := make([]byte, 4)
	err = ba.ReadBytes(buf, ba.Length()-4, 4)
	if err != nil {
		t.Fatalf("ban.ReadUnsignedInt err:%v", err)
		t.FailNow()
	}
	if buf[3] != 0xdd {
		t.Fatalf("buf[3](0x%x) != 0xdd", buf[3])
	}
}
