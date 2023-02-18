package gamemap

import "github.com/VxVxN/game/internal/base"

type Layer [][]MapTile

func NewLayer(width, height int) Layer {
	tiles := make([][]MapTile, width)
	for x := 0; x < width; x++ {
		tiles[x] = make([]MapTile, height)
	}
	return tiles
}

type LayerContainer struct {
	width          int
	height         int
	currentElement int
	layers         []Layer
}

func NewLayerContainer(width, height int) *LayerContainer {
	return &LayerContainer{
		width:          width,
		height:         height,
		currentElement: 0,
		layers: []Layer{
			{},
		},
	}
}

func (container *LayerContainer) Elements() []Layer {
	return container.layers
}

func (container *LayerContainer) Next() Layer {
	container.currentElement++
	if len(container.layers) == container.currentElement {
		container.layers = append(container.layers, NewLayer(container.width, container.height))
	}
	return container.layers[container.currentElement]
}

func (container *LayerContainer) SetCurrent(layer Layer) {
	container.layers[container.currentElement] = layer
}

func (container *LayerContainer) GetLayerWithoutCollisions(positions []base.Position) Layer {
	var collisionLayend int
	for i, layer := range container.layers {
		if i == 0 {
			continue // skip background
		}
		for _, position := range positions {
			if existTile(layer, position.X, position.Y) {
				collisionLayend = i
				break
			}
		}
	}
	if collisionLayend != 0 {
		container.currentElement = collisionLayend
		container.currentElement++

		if len(container.layers) == container.currentElement {
			container.layers = append(container.layers, NewLayer(container.width, container.height))
		}
		return container.layers[container.currentElement]
	}
	return container.layers[container.currentElement]
}

func (container *LayerContainer) GetIndex() int {
	return container.currentElement
}

func (container *LayerContainer) SetIndex(index int) {
	container.currentElement = index
}
