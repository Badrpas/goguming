package foight

type UpdateMessage struct {
  dx, dy int8

}

func decodeUpdateMessage(data []byte) UpdateMessage {
  um := UpdateMessage{

  }
  um.decode(data);
  return um;
}

func (um *UpdateMessage) decode(data []byte) {
  um.dx = int8(data[0]) - 1;
  um.dy = int8(data[1]) - 1;

}