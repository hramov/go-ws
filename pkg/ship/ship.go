package ship

type Ship struct {
	ID         int  `json: id`
	Player     int  `json: "player"`
	Length     int  `json: "length"`
	StartX     int  `json: "startX"`
	StartY     int  `json: "startY"`
	Direction  int  `json: "direction"`
	Hit        bool `json: "hit"`
	LivePoints int  `json: "livePoints"`
	Live       bool `json: "live"`
}
