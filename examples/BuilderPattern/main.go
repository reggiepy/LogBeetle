package main

import (
	"fmt"
)

type Speed float64

const (
	MPH Speed = 1
	KPH       = 1.60934
)

type Color string

const (
	BlueColor  Color = "blue"
	GreenColor       = "green"
	RedColor         = "red"
)

type Wheels string

const (
	SportsWheels Wheels = "sports"
	SteelWheels         = "steel"
)

type Builder interface {
	Color(Color) Builder
	Wheels(Wheels) Builder
	TopSpeed(Speed) Builder
	Build() Interface
}

type Interface interface {
	Drive() error
	Stop() error
}

type MyCarBuilder struct {
	color    Color
	wheels   Wheels
	topSpeed Speed
}

func NewMyCarBuilder() *MyCarBuilder {
	return &MyCarBuilder{}
}

func (b *MyCarBuilder) Color(c Color) Builder {
	b.color = c
	return b
}

func (b *MyCarBuilder) Wheels(w Wheels) Builder {
	b.wheels = w
	return b
}

func (b *MyCarBuilder) TopSpeed(s Speed) Builder {
	b.topSpeed = s
	return b
}

func (b *MyCarBuilder) Build() Interface {
	return &MyCar{
		color:    b.color,
		wheels:   b.wheels,
		topSpeed: b.topSpeed,
	}
}

type MyCar struct {
	color    Color
	wheels   Wheels
	topSpeed Speed
}

func (c *MyCar) Drive() error {
	// 实现驾驶汽车的逻辑
	return nil
}

func (c *MyCar) Stop() error {
	// 实现停止汽车的逻辑
	return nil
}

func main() {
	builder := NewMyCarBuilder()
	car := builder.Color(BlueColor).Wheels(SportsWheels).TopSpeed(MPH).Build()
	fmt.Println(car)
}
