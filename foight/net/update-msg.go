package net

type UpdateMessage struct {
	Dx, Dy int8
	Tx, Ty int8
}

func DecodeUpdateMessage(data []byte) UpdateMessage {
	um := UpdateMessage{}
	um.decode(data)
	return um
}

func (um *UpdateMessage) decode(data []byte) {
	um.Dx = int8(uint8(data[0]) - 50)
	um.Dy = int8(uint8(data[1]) - 50)
	um.Tx = int8(uint8(data[2]) - 50)
	um.Ty = int8(uint8(data[3]) - 50)
}
