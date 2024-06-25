package main

type Point struct {
    X int
    Y int
}

type Snake struct {
    Body        []Point
    Direction   Point
    GrowCounter int
}

func NewSnake() *Snake {
    return &Snake{
        Body: []Point{
            {X: screenWidth / tileSize / 2, Y: screenHeight / tileSize / 2},
        },
        Direction: Point{X: 1, Y: 0},
    }
}

func (s *Snake) Move() {
    newHead := Point{
        X: s.Body[0].X + s.Direction.X,
        Y: s.Body[0].Y + s.Direction.Y,
    }
    s.Body = append([]Point{newHead}, s.Body...)

    if s.GrowCounter > 0 {
        s.GrowCounter--
    } else {
        s.Body = s.Body[:len(s.Body)-1]
    }
}
