package storagesc

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *BlobberPartitionsWeights) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "ps"
	o = append(o, 0x81, 0xa2, 0x70, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Parts)))
	for za0001 := range z.Parts {
		// map header, size 2
		// string "i"
		o = append(o, 0x82, 0xa1, 0x69)
		o = msgp.AppendInt(o, z.Parts[za0001].Index)
		// string "w"
		o = append(o, 0xa1, 0x77)
		o = msgp.AppendInt(o, z.Parts[za0001].Weight)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlobberPartitionsWeights) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ps":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Parts")
				return
			}
			if cap(z.Parts) >= int(zb0002) {
				z.Parts = (z.Parts)[:zb0002]
			} else {
				z.Parts = make([]PartitionWeightBlobber, zb0002)
			}
			for za0001 := range z.Parts {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Parts", za0001)
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						err = msgp.WrapError(err, "Parts", za0001)
						return
					}
					switch msgp.UnsafeString(field) {
					case "i":
						z.Parts[za0001].Index, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							err = msgp.WrapError(err, "Parts", za0001, "Index")
							return
						}
					case "w":
						z.Parts[za0001].Weight, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							err = msgp.WrapError(err, "Parts", za0001, "Weight")
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							err = msgp.WrapError(err, "Parts", za0001)
							return
						}
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *BlobberPartitionsWeights) Msgsize() (s int) {
	s = 1 + 3 + msgp.ArrayHeaderSize + (len(z.Parts) * (5 + msgp.IntSize + msgp.IntSize))
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BlobberWeight) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "bid"
	o = append(o, 0x82, 0xa3, 0x62, 0x69, 0x64)
	o = msgp.AppendString(o, z.BlobberID)
	// string "w"
	o = append(o, 0xa1, 0x77)
	o = msgp.AppendInt(o, z.Weight)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlobberWeight) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "bid":
			z.BlobberID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlobberID")
				return
			}
		case "w":
			z.Weight, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Weight")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z BlobberWeight) Msgsize() (s int) {
	s = 1 + 4 + msgp.StringPrefixSize + len(z.BlobberID) + 2 + msgp.IntSize
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PartitionWeightBlobber) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "i"
	o = append(o, 0x82, 0xa1, 0x69)
	o = msgp.AppendInt(o, z.Index)
	// string "w"
	o = append(o, 0xa1, 0x77)
	o = msgp.AppendInt(o, z.Weight)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PartitionWeightBlobber) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "i":
			z.Index, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Index")
				return
			}
		case "w":
			z.Weight, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Weight")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z PartitionWeightBlobber) Msgsize() (s int) {
	s = 1 + 2 + msgp.IntSize + 2 + msgp.IntSize
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PartitionWeightsBlobber) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "ps"
	o = append(o, 0x81, 0xa2, 0x70, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Parts)))
	for za0001 := range z.Parts {
		// map header, size 2
		// string "i"
		o = append(o, 0x82, 0xa1, 0x69)
		o = msgp.AppendInt(o, z.Parts[za0001].Index)
		// string "w"
		o = append(o, 0xa1, 0x77)
		o = msgp.AppendInt(o, z.Parts[za0001].Weight)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PartitionWeightsBlobber) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ps":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Parts")
				return
			}
			if cap(z.Parts) >= int(zb0002) {
				z.Parts = (z.Parts)[:zb0002]
			} else {
				z.Parts = make([]PartitionWeightBlobber, zb0002)
			}
			for za0001 := range z.Parts {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Parts", za0001)
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						err = msgp.WrapError(err, "Parts", za0001)
						return
					}
					switch msgp.UnsafeString(field) {
					case "i":
						z.Parts[za0001].Index, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							err = msgp.WrapError(err, "Parts", za0001, "Index")
							return
						}
					case "w":
						z.Parts[za0001].Weight, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							err = msgp.WrapError(err, "Parts", za0001, "Weight")
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							err = msgp.WrapError(err, "Parts", za0001)
							return
						}
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PartitionWeightsBlobber) Msgsize() (s int) {
	s = 1 + 3 + msgp.ArrayHeaderSize + (len(z.Parts) * (5 + msgp.IntSize + msgp.IntSize))
	return
}
