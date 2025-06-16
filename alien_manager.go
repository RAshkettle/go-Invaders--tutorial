package main

const (
	NUMBER_OF_ALIENS_IN_ROW = 12
	ALIEN_SIZE              = 16
	PADDING                 = 64
)

func SpawnAlienWave() []*Alien {
	aliens := make([]*Alien, 0)

	for i := range NUMBER_OF_ALIENS_IN_ROW {
		alien := NewAlien(SquidAlien)
		alien.X = (i * ALIEN_SIZE) + PADDING
		alien.Y = ALIEN_SIZE

		aliens = append(aliens, alien)
	}
	return aliens
}
